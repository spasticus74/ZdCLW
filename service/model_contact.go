/*
 * ZdCLSync
 *
 * This server is used to simplify keeping instances of the Zendesk Call Logger application up to date with the configured assets.
 * It also acts as an online backup of the calls logged by users of the app, and a searchable reporitory of site contacts.
 * The Zendesk Call Logger application can be found [here](https://gitea.rantorium.com/MTS_ISC/zdcl).
 *
 * API version: 0.0.2
 * Contact: troy.campbell@global.komatsu
 */
package service

type Contact struct {
	// identifier
	Id int64 `json:"id"`
	// contact's name
	Name string `json:"name"`
	// contact's site
	Site string `json:"site"`
	// contact's company (employer)
	Company string `json:"company"`
	// contact's role on site
	Role string `json:"role"`
	// contact's email address
	Email string `json:"email"`
	// contact's primary phone number
	Phone1 string `json:"phone1"`
	// contact's secondaryphone number
	Phone2 string `json:"phone2"`
	// contextual notes about the contact
	Notes string `json:"notes,omitempty"`
	// organisation identifier
	OrgId int64 `json:"orgid"`
}
