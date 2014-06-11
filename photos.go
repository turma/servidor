package main

import (
	"log"
	"time"
)

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

	//log.Printf("Photo: %dx%d %s ", fbImage.Width, fbImage.Height, fbImage.Id)

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
			Likes:       fbLikesSummary.Summary.TotalCount,
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
	} else {
		log.Printf("Photo '%s' already exists.", fbImage.Id)
		return
	}

	//
	// I'm not saving that this user shared this photo
	//

}
