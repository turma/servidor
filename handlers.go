package main

import (
	"log"
	"net/http"

	"github.com/martini-contrib/render"
)

func PostMeHandler(fb FB, r render.Render) {
	log.Println("Received PostMeHandler")

	id, err := fb.User()
	if err != nil {
		log.Println("Token not accepted by Facebook")
		// The request requires User Authentication, return 401
		r.JSON(http.StatusUnauthorized, false)
		return
	}

	log.Printf("User will be insert by our FBMiddleware: %s", id)

	// Request accepted!
	r.JSON(http.StatusOK, true)
	return
}
