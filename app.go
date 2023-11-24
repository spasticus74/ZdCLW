package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spasticus74/ZdCLW/service"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// search for a contact by name
func (a *App) SearchContactByName(searchTerm string) string {
	c := service.NewConfiguration()
	s := service.NewAPIClient(c)
	contacts, _, err := s.ContactApi.GetContactsByName(context.TODO(), searchTerm)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Print(contacts)
	}
	return contacts
}

// search for a contact by name
func (a *App) SearchContactByOrg(searchTerm string) string {
	c := service.NewConfiguration()
	s := service.NewAPIClient(c)
	contacts, _, err := s.ContactApi.GetContactsByOrg(context.TODO(), searchTerm)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Print(contacts)
	}
	return contacts
}
