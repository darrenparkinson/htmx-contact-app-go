package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type errorList struct {
	First string
	Last  string
	Phone string
	Email string
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
		errs := errorList{Email: err.Error()} // currently only email errors given, otherwise we could check errors.Is?
		c.HTML(http.StatusOK, "new.html", gin.H{"contact": contact, "errors": errs})
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

func (app *application) contactsEditGet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	contact, _ := app.contactsUseCase.Find(id)
	c.HTML(http.StatusOK, "edit.html", gin.H{"contact": contact, "messages": flashes(c)})
}

func (app *application) contactsEditPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var f newContactForm
	c.Bind(&f)
	contact := domain.Contact{
		ID:    id,
		First: f.FirstName,
		Last:  f.LastName,
		Phone: f.Phone,
		Email: f.Email,
	}
	err := app.contactsUseCase.Update(contact)
	if err != nil {
		errs := errorList{Email: err.Error()} // currently only email errors given, otherwise we could check errors.Is?
		c.HTML(http.StatusOK, "edit.html", gin.H{"contact": contact, "errors": errs})
		return
	}
	flashMessage(c, "Updated Contact!")
	c.Redirect(http.StatusFound, fmt.Sprintf("/contacts/%d", id))
}

func (app *application) contactsEmailGet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	email := c.Query("email")
	contact, _ := app.contactsUseCase.Find(id)
	contact.Email = email
	err := app.contactsUseCase.Validate(*contact)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "")
}

func (app *application) contactsDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	app.contactsUseCase.Delete(id)
	if c.GetHeader("HX-Trigger") == "delete-btn" {
		flashMessage(c, "Deleted Contact!")
		c.Redirect(http.StatusSeeOther, "/contacts")
		return
	}
	c.String(http.StatusOK, "")
}

func (app *application) contactsDeleteAll(c *gin.Context) {
	// as per the docs, the body doesn't get parsed for DELETE
	// https://pkg.go.dev/net/http#Request.ParseForm
	b, _ := io.ReadAll(c.Request.Body)
	q, _ := url.ParseQuery(string(b))
	contact_ids := q["selected_contact_ids"]
	for _, cid := range contact_ids {
		id, _ := strconv.Atoi(cid)
		app.contactsUseCase.Delete(id)
	}
	flashMessage(c, "Deleted Contacts!")
	contacts, _ := app.contactsUseCase.List()
	c.HTML(http.StatusOK, "index.html", gin.H{"contacts": contacts, "messages": flashes(c)})
}
