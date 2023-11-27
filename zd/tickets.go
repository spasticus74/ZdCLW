package zd

import (
	"fmt"
	"strings"

	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/rs/zerolog/log"
)

func CreateTicket(caller, email, phone, org, problem, issue, asset, product string) (ticketId string, err error) {
	log.Info().Str("caller", caller).Str("email", email).Str("phone", phone).Str("orgName", org).
		Str("problem", problem).Str("issue", issue).Str("asset", asset).Str("product", product).Msg("Calling CreateTicket")

	t := zendesk.Ticket{}

	// get the id of the caller
	callerId, err := getUserIdByNameAndOrg(caller, org)
	if err != nil {
		return
	}

	// get the id of the caller's organisation
	orgId, err := getOrgIdByName(org)
	if err != nil {
		return
	}

	// Append the caller's details to the problem description
	p := problem + fmt.Sprintf("\n\nCaller: %s\n Email: %s\n Phone: %s", caller, email, phone)

	// get any special fields that need to be set
	extraFields, err := getTicketFields()
	if err != nil {
		return "", err
	}

	for i, f := range extraFields.customfields {
		switch *(f.ID) {
		case 4926803905179:
			a, err2 := getAssetTag(asset)
			if err2 != nil {
				return "", err2
			}
			f.Value = a
		case 21002523:
			f.Value = strings.ToLower(product)

		}
		extraFields.customfields[i] = f
	}

	t.Subject = &issue
	t.Description = &p
	t.RequesterID = &callerId
	t.OrganizationID = &orgId
	t.TicketFormID = &extraFields.formId
	t.CustomFields = extraFields.customfields

	// get our credentials
	conxDetails, err := getZdCredentials()
	if err != nil {
		return
	}

	// create a ZD client
	client, err := zendesk.NewClient(conxDetails.instance, conxDetails.email, conxDetails.password)
	if err != nil {
		log.Err(err).Str("instance", conxDetails.instance).Str("email", conxDetails.email).Msg("couldn't create Zendesk client")
		return
	}

	newT, err := client.CreateTicket(&t)
	if err != nil {
		log.Err(err).Msg("couldn't create Zendesk Ticket")
		return
	}

	ticketId = fmt.Sprintf("https://%s.zendesk.com/agent/tickets/%d", conxDetails.instance, *newT.ID)

	return
}
