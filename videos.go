package main

import (
	"log"
	"time"
)

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

	//log.Printf("Video: %s", fbVideo.Id)

	videoobj, err := db.Get(Video{}, fbVideo.Id)
	if err != nil {
		log.Printf("Error getting the video '%s'. %s", fbVideo.Id, err)
		return
	}

	if videoobj == nil {
		// Try to get the Like Count of this image
		res, err := fb.Get("/"+fbVideo.Id+"/likes", FBSummaryParams)
		if err != nil {
			log.Printf("Error getting video '%s' likes! %s.", fbVideo.Id, err)
			return
		}

		var fbLikesSummary FBLikesSummary
		err = res.Decode(&fbLikesSummary)
		if err != nil {
			log.Printf("Error decoding video '%s' summary! %s.", fbVideo.Id, err)
			return
		}

		// Extract a time format from the facebook json time format
		createdTime, err := time.Parse(FBTimeLayout, fbVideo.CreatedTime)
		if err != nil {
			log.Printf("Error extracting CreatedTime from a video. %s", err)
			return
		}

		video := &Video{
			Id:          fbVideo.Id,
			Likes:       fbLikesSummary.Summary.TotalCount,
			Name:        fbVideo.Name,
			Description: fbVideo.Description,
			Source:      fbVideo.Source,
			Picture:     fbVideo.Picture,
			CreatedTime: createdTime,
			FromId:      fbVideo.From.Id,
			FromName:    fbVideo.From.Name,
			FromPicture: fbVideo.From.Picture.Data.Url,
		}

		err = db.Insert(video)
		if err != nil {
			log.Printf("Error insterting a new video '%s'! %s.", fbVideo.Id, err)
			return
		}

		log.Printf("New video '%s' likes: %d!", video.Id, video.Likes)
	} else {
		log.Printf("Video '%s' already exists.", fbVideo.Id)
		return
	}

	//
	// I'm not saving that this user shared this photo
	//
}
