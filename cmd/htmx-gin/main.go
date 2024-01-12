package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darrenparkinson/htmx-contact-app-go/internal/ports"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/repositories"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/usecases"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const sessionKey = "hypermedia rocks" // obviously insecure, don't do this!

type application struct {
	contactsUseCase ports.ContactUseCase
}

func main() {
	contactRepo := repositories.NewContactJSONRepository()
	contactUseCase := usecases.NewContactUseCase(contactRepo)
	app := application{
		contactsUseCase: contactUseCase,
	}
	router := gin.Default()

	store := cookie.NewStore([]byte(sessionKey))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/contacts") })
	contactsGroup := router.Group("contacts")
	contactsGroup.GET("/", app.contacts)
	contactsGroup.GET("/count", app.count)
	contactsGroup.GET("/new", app.contactsNewGet)
	contactsGroup.POST("/new", app.contactsNew)
	contactsGroup.GET("/:id", app.contactsView)
	contactsGroup.GET("/:id/edit", app.contactsEditGet)

	log.Println("listening...")
	router.Run(fmt.Sprintf(":%d", 8080))
}

func flashMessage(c *gin.Context, message string) {
	session := sessions.Default(c)
	session.AddFlash(message)
	if err := session.Save(); err != nil {
		log.Printf("error in flashMessage saving session: %s", err)
	}
}

func flashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()
	if len(flashes) != 0 {
		if err := session.Save(); err != nil {
			log.Printf("error in flashes saving session: %s", err)
		}
	}
	return flashes
}
