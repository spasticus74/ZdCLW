package service

type Call struct {
	// identifier from the agent's local DB
	Id int64 `json:"id"`
	// seconds since epoch
	Time int64 `json:"time"`
	// name of the person that called
	Caller string `json:"caller"`
	// name of the organisation of which the caller is a member
	Org string `json:"org"`
	// extended description of the problem
	Problem string `json:"problem"`
	// URL of the Zendesk ticket
	Zendesk string `json:"zendesk"`
	// is this a new ticket?
	Newticket bool `json:"newticket"`
	// name of the agent that took the call
	Agent string `json:"agent"`
	// name of the asset that is the subject of the call
	Asset string `json:"asset,omitempty"`
	// brief description of the problem
	Issue string `json:"issue"`
	// agent's notes (not actually stored)
	Notes string `json:"notes,omitempty"`
	// the product or sub-system affected by the issue
	Product string `json:"product"`
	// phone number of caller
	Phone string `json:"phone,omitempty"`
	// reported time that issue began
	IssueStartDate string `json:"issuestartdate"`
}
