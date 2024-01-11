package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darrenparkinson/htmx-contact-app-go/internal/ports"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/repositories"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/usecases"
	"github.com/gin-gonic/gin"
)

type application struct {
	usecase ports.ContactUseCase
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	contactRepo := repositories.NewContactJSONRepository()
	contactUseCase := usecases.NewContactUseCase(contactRepo)
	app := application{
		usecase: contactUseCase,
	}
	router := gin.Default()
	router.GET("/contacts", app.listContacts)
	log.Println("listening...")
	router.Run(fmt.Sprintf(":%d", 8080))
}

func (app *application) listContacts(c *gin.Context) {
	contacts, err := app.usecase.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, contacts)
}
