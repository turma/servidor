package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var Env struct {
	Url        string
	Port       int
	Production bool
	AppId      string
	AppSecret  string
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
	if martini.Env == "production" {
		log.Println("Server in production")
		Env.Port = 8000
		Env.Url = "http://tur.ma/"
		Env.Production = true
		Env.AppId = "486327838135438"
		Env.AppSecret = "327e1e72f15e4debca506a3a2580f02d"
	} else {
		log.Println("Server in development")
		Env.Port = 8000
		Env.Url = fmt.Sprintf("http://localhost:%d/", Env.Port)
		Env.Production = false
		Env.AppId = "494483310653224"
		Env.AppSecret = "ea444f34fa4a92c6fd13707b83924ddc"
	}

	m = martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	// Serving public directory
	m.Use(martini.Static("public"))

	// Render html templates from /views directory
	m.Use(render.Renderer(renderOptions))

	// Add the OrmMiddleware
	//m.Use(OrmMiddleware)

	// Add the AuthMiddleware
	//m.Use(AuthMiddleware)

	// Setup routes
	r := martini.NewRouter()

	r.Get("/", func(r render.Render) {
		template := make(map[string]interface{})
		template["Env"] = Env

		r.HTML(200, "index", template)
	})

	// Api calls
	r.Post("/api/me", PostMeHandler)

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
	// Starting de HTTP server
	log.Println("Starting HTTP server in " + Env.Url + " ...")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Env.Port), m); err != nil {
		log.Fatal(err)
	}
}
