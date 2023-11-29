package main

import (
	"context"
	"log"

	"github.com/spasticus74/ZdCLW/service"
	"github.com/spasticus74/ZdCLW/zd"
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
	}

	return contacts
}

// return a list of organisation names
func (a *App) GetOrgNames() []string {
	return zd.GetOrgs()
}

// get all defined callers
func (a *App) GetAllCallers() []string {
	c := service.NewConfiguration()
	s := service.NewAPIClient(c)
	callers, _, err := s.CallerApi.GetAllCallers(context.TODO())
	if err != nil {
		log.Println(err)
	}

	clrs := make([]string, 0)
	for _, caller := range callers {
		clrs = append(clrs, caller.Name)
	}
	return clrs
}

// return all callers from an organisation
func (a *App) GetCallersByOrg(orgName string) []string {
	c := service.NewConfiguration()
	s := service.NewAPIClient(c)
	callers, _, err := s.CallerApi.GetCallersByOrg(context.TODO(), orgName)
	if err != nil {
		log.Println(err)
	}

	clrs := make([]string, 0)
	for _, caller := range callers {
		clrs = append(clrs, caller.Name)
	}
	return clrs
}

// Create a new zendesk ticket from a call
func (a *App) CreateTicket(caller, email, phone, org, problem, issue, asset, product string) string {
	ticketId, err := zd.CreateTicket(caller, email, phone, org, problem, issue, asset, product)
	if err != nil {
		log.Println(err)
		ticketId = "Error creating ticket"
	}
	return ticketId
}

func (a *App) StoreAgentCredentials(email, password string) string {
	err := zd.StoreAgentCredentials(email, password)
	if err != nil {
		return "error"
	}
	return ""
}
