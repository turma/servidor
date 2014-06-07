package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/huandu/facebook"
	"github.com/martini-contrib/render"
)

type postMeHandlerData struct {
	AccessToken string
	UserId      string
}

func PostMeHandler(r render.Render, req *http.Request) {
	var data postMeHandlerData
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding body PostMeHandler!!")
		r.JSON(http.StatusNotAcceptable, fmt.Sprintf(
			"Nao foi possivel decodificar o objeto Json! %s.", err))
		return
	}

	log.Println("Received PostMeHandler")

	// create a global App var to hold your app id and secret.
	var fb = facebook.New(Env.AppId, Env.AppSecret)
	fb.RedirectUri = Env.Url

	session := fb.Session(data.AccessToken)

	log.Println("Testing the time...")
	id, err := session.User()
	if err != nil {
		log.Println("Sua sessao nao e valida...")
		r.JSON(http.StatusNotAcceptable, fmt.Sprintf(
			"Sua sessao nao e valida...! %s.", err))
		return
	}
	log.Printf("UserId: %s", id)

	token, expires, err := fb.ExchangeToken(data.AccessToken)
	if err != nil {
		r.JSON(http.StatusNotAcceptable, fmt.Sprintf(
			"Erro ao tentar adquirir um Token de vida longa! %s.", err))
		return
	}

	fmt.Printf("Tok: %s, exp: %d.", token, expires)

	r.JSON(http.StatusOK, data)
	return
}
