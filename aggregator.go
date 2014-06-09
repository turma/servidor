package main

import (
	"log"
	"net/http"
	"time"

	"github.com/huandu/facebook"
	"github.com/martini-contrib/render"
)

const FBTimeLayout = "2006-01-02T15:04:05+0000"

//
// This aggregator get all links and videos posted by users in the last hour
// and update it in our Data Base
// It should refresh the an cache with the first 10 pages, or more, of information.
//
func AggregatorHandler(db DB, fb FB, r render.Render) {
	log.Println("Received AggregatorHandler")

	// List of users who shared his contents and
	// we still have a valid long lived access token
	query := "select * from user where read_stream_permission = true and token_expires_at >= now()"
	var users []*User
	_, err := db.Select(&users, query)
	if err != nil {
		log.Printf("Error selecting the list of users. %s", err)
	}

	for _, user := range users {
		session := fb.App().Session(user.LongLivedAccessToken)
		go Aggregator(session, user)
	}

	// Request accepted!
	r.JSON(http.StatusOK, true)
	return
}

var FBFeedParams = facebook.Params{
	"limit":  50,
	"fields": "id,name,description,message,link,source,picture,caption,created_time,privacy,type,object_id,likes.fields(id,name,picture)",
}

//
// This function will receives an FB session,
// search and aggregate all links and videos in this session
//
func Aggregator(fb FB, user *User) {
	res, err := fb.Get("/me/feed", FBFeedParams)
	if err != nil {
		log.Printf("Error getting user feeds. %s", err)
		return
	}

	var fbFeed FBFeed
	err = res.Decode(&fbFeed)
	if err != nil {
		log.Printf("Error decoding user '%s' feed! %s.", user.Id, err)
		return
	}

	for _, feed := range fbFeed.Data {
		if feed.Type == "link" {
			log.Printf("Feed type:%s, id:%s", feed.Type, feed.Id)
			InsertLink(fb, user, feed)
		}
		// if feed.ObjectId != "" {
		// 	log.Printf("Object_id:%s", feed.ObjectId)
		// }
	}
}

//
// This function inserts and/or updates links
//
func InsertLink(fb FB, user *User, feed FBFeedData) {
	linkobj, err := db.Get(Link{}, feed.Link)
	if err != nil {
		log.Printf("Error getting the link '%s'. %s", feed.Link, err)
		return
	}

	// Extract a time format from the facebook json time format
	createdTime, err := time.Parse(FBTimeLayout, feed.CreatedTime)
	if err != nil {
		log.Printf("Error extracting CreatedTime from a link feed. %s", err)
		return
	}

	// If link doesn't exists, Insert it in DB,
	// else if this user shared it first, he has priority
	if linkobj == nil {

		// Getting how many peopple liked+shared+commented this link
		shares, err := GetShares(fb, feed.Link)
		if err != nil {
			log.Printf("Error getting the link shares '%s'! %s.", feed.Link, err)
			return
		}

		link := &Link{
			Link:        feed.Link,
			Shares:      shares,
			Name:        feed.Name,
			Description: feed.Description,
			Caption:     feed.Caption,
			Picture:     feed.Picture,
			CreatedTime: createdTime,
		}

		err = db.Insert(link)
		if err != nil {
			log.Printf("Error insterting a new link '%s'! %s.", link.Link, err)
			return
		}
	} else {
		// Link already exists in DB
		link := linkobj.(*Link)

		// If user shared it first
		// his information about the link will be shown
		if createdTime.After(link.CreatedTime) {
			link.Name = feed.Name
			link.Description = feed.Description
			link.Caption = feed.Caption
			link.Picture = feed.Picture
			link.CreatedTime = createdTime

			count, err := db.Update(link)
			if err != nil {
				log.Printf("Error updating link '%s'! %s.", link.Link, err)
				return
			}
			if count == 0 {
				log.Printf("Link not updated '%s'! %s.", link.Link, err)
				return
			}
		}
	}

	sharedobj, err := db.Get(Shared{}, feed.Id)
	if err != nil {
		log.Printf("Error getting the shared '%s'. %s", feed.Id, err)
		return
	}

	if sharedobj == nil {

		shared := &Shared{
			Id:          feed.Id,
			Link:        feed.Link,
			From:        user.Id,
			Name:        feed.Name,
			Description: feed.Description,
			Message:     feed.Message,
			Caption:     feed.Caption,
			Picture:     feed.Picture,
			CreatedTime: createdTime,
		}

		err = db.Insert(shared)
		if err != nil {
			log.Printf("Error insterting a new shared '%s'! %s.", feed.Id, err)
			return
		}
	}

}

func GetShares(fb FB, id string) (int, error) {
	res, err := fb.Get(id, nil)
	if err != nil {
		return 0, err
	}

	var shares FBShares
	err = res.Decode(&shares)
	if err != nil {
		return 0, err
	}

	return shares.Shares, nil
}

type FBFeed struct {
	Data []FBFeedData
}

type FBFeedData struct {
	Id          string // user id _ post id
	Name        string
	Description string
	Message     string
	Link        string
	Source      string // For video player
	Picture     string
	Caption     string // www.youtube.com for videos
	CreatedTime string
	Privacy     FBPrivacy
	Type        string // photo, video, swf
	ObjectId    string // If exist's, just share if it's accessible

	Likes FBLikes // We are not using it yet
}

type FBPrivacy struct {
	Value string
}

type FBLikes struct {
	Data []FBLikesData
}

type FBLikesData struct {
	Id      string
	Name    string
	Picture FBPictureData
}

type FBShares struct {
	Shares int
}

// "status_type": "tagged_in_photo", type="photo", object_id not work...
