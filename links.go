package main

import (
	"log"
	"net/url"
	"time"
)

//
// This function inserts and/or updates links
// If this link doesn't exist in our DB, it will be inserted
// If there is an object associated with this link, this object info will be saved
// If it is a normal link, the oldest information about this link will be saved.
//
// it means that if an user saved it before, and the current user shared it first,
// the information added by this user who shared it first will be shown
//
func InsertLink(fb FB, user *User, feed FBFeedData) {

	// Escapping URL for safely place it inside the facebook api query
	escapedLink := url.QueryEscape(feed.Link)

	// Getting info about this link
	// We don't know if it will return an facebook object or not
	res, err := fb.Get(escapedLink, nil)
	if err != nil {
		log.Printf("Error getting link '%s'! %s.", feed.Link, err)
		return
	}

	var fbLinkInfo FBLinkInfo
	err = res.Decode(&fbLinkInfo)
	if err != nil {
		log.Printf("Error decodding link '%s'! %s.", feed.Link, err)
		return
	}

	// If it is an Facebook object, so we should use its link
	if fbLinkInfo.Link != "" {
		feed.Link = fbLinkInfo.Link
	}

	// Check if this link is already saved
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

	// If link doesn't exists, save it in our DB
	if linkobj == nil {

		// Checks if there is some facebook object associated with this link
		if fbLinkInfo.Link != "" {

			// Get this Link object in facebook
			res, err := fb.Get(fbLinkInfo.Id, FBLinkParams)
			if err != nil {
				log.Printf("Error getting link '%s'! %s.", fbLinkInfo.Id, err)
				return
			}

			var fbLink FBLink
			err = res.Decode(&fbLink)
			if err != nil {
				log.Printf("Error decoding link '%s'! %s.", fbLinkInfo.Id, err)
				return
			}

			link := &Link{
				Link:        fbLink.Link,
				Id:          fbLink.Id,    // Some links is a facebook object
				Likes:       fbLink.Likes, // Likes for facebook objects
				Name:        fbLink.Name,
				Description: fbLink.Description,
				Caption:     feed.Caption,
				Picture:     fbLink.Picture.Data.Url,
				CreatedTime: createdTime,
			}

			err = db.Insert(link)
			if err != nil {
				log.Printf("Error insterting a new link '%s'! %s.", link.Link, err)
				return
			}

			log.Printf("Link inserted '%s'!.", link.Link)

		} else { // This link has not an object id associated with it

			link := &Link{
				Link:        feed.Link,
				Id:          "",                // It's a normal link
				Likes:       fbLinkInfo.Shares, // Shares for links
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

			log.Printf("Link inserted '%s'!.", link.Link)

		}
	} else {

		// This link already exists in DB
		link := linkobj.(*Link)

		// If it's not an facebook object and this user has shared it first
		// so his information about the link will be shown
		if link.Id == "" && createdTime.Before(link.CreatedTime) {
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

			log.Printf("Link updated '%s'!.", link.Link)
		}
	}

	// After save the link,
	// We will save the info shared by this user about the current link
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

		log.Printf("Saved shared info for '%s'!.", feed.Link)
	}
}
