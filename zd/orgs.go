package zd

import (
	"database/sql"

	"github.com/MEDIGO/go-zendesk/zendesk"

	"github.com/rs/zerolog/log"
)

type userCount struct {
	userId int64
	count  int64
}

func LookupOrgs(orgName string) (orgs []zendesk.Organization, err error) {
	log.Info().Str("orgName", orgName).Msg("Calling LookupOrgs")
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

	orgs, err = client.AutocompleteOrganizations(orgName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't get autocomplete orgs")
		return
	}

	return
}

// Add an organisation to out local DB
func AddOrg(org zendesk.Organization) (err error) {
	log.Info().Str("orgName", *org.Name).Msg("Calling AddOrg")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	_, err = db.Exec("INSERT INTO orgs VALUES( ?, ?) ON CONFLICT DO NOTHING;", *(org.ID), *(org.Name))
	if err != nil {
		log.Err(err).Int64("orgId", *(org.ID)).Str("orgNme", *(org.Name)).Msg("couldn't insert org")
		return
	}

	return
}

// Remove an organisation and its users from our local DB
func DeleteOrg(orgName string) error {
	log.Info().Str("orgName", orgName).Msg("Calling DeleteOrg")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return err
	}

	orgId, err := getOrgIdByName(orgName)
	if err != nil {
		return err
	}

	users, err := db.Query("WITH orgusers AS (SELECT DISTINCT user FROM orguser WHERE org = ?) SELECT user, count(1) count FROM orguser WHERE user IN (SELECT * FROM orgusers) GROUP BY  user", orgId)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't select org users")
		return err
	}
	defer users.Close()

	var userCounts []userCount
	for users.Next() {
		var uc userCount
		err = users.Scan(&uc.userId, &uc.count)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
		userCounts = append(userCounts, uc)
	}

	// delete the users
	for _, U := range userCounts {
		if U.count == 1 { // if the user is unique to this org then delete it
			_, err = db.Exec("DELETE FROM users WHERE id = ?;", U.userId)
			if err != nil {
				log.Err(err).Int64("userId", U.userId).Msg("couldn't delete user")
				return err
			}
		}

		_, err = db.Exec("DELETE FROM orguser WHERE user = ? AND org = ?;", U.userId, orgId)
		if err != nil {
			log.Err(err).Int64("userId", U.userId).Int64("orgId", orgId).Msg("couldn't delete orguser")
			return err
		}
	}

	// delete the org
	_, err = db.Exec("DELETE FROM orgs WHERE id = ?;", orgId)
	if err != nil {
		log.Err(err).Int64("orgId", orgId).Msg("couldn't delete org")
		return err
	}

	return nil
}

// Pull the names of orgs from our local DB
func GetOrgs() (orgs []string) {
	log.Info().Msg("Calling GetOrgs")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT name FROM orgs ORDER BY name ASC")
	if err != nil {
		log.Err(err).Msg("couldn't select org names")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var org string
		err = rows.Scan(&org)
		if err != nil {
			log.Err(err).Msg("couldn't scan rows")
			continue
		}
		orgs = append(orgs, org)
	}

	return
}

func GetOrgsLike(orgName string) (orgs []string) {
	log.Info().Str("orgName", orgName).Msg("Calling GetOrgsLike")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT name FROM orgs where name LIKE '%" + orgName + "%' ORDER BY name ASC")
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't select orgs")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var org string
		err = rows.Scan(&org)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
		orgs = append(orgs, org)
	}

	return
}

func getOrgIdByName(orgName string) (id int64, err error) {
	log.Info().Str("orgName", orgName).Msg("Calling getOrgIdByName")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT Id FROM orgs WHERE name = ?;", orgName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't select org by name")
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}
	return
}
