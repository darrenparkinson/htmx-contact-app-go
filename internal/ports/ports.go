package ports

import "github.com/darrenparkinson/htmx-contact-app-go/internal/domain"

type ContactUseCase interface {
	Create(first, last, phone, email string) error
	List() ([]domain.Contact, error)
	Update(first, last, phone, email string) error
	Delete(id int) error
	Search(query string) ([]domain.Contact, error)
	Find(id int) (*domain.Contact, error)
	Count() int
}

type ContactRepository interface {
	Create(contact domain.Contact) error
	List() ([]domain.Contact, error)
	Update(first, last, phone, email string) error
	Delete(id int) error
	Search(query string) ([]domain.Contact, error)
	Find(id int) (*domain.Contact, error)
	Count() int
}
