package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/martini-contrib/render"
)

//
// This aggregator get all links and videos posted by users in the last hour
// and update it in our Data Base
// It should refresh the cache with the first 10 pages, or more, of information.
//
func AggregatorHandler(db DB, fb FB, r render.Render) {
	log.Println("Received AggregatorHandler")

	// List of users who shared his contents and
	// we still have a valid long lived access token

	// maybe not needed: read_stream_permission = true and
	query := "select * from user where token_expires_at >= now()"
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

// This regex checks if link url is inside facebook
var isFacebookObject = regexp.MustCompile(`^https?:\/\/www.facebook.com\/`)

// This regex checks if source is from an youtube video
var isYouTubeVideo = regexp.MustCompile(`^https?:\/\/www.youtube.com\/v\/`)

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
		log.Printf("POST: %s, Type: ", feed.Id, feed.Type)
		if feed.Type == "link" {
			log.Printf("Feed type:%s, id:%s", feed.Type, feed.Id)

			// It will insert links that isn't a link to an facebook object
			if !isFacebookObject.MatchString(feed.Link) {
				InsertLink(fb, user, feed)
			} else {
				//
				// I'm not interested in saving facebook communities and pages
				//
				log.Printf("Facebook object not saved link:%s", feed.Type, feed.Link)
			}

		}
		if feed.ObjectId != "" {
			// There is an Facebook Object shared in this post

			log.Printf("OID Feed type:%s, object_id:%s", feed.Type, feed.ObjectId)
			if feed.Type == "photo" {
				InsertPhoto(fb, user, feed)
			} else if feed.Type == "video" {
				InsertVideo(fb, user, feed)
			}
		}
		if feed.ObjectId == "" && (feed.Type == "swf" || feed.Type == "video") {
			log.Printf("YOUTUBE VIDEO type:%s, object_id:%s, Link:%s", feed.Type, feed.ObjectId, feed.Link)

			// If it's an YouTube video, insert it too
			if isYouTubeVideo.MatchString(feed.Source) {
				InsertYouTubeVideo(fb, user, feed)
			}

		}
	}
}
