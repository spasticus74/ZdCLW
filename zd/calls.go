package zd

import (
	"database/sql"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type CALL struct {
	Id         int
	Timeofcall string
	Caller     string
	Org        string
	Problem    string
	ZendeskURL string
	Email      string
	Phone      string
	Asset      string
	IssueStart string
	Issue      string
	Notes      string
	Product    string
}

// Record a call into the local database
func SubmitCall(caller, email, orgName, problem, issue, zendeskURL, product, asset, phone, issueStart string, isNewTicket bool) (callId int, err error) {
	log.Info().Str("caller", caller).Str("email", email).Str("orgName", orgName).Str("problem", problem).
		Str("issue", issue).Str("zendeskURL", zendeskURL).Str("product", product).Str("asset", asset).Str("issueStart", issueStart).Bool("isNew", isNewTicket).Msg("Calling SubmitCall")

	// get our credentials
	conxDetails, err := getZdCredentials()
	if err != nil {
		return
	}

	// add it to the local DB
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	_, err = db.Exec("INSERT INTO calls VALUES(null, unixepoch('now', 'localtime'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING;", caller, orgName, problem, zendeskURL, isNewTicket, conxDetails.email, asset, issue, "", product, phone, issueStart)
	if err != nil {
		log.Err(err).Str("caller", caller).Str("orgName", orgName).Str("problem", problem).Str("zendeskURL", zendeskURL).
			Bool("isNew", isNewTicket).Str("email", email).Str("asset", asset).Str("issue", issue).Str("product", product).Str("phone", phone).Str("issueStart", issueStart).Msg("couldn't insert call")
		return
	}

	// get the id of the record just added
	rows, err := db.Query("SELECT last_insert_rowid();")
	if err != nil {
		log.Err(err).Msg("couldn't get last_insert_rowid")
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&callId)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}

	return
}

// Write all call records to file
func ExportCalls(path string) (err error) {
	log.Info().Str("path", path).Msg("Calling ExportCalls")
	csvFile, err := os.Create(path)
	if err != nil {
		log.Err(err).Str("path", path).Msg("couldn't create file")
		return
	}
	defer csvFile.Close()

	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT id, timeofcall, caller, org, problem, zendesk, newticket, agent, asset, issue, notes, product, phone, issuestartdate FROM calls ORDER BY id ASC")
	if err != nil {
		log.Err(err).Msg("couldn't select calls")
		return
	}
	defer rows.Close()

	csvwriter := csv.NewWriter(csvFile)
	err = csvwriter.Write([]string{"id", "timestamp", "caller", "phone", "org", "problem", "zendesk", "isNew", "agent", "asset", "issue", "notes", "product", "issueStartDate"})
	if err != nil {
		log.Err(err).Str("path", csvFile.Name()).Msg("couldn't write header to file")
		return
	}

	defer csvwriter.Flush()
	for rows.Next() {
		var id string
		var timestamp int
		var caller string
		var org string
		var issue string
		var desc string
		var link string
		var isNew string
		var agent string
		var asset string
		var notes string
		var product string
		var phone string
		var issueStart string
		err = rows.Scan(&id, &timestamp, &caller, &org, &desc, &link, &isNew, &agent, &asset, &issue, &notes, &product, &phone, &issueStart)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			return
		}

		err = csvwriter.Write([]string{id, strconv.Itoa(timestamp), caller, phone, org, desc, link, isNew, agent, asset, issue, notes, product, issueStart})
		if err != nil {
			log.Err(err).Str("path", csvFile.Name()).Str("id", id).Int("timestam", timestamp).Str("caller", caller).Str("org", org).Str("desc", desc).
				Str("link", link).Str("isNew", isNew).Str("agent", agent).Str("issue", issue).Str("notes", notes).Str("product", product).Str("phone", phone).Str("issueStart", issueStart).Msg("couldn't write row to file")
			return
		}
	}

	return
}

