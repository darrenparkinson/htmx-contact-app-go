package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/darrenparkinson/htmx-contact-app-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type newContactForm struct {
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Phone     string `form:"phone"`
	Email     string `form:"email"`
}

func (app *application) contacts(c *gin.Context) {
	var contacts []domain.Contact
	search := c.Query("q")
	if search != "" {
		contacts, _ = app.contactsUseCase.Search(search)
	} else {
		contacts, _ = app.contactsUseCase.List()
	}
	if c.GetHeader("HX-Trigger") == "search" {
		c.HTML(http.StatusOK, "rows.html", gin.H{"contacts": contacts, "search": search})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"contacts": contacts, "messages": flashes(c)})
}

func (app *application) count(c *gin.Context) {
	count := app.contactsUseCase.Count()
	c.String(http.StatusOK, fmt.Sprintf("(%d total Contacts)", count))
}

func (app *application) contactsNewGet(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", gin.H{"contact": domain.Contact{}, "messages": flashes(c)})
}

func (app *application) contactsNew(c *gin.Context) {
	var f newContactForm
	c.Bind(&f)
	err := app.contactsUseCase.Create(f.FirstName, f.LastName, f.Phone, f.Email)
	if err != nil {
		contact := domain.Contact{
			First: f.FirstName,
			Last:  f.LastName,
			Phone: f.Phone,
			Email: f.Email,
		}
		c.HTML(http.StatusOK, "new.html", gin.H{"contact": contact})
		return
	}
	flashMessage(c, "Created New Contact!")
	c.Redirect(http.StatusFound, "/contacts")
}

func (app *application) contactsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contact, _ := app.contactsUseCase.Find(id)
	c.HTML(http.StatusOK, "show.html", gin.H{"contact": contact, "messages": flashes(c)})
}
