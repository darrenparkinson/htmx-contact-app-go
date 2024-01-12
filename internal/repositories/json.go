package repositories

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/darrenparkinson/htmx-contact-app-go/internal/cache"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/domain"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/helpers"
	"github.com/darrenparkinson/htmx-contact-app-go/internal/ports"
)

type contactJSONRepository struct {
	contactCache *cache.Cache[domain.Contact, int]
}

func NewContactJSONRepository() ports.ContactRepository {
	var contacts []domain.Contact
	err := helpers.LoadJSON("contacts.json", &contacts)
	if err != nil {
		log.Fatal(err)
	}
	contactsMap := make(map[int]domain.Contact)
	for _, c := range contacts {
		contactsMap[c.ID] = c
	}
	contactCache := cache.NewCache(contacts, contactsMap)
	return &contactJSONRepository{
		contactCache: contactCache,
	}
}

func (r *contactJSONRepository) Create(contact domain.Contact) error {
	err := r.Validate(contact)
	if err != nil {
		return err
	}
	var ids []int
	contacts := r.contactCache.Retrieve()
	for _, c := range contacts {
		ids = append(ids, c.ID)
	}
	contactsMap := r.contactCache.RetrieveMap()
	maxid := 0
	if len(contacts) != 0 {
		max := slices.Max(ids)
		fmt.Println("MAX:", max)
		maxid = max + 1

	}
	contact.ID = maxid
	contacts = append(contacts, contact)
	contactsMap[maxid] = contact
	r.contactCache.Update(contacts, contactsMap)
	helpers.SaveJSON("contacts.json", contacts)
	return nil
}

func (r *contactJSONRepository) List() ([]domain.Contact, error) {
	return r.contactCache.Retrieve(), nil
}

func (r *contactJSONRepository) Update(contact domain.Contact) error {
	err := r.Validate(contact)
	if err != nil {
		return err
	}
	contacts := r.contactCache.Retrieve()
	var contactPos int
	for i, c := range contacts {
		if c.ID == contact.ID {
			contactPos = i
		}
	}
	contacts[contactPos] = contact
	contactsMap := r.contactCache.RetrieveMap()
	contactsMap[contact.ID] = contact
	r.contactCache.Update(contacts, contactsMap)
	helpers.SaveJSON("contacts.json", contacts)
	return nil
}

func (r *contactJSONRepository) Delete(id int) error {
	return nil
}

func (r *contactJSONRepository) Search(query string) ([]domain.Contact, error) {
	var results []domain.Contact
	contacts := r.contactCache.Retrieve()
	for _, c := range contacts {
		matchFirst := strings.Contains(strings.ToLower(c.First), strings.ToLower(query))
		matchLast := strings.Contains(strings.ToLower(c.Last), strings.ToLower(query))
		matchEmail := strings.Contains(strings.ToLower(c.Email), strings.ToLower(query))
		matchPhone := strings.Contains(strings.ToLower(c.Phone), strings.ToLower(query))
		if matchFirst || matchLast || matchEmail || matchPhone {
			results = append(results, c)
		}
	}
	return results, nil
}

func (r *contactJSONRepository) Find(id int) (*domain.Contact, error) {
	contact, _ := r.contactCache.RetrieveMapEntry(id)
	return &contact, nil
}

func (r *contactJSONRepository) Count() int {
	return len(r.contactCache.Retrieve())
}

func (r *contactJSONRepository) Validate(contact domain.Contact) error {
	if contact.Email == "" {
		return domain.ErrMissingEmail
	}
	contacts := r.contactCache.Retrieve()
	for _, c := range contacts {
		if strings.EqualFold(c.Email, contact.Email) && c.ID != contact.ID {
			return domain.ErrExistingContact
		}
	}
	return nil
}