// Read call records from file and insert them into local db
func ImportCalls(path string) (err error) {
	log.Info().Str("path", path).Msg("Calling ImportCalls")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.Err(err).Str("path", path).Msg("couldn't open calls file")
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Read() // header
	for {
		rec, err2 := r.Read()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			log.Err(err).Str("path", path).Msg("couldn't read line in calls file")
			return err2
		}

		_, err = db.Exec("INSERT INTO calls VALUES(null, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING;", rec[1], rec[2], rec[4], rec[5], rec[6], rec[7], rec[8], rec[9], rec[10], rec[11], rec[12], rec[3], rec[13])
		if err != nil {
			log.Err(err).Str("timeOfCall", rec[1]).Str("caller", rec[2]).Str("phone", rec[3]).Str("orgName", rec[4]).Str("problem", rec[5]).Str("zendeskLink", rec[6]).Str("isNew", rec[7]).
				Str("agent", rec[8]).Str("asset", rec[9]).Str("issue", rec[10]).Str("notes", rec[11]).Str("product", rec[12]).Str("issueStart", rec[13]).Msg("couldn't insert call")
			return
		}
	}
	return
}

// Return a previous recorded call by its Id
func GetCallById(callId int) (call CALL, err error) {
	log.Info().Int("callId", callId).Msg("Calling GetCallById")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT C.id, datetime(C.timeofcall, 'unixepoch') timeofcall, C.caller, C.org, C.problem, C.zendesk, U.email, U.phone, C.asset, C.issue, C.notes, C.product, C.issuestartdate FROM calls C LEFT JOIN users U ON U.Name = C.caller WHERE C.id = ?", callId)
	if err != nil {
		log.Err(err).Int("callId", callId).Msg("couldn't select call by id")
		return
	}
	defer rows.Close()

	// there should only ever be one row
	for rows.Next() {
		err = rows.Scan(&(call.Id), &(call.Timeofcall), &(call.Caller), &(call.Org), &(call.Problem), &(call.ZendeskURL), &(call.Email), &(call.Phone), &(call.Asset), &(call.Issue), &(call.Notes), &(call.Product), &(call.IssueStart))
		if err != nil {
			log.Err(err).Msg("error scanning row")
			return
		}
	}

	return
}

// Update the local notes field of a call
func UpdateNotes(callId int, notes string) (err error) {
	log.Info().Int("callId", callId).Str("notes", notes).Msg("Calling UpdateNotes")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}
	_, err = db.Exec("UPDATE calls SET notes = ? WHERE id = ?;", notes, callId)
	if err != nil {
		log.Err(err).Int("callId", callId).Str("notes", notes).Msg("couldn't update notes")
		return
	}
	return
}

// Return all defined products
func GetProducts() (products []string, err error) {
	log.Info().Msg("Calling GetProducts")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT name FROM products ORDER BY name ASC;")
	if err != nil {
		log.Err(err).Msg("couldn't select products")
		return
	}
	defer rows.Close()

	var p string
	for rows.Next() {
		err = rows.Scan(&p)
		if err != nil {
			log.Err(err).Msg("error scanning row")
			return
		}
		products = append(products, p)
	}

	return
}

// Return all calls
func GetCalls() (calls []CALL, err error) {
	log.Info().Msg("Calling GetCalls")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT C.id, datetime(C.timeofcall, 'unixepoch') timeofcall, C.caller, C.org, C.problem, C.zendesk, U.email, U.phone, C.asset, C.issue, C.notes, C.product, C.issuestartdate FROM calls C LEFT JOIN users U ON U.Name = C.caller ORDER BY C.id DESC")
	if err != nil {
		log.Err(err).Msg("couldn't select calls")
		return
	}
	defer rows.Close()

	for rows.Next() {
		call := CALL{}
		err = rows.Scan(&(call.Id), &(call.Timeofcall), &(call.Caller), &(call.Org), &(call.Problem), &(call.ZendeskURL), &(call.Email), &(call.Phone), &(call.Asset), &(call.Issue), &(call.Notes), &(call.Product), &(call.IssueStart))
		if err != nil {
			log.Err(err).Msg("error scanning row")
			continue
		}
		calls = append(calls, call)
	}
	return
}
