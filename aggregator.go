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
	"limit": 100,
	// Removed LIKES from the fields: likes.fields(id,name,picture)
	"fields": "id,name,description,message,link,source,picture,caption,created_time,privacy,type,object_id",
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
			//log.Printf("Feed type:%s, id:%s", feed.Type, feed.Id)
			//InsertLink(fb, user, feed)
		}
		if feed.ObjectId != "" {
			//log.Printf("OID Feed type:%s, object_id:%s", feed.Type, feed.ObjectId)
			if feed.Type == "photo" {
				InsertPhoto(fb, user, feed)
				log.Printf("POST: %s", feed.Id)
			} else if feed.Type == "video" {
				//InsertVideo(fb, user, feed)
			}
		}
		if feed.ObjectId == "" && (feed.Type == "swf" || feed.Type == "video") {
			//log.Printf("YOUTUBE VIDEO type:%s, object_id:%s, Link:%s", feed.Type, feed.ObjectId, feed.Link)
			// Check if Link is from Facebook, and take it there
		}
	}
}

var FBVideoParams = facebook.Params{
	"fields": "id,from,name,description,embed_html,source,picture,created_time,format",
}

//
// This function inserts and/or updates videos
//
func InsertVideo(fb FB, user *User, feed FBFeedData) {
	// Try to get this facebook object
	res, err := fb.Get("/"+feed.ObjectId, FBVideoParams)
	if err != nil { // Facebook object_id not exists
		//log.Printf("Error getting object '%s' from facebook! %s.", feed.ObjectId, err)
		return
	}

	var fbVideo FBVideo
	err = res.Decode(&fbVideo)
	if err != nil {
		log.Printf("Error decoding video object '%s'! %s.", feed.ObjectId, err)
		return
	}

	log.Printf("Video: %s", fbVideo.Id)
}

var FBPhotoParams = facebook.Params{
	// Another image sizes not in use, removed from fields: images
	"fields": "id,from.fields(id,name,picture),name,height,width,link,source,picture,created_time",
}

var FBSummaryParams = facebook.Params{
	"limit":   0, // I don't need the likes information
	"summary": true,
}

//
// This function inserts and/or updates photos
//
func InsertPhoto(fb FB, user *User, feed FBFeedData) {
	// Try to get this facebook object
	res, err := fb.Get("/"+feed.ObjectId, FBPhotoParams)
	if err != nil { // Facebook object_id not exists
		//log.Printf("Error getting object '%s' from facebook! %s.", feed.ObjectId, err)
		return
	}

	var fbImage FBImage
	err = res.Decode(&fbImage)
	if err != nil {
		log.Printf("Error decoding image object '%s'! %s.", feed.ObjectId, err)
		return
	}

	log.Printf("Photo: %dx%d %s ", fbImage.Width, fbImage.Height, fbImage.Id)

	photoobj, err := db.Get(Photo{}, fbImage.Id)
	if err != nil {
		log.Printf("Error getting the photo '%s'. %s", fbImage.Id, err)
		return
	}

	if photoobj == nil {

		// Try to get the Like Count of this image
		res, err := fb.Get("/"+fbImage.Id+"/likes", FBSummaryParams)
		if err != nil {
			log.Printf("Error getting photo '%s' likes! %s.", fbImage.Id, err)
			return
		}

		var fbLikesSummary FBLikesSummary
		err = res.Decode(&fbLikesSummary)
		if err != nil {
			log.Printf("Error decoding photo '%s' summary! %s.", fbImage.Id, err)
			return
		}

		// Extract a time format from the facebook json time format
		createdTime, err := time.Parse(FBTimeLayout, fbImage.CreatedTime)
		if err != nil {
			log.Printf("Error extracting CreatedTime from a photo. %s", err)
			return
		}

		photo := &Photo{
			Id:          fbImage.Id,
			Likes:       fbLikesSummary.Summary.TotalCount, // NEED TO GET IT YET
			Name:        fbImage.Name,
			Height:      fbImage.Height,
			Width:       fbImage.Width,
			Link:        fbImage.Link,
			Source:      fbImage.Source,
			Picture:     fbImage.Picture,
			CreatedTime: createdTime,
			FromId:      fbImage.From.Id,
			FromName:    fbImage.From.Name,
			FromPicture: fbImage.From.Picture.Data.Url,
		}

		err = db.Insert(photo)
		if err != nil {
			log.Printf("Error insterting a new photo '%s'! %s.", fbImage.Id, err)
			return
		}

		log.Printf("New photo '%s' likes: %d!", photo.Id, photo.Likes)
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

}

type FBPrivacy struct {
	Value string
}

type FBShares struct {
	Shares int
}

type FBImage struct {
	Id          string // user id _ post id
	Name        string
	Height      int
	Width       int
	Link        string
	Source      string // Photo URL
	Picture     string
	CreatedTime string
	From        FBFrom

	// Images []FBImages
}

// Not in use, we will use de default 720x720 max image
//
// type FBImages struct {
// 	Source string
// 	Height int
// 	Width  int
// }

type FBVideo struct {
	Id          string // user id _ post id
	Name        string
	Description string
	EmbedHtml   string
	Source      string // For video player
	Picture     string
	CreatedTime string

	Format []FBFormat

	From FBFrom
}

type FBFormat struct {
	EmbedHtml string
	Width     int
	Height    int
	Filter    string
	// 130x130 (real:130x98), 480x480 (real:480x380), native (real:640x480)
}

type FBFrom struct {
	Name    string
	Id      string
	Picture FBPicture
}

type FBLikesSummary struct {
	Summary FBSummary
}

type FBSummary struct {
	TotalCount int
}

//
// LIKES REMOVED
//

// type FBLikes struct {
// 	Data []FBLikesData
// }

// type FBLikesData struct {
// 	Id      string
// 	Name    string
// 	Picture FBPictureData
// }
