package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Environment struct {
	Url        string
	Port       int
	Production bool
	AppId      string
	AppSecret  string
	DBName     string
	DBUser     string
	DBPass     string

	Youtube string
}

// My Production configs
var EnvProd = Environment{
	Port:       8000,
	Url:        "http://tur.ma/",
	Production: true,
	AppId:      "486327838135438",
	AppSecret:  "327e1e72f15e4debca506a3a2580f02d",
	DBName:     "turmadatabase",
	DBUser:     "app",
	DBPass:     "SecretPassword!",

	Youtube: "AIzaSyAd8uyRU5NeutRtfy9bOrYU1slvq5WSb9g",
}

// My Development configs
var EnvDev = Environment{
	Port:       8000,
	Url:        "http://localhost:8000/",
	Production: false,
	AppId:      "494483310653224",
	AppSecret:  "ea444f34fa4a92c6fd13707b83924ddc",
	DBName:     "turmadatabase",
	DBUser:     "app",
	DBPass:     "SecretPassword!",

	Youtube: "AIzaSyAd8uyRU5NeutRtfy9bOrYU1slvq5WSb9g",
}

// The only one martini instance
var m *martini.Martini

var renderOptions = render.Options{
	// Directory to load templates. Default is "templates"
	Directory: "views",
	// Extensions to parse template files from. Defaults to [".tmpl"]
	Extensions: []string{".html"},
	// Delims sets the action delimiters to the specified strings in the Delims struct.
	Delims: render.Delims{Left: "{#", Right: "#}"},
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset: "UTF-8",
	// Outputs human readable JSON
	IndentJSON: true,
}

func init() {
	m = martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	// Serving public directory
	m.Use(martini.Static("public"))

	// Render html templates from /views directory
	m.Use(render.Renderer(renderOptions))

	// Setup routes
	r := martini.NewRouter()

	r.Get("/", func(r render.Render) {
		template := make(map[string]interface{})
		if martini.Env == "production" {
			template["AppId"] = EnvProd.AppId
		} else {
			template["AppId"] = EnvDev.AppId
		}
		r.HTML(200, "index", template)
	})

	// Api calls
	// Explicity where we are using each middleare
	r.Post("/api/me", OrmMiddleware, FbMiddleware, PostMeHandler)

	//
	// Run the aggregator bot
	//
	r.Get("/aggregator", OrmMiddleware, FbMiddleware, AggregatorHandler)

	// Just a ping route
	r.Get("/ping", func() string {
		return "pong!"
	})

	r.NotFound(func(r render.Render, req *http.Request) {
		// !!! Need to check if the call isn't to /api, so return default html not found
		r.JSON(http.StatusNotFound, fmt.Sprintf(
			"Desculpe, mas nao ha nada no endereço requisitado. [%s] %s", req.Method, req.RequestURI))
	})

	// Add the routers to my martini instances
	m.Action(r.Handle)
}

func main() {
	var port int
	if martini.Env == "production" {
		log.Println("Server starting in production in " + EnvProd.Url)
		port = EnvProd.Port
	} else {
		log.Println("Server starting in development in " + EnvDev.Url)
		port = EnvProd.Port
	}

	// Starting de HTTP server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), m); err != nil {
		log.Fatal(err)
	}
}
