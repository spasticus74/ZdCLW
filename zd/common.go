package zd

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/MEDIGO/go-zendesk/zendesk"

	"github.com/rs/zerolog/log"
)

type ZDCREDS struct {
	email    string
	password string
	instance string
}

type TICKETFIELDS struct {
	formId       int64
	customfields []zendesk.CustomField
}

func getZdCredentials() (creds ZDCREDS, err error) {
	log.Info().Msg("Calling getZdCredentials")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	// get the required credentials for connecting to Zendesk
	rows, err := db.Query("select A.email, A.password, S.zdinstance from agent A left join settings S ON A.id = S.id AND A.id = 1")
	if err != nil {
		log.Err(err).Msg("couldn't select credentials")
		return
	}
	defer rows.Close()

	// there should only ever be one row
	for rows.Next() {
		err = rows.Scan(&(creds.email), &(creds.password), &(creds.instance))
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
	}

	pw, err := decrypt([]byte(creds.password))
	if err != nil {
		return
	}

	creds.password = string(pw)

	return
}

func getTicketFields() (fields TICKETFIELDS, err error) {
	log.Info().Msg("Callnig getTicketFields")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT S.formId FROM settings S WHERE S.id = 1")
	if err != nil {
		log.Err(err).Msg("couldn't select form id")
		return
	}
	defer rows.Close()

	// there should only ever be one row
	for rows.Next() {
		err = rows.Scan(&(fields.formId))
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
	}

	rows2, err := db.Query("SELECT C.id, C.value FROM customfields C")
	if err != nil {
		log.Err(err).Msg("couldn't select custom fields")
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		cf := zendesk.CustomField{}
		err = rows2.Scan(&(cf.ID), &(cf.Value))
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
		fields.customfields = append(fields.customfields, cf)
	}

	return
}

func FormatZendeskLink(zdLink string) (string, error) {
	log.Info().Str("zdLink", zdLink).Msg("Calling FormatZendeskLink")
	if len(zdLink) == 0 {
		return zdLink, nil
	}

	reZendeskLink := regexp.MustCompile(`^https:\/\/(.*).zendesk.com\/agent\/tickets\/(\d*)$`)
	reZendeskId := regexp.MustCompile(`^(\d{6,})$`)

	m := reZendeskLink.FindAllStringSubmatch(zdLink, -1)
	if m != nil { // match full link
		return zdLink, nil
	}

	m = reZendeskId.FindAllStringSubmatch(zdLink, -1)
	if m != nil { // match just an id
		creds, err := getZdCredentials()
		if err != nil {
			return zdLink, err
		}
		return fmt.Sprintf("https://%s.zendesk.com/agent/tickets/%s", creds.instance, m[0][1]), nil
	}

	err := errors.New("invalid link")
	log.Err(err).Str("zdLink", zdLink).Msg("could't parse zendesk link")

	return zdLink, err
}

func StoreAgentCredentials(email, pw string) (err error) {
	log.Info().Str("email", email).Msg("Calling StoreAgentCredentials")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	enc, err := encrypt([]byte(pw))
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO agent VALUES(1, ?, ?, unixepoch()) ON CONFLICT DO UPDATE SET email = ?, password = ?, lastupdate = unixepoch() WHERE id = 1;", email, enc, email, enc)
	if err != nil {
		log.Err(err).Str("email", email).Msg("couldn't insert agent credentials")
		return
	}
	return
}

func StoreInstanceName(instance string) (err error) {
	log.Info().Str("instance", instance).Msg("Calling storeInstanceName")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't open database")
		return
	}

	_, err = db.Exec("UPDATE settings SET zdinstance = ? WHERE id = 1;", instance)
	if err != nil {
		log.Err(err).Str("instance", instance).Msg("couldn't insert instance name")
		return
	}

	return
}

// Returns true if we have a user defined
func HaveAgent() (hasAgent bool, err error) {
	log.Info().Msg("Calling HaveAgent")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't open database")
		return
	}

	// first see if there's any record at all
	rows, err := db.Query("SELECT COUNT(1) FROM Agent")
	if err != nil {
		log.Err(err).Msg("couldn't select agent count")
		return
	}
	defer rows.Close()

	var agentCount int
	rows.Next()
	err = rows.Scan(&agentCount)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}
	if agentCount < 1 {
		log.Info().Msg("No agent defined")
		return false, errors.New("no agent defined")
	}

	// if a record exists make sure we've stored a password
	rows2, err := db.Query("SELECT length(password) FROM Agent")
	if err != nil {
		log.Err(err).Msg("couldn't select agent count")
		return
	}
	defer rows2.Close()

	var pwLength int
	rows2.Next()
	err = rows2.Scan(&pwLength)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}
	if pwLength > 0 {
		hasAgent = true
	}

	return
}
