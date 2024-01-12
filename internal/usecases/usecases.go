package usecases

import (
	"time"

	"github.com/darrenparkinson/htmx-contact-app-go/internal/domain"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/ports"
)

type contactUseCase struct {
	contactRepo ports.ContactRepository
}

func NewContactUseCase(contactRepo ports.ContactRepository) ports.ContactUseCase {
	return &contactUseCase{
		contactRepo: contactRepo,
	}
}

func (c *contactUseCase) Create(first, last, phone, email string) error {
	return c.contactRepo.Create(domain.Contact{
		First: first,
		Last:  last,
		Phone: phone,
		Email: email,
	})
}

func (c *contactUseCase) List() ([]domain.Contact, error) {
	return c.contactRepo.List()
}

func (c *contactUseCase) Update(contact domain.Contact) error {
	return c.contactRepo.Update(contact)
}

func (c *contactUseCase) Delete(id int) error {
	return c.contactRepo.Delete(id)
}

func (c *contactUseCase) Search(query string) ([]domain.Contact, error) {
	return c.contactRepo.Search(query)
}

func (c *contactUseCase) Find(id int) (*domain.Contact, error) {
	return c.contactRepo.Find(id)
}

func (c *contactUseCase) Count() int {
	time.Sleep(2 * time.Second)
	return c.contactRepo.Count()
}
