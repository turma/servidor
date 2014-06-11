package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/huandu/facebook"
)

// The interface to an Facebook API session
type FB interface {
	AccessToken() string
	Api(string, facebook.Method, facebook.Params) (facebook.Result, error)
	App() *facebook.App
	Batch(facebook.Params, ...facebook.Params) ([]facebook.Result, error)
	BatchApi(...facebook.Params) ([]facebook.Result, error)
	Delete(string, facebook.Params) (facebook.Result, error)
	Get(string, facebook.Params) (facebook.Result, error)
	Inspect() (facebook.Result, error)
	Post(string, facebook.Params) (facebook.Result, error)
	Put(string, facebook.Params) (facebook.Result, error)
	SetAccessToken(string)
	User() (string, error)
	Validate() error
}

// The main instance of my Facebook client holding my AppId and AppSecret
var fbClient *facebook.App

// Injects an Facebook Session in the handler
func FbMiddleware(db DB, c martini.Context, req *http.Request) {
	accessToken := req.URL.Query().Get("accessToken")

	var fb FB
	fb = fbClient.Session(accessToken)

	// If user sent us an access token,
	// We have to store it in our DB
	if accessToken != "" {
		// We save it in a new context
		// user don't have to wait for the response
		go saveSession(db, fb)
	}

	// Inject the facebook session in the handlers
	c.MapTo(fb, (*FB)(nil))
}

func init() {
	if martini.Env == "production" {
		log.Println("Facebook starting in production")
		fbClient = facebook.New(EnvProd.AppId, EnvProd.AppSecret)
		fbClient.RedirectUri = EnvProd.Url
	} else {
		log.Println("Facebook starting in development")
		fbClient = facebook.New(EnvDev.AppId, EnvDev.AppSecret)
		fbClient.RedirectUri = EnvDev.Url
	}
}

func saveSession(db DB, fb FB) {

	log.Println("Testing this session before save it.")
	id, err := fb.User()
	if err != nil {
		log.Println("Session dennied by Facebook's Server")
		return
	}

	log.Printf("User to be saved in our DB: %s", id)

	userobj, err := db.Get(User{}, id)
	if err != nil {
		log.Printf("Erro searching the user '%s' in our DB! %s.", id, err)
	}

	// It's the first access of this user
	if userobj == nil {
		log.Printf("New user '%s'.", id)

		res, err := fb.Get("/me", FBUserParams)
		if err != nil {
			log.Printf("Error getting user '%s' from Facebook! %s.", id, err)
			return
		}

		var fbUser FBUser
		err = res.Decode(&fbUser)
		if err != nil {
			log.Printf("Error decoding user '%s' data! %s.", id, err)
			return
		}

		user := decodeFBUser(&fbUser)

		user.CreatedTime = time.Now()

		// Exchange a short lived access token to a long lived access token
		longLivedAccessToken, secoundsToExpiry, err := fb.App().ExchangeToken(fb.AccessToken())
		if err != nil {
			log.Printf("Error exchanging access token '%s'! %s.", id, err)
			return
		}

		user.AccessToken = fb.AccessToken()
		user.LongLivedAccessToken = longLivedAccessToken
		user.TokenExpiresAt = time.Now().AddDate(0, 0, secoundsToExpiry/60/60/24)

		//log.Printf("\nUser= %#v", user)

		err = db.Insert(user)
		if err != nil {
			log.Printf("Error insterting the new user '%s'! %s.", id, err)
			return
		}

		log.Printf("\nUser: %s Id: %s added!", user.Name, user.Id)

	} else {
		log.Printf("Returning user '%s'.", id)

		user := userobj.(*User)

		// If user gets a sent us a new access token
		// so let's update our his long lived access token
		if user.AccessToken != fb.AccessToken() {
			log.Printf("User sent as a new access token: %v.", fb.AccessToken())

			// Exchange a short lived access token to a long lived access token
			longLivedAccessToken, secoundsToExpiry, err := fb.App().ExchangeToken(fb.AccessToken())
			if err != nil {
				log.Printf("Error exchanging access token '%s'! %s.", id, err)
				return
			}

			user.AccessToken = fb.AccessToken()
			user.LongLivedAccessToken = longLivedAccessToken
			user.TokenExpiresAt = time.Now().AddDate(0, 0, secoundsToExpiry/60/60/24)

			count, err := db.Update(user)
			if err != nil {
				log.Printf("Error updating user '%s' access token ! %s.", id, err)
				return
			}

			if count == 0 {
				log.Printf("User '%s' access token wasn't updated ! %s.", id, err)
				return
			}

			log.Printf("\nUser: %s Id: %s updated it's access token!", user.Name, user.Id)
		}

		// Should we update our long lived access token just in cases?
		// Should we update just if the last is next to expire?
		// Should we update user permissions if it changed?...
	}
}

func decodeFBUser(fbUser *FBUser) *User {
	user := &User{
		Id:       fbUser.Id,
		Email:    fbUser.Email,
		Name:     fbUser.Name,
		Picture:  fbUser.Picture.Data.Url,
		Gender:   fbUser.Gender,
		Link:     fbUser.Link,
		Locale:   fbUser.Locale,
		Timezone: fbUser.Timezone,
		Verified: fbUser.Verified,
	}

	// Scanning the permissions
	for _, perm := range fbUser.Permissions.Data {
		switch perm.Permission {
		case "email":
			if perm.Status == "granted" {
				user.EmailPermission = true
			} else {
				user.EmailPermission = false
			}
		case "read_stream":
			if perm.Status == "granted" {
				user.ReadStreamPermission = true
			} else {
				user.ReadStreamPermission = false
			}
		case "user_friends":
			if perm.Status == "granted" {
				user.UserFriendsPermission = true
			} else {
				user.UserFriendsPermission = false
			}
		}
	}

	return user
}
