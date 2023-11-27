package zd

import (
	"database/sql"

	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/rs/zerolog/log"
)

// Return the name of all users for the specified org
func getUsersByOrgId(c zendesk.Client, orgId int64) (users []zendesk.User, err error) {
	log.Info().Int64("orgId", orgId).Msg("Calling getUsersByOrgId")
	z := zendesk.ListOptions{PerPage: 100}
	opts := zendesk.ListUsersOptions{}
	opts.ListOptions = z

	for {
		u, err2 := c.ListOrganizationUsers(orgId, &opts)
		if err2 != nil {
			log.Err(err2).Int64("orgId", orgId).Msg("couldn't get org users from Zendesk")
			return users, err2
		}
		users = append(users, u...)

		if len(u) < 100 {
			break
		} else {
			opts.Page += 1
		}
	}

	return
}

// Add active users for an organisation
func AddOrgUsers(orgId int64) (userCount int, err error) {
	log.Info().Int64("orgId", orgId).Msg("Calling AddOrgUsers")
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

	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	// get the new users
	users, err := getUsersByOrgId(client, orgId)

	// insert the active users into our local db
	for _, user := range users {
		if !*(user.Suspended) {
			ph := ""
			if user.Phone != nil {
				ph = *(user.Phone)
			}

			_, err = db.Exec("INSERT INTO users VALUES(?, ?, ?, ?) ON CONFLICT DO UPDATE SET name = ?, email = ?, phone = ?;", *(user.ID), *(user.Name), *(user.Email), ph, *(user.Name), *(user.Email), ph)
			if err != nil {
				log.Err(err).Int64("userId", *(user.ID)).Str("email", *(user.Email)).Str("phone", ph).Msg("couldn't insert user")
				return
			}
			_, err = db.Exec("INSERT INTO orguser VALUES(null, ?, ?) ON CONFLICT DO NOTHING;", *(user.ID), orgId)
			if err != nil {
				log.Err(err).Int64("userId", *(user.ID)).Int64("orgId", orgId).Msg("couldn't create orguser")
				return
			}
			userCount++
		}
	}

	return
}

// Return the name of all users associated with the org
func GetUserNamesByOrg(orgName string) (users []string) {
	log.Info().Str("orgName", orgName).Msg("Calling GetUserNamesByOrg")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT U.name FROM users U LEFT JOIN orguser OU ON U.ID = OU.user LEFT JOIN Orgs O ON O.id = OU.org WHERE O.name = ? ORDER BY U.name;", orgName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't select users by org name")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
		users = append(users, user)
	}

	return
}

// Return only the users for this org with a name like userName
func GetOrgUsersLike(orgName, userName string) (usernames []string) {
	log.Info().Str("orgName", orgName).Str("userName", userName).Msg("Calling GetOrgUsersLike")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT U.name FROM users U LEFT JOIN orguser OU ON U.ID = OU.user LEFT JOIN Orgs O ON O.id = OU.org WHERE O.name = ? AND u.name LIKE '%"+userName+"%' ORDER BY U.name;", orgName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Str("userName", userName).Msg("couldn't select org user")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			continue
		}
		usernames = append(usernames, username)
	}

	return
}

// Return the email address and phone number of the user matching the supplied name and org
func GetUserContactsByName(userName, orgName string) (email string, phone string, err error) {
	log.Info().Str("orgName", orgName).Str("userName", userName).Msg("Calling GetUserContactsByName")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT email, phone FROM users U LEFT JOIN orguser OU ON U.id = OU.user LEFT JOIN orgs O ON OU.org = O.Id WHERE O.name = ? AND U.name = ?;", orgName, userName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Str("userName", userName).Msg("couldn't select user contact details")
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&email, &phone)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}
	return
}

// Return the Id of the user matching the name and organisation
func getUserIdByNameAndOrg(userName, orgName string) (id int64, err error) {
	log.Info().Str("orgName", orgName).Str("userName", userName).Msg("Calling getUserIdByNameAndOrg")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT U.id FROM users U LEFT JOIN orguser OU ON U.id = OU.user LEFT JOIN orgs O ON OU.org = O.Id WHERE O.name = ? AND U.name = ?;", orgName, userName)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Str("userName", userName).Msg("couldn't select user")
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

// Update the user's phone field in our local database
func UpdateUserPhone(userName, phoneNumber string) (err error) {
	log.Info().Str("userName", userName).Str("phoneNumber", phoneNumber).Msg("Calling UpdateUserPhone")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT DISTINCT U.id FROM users U WHERE U.name = ?;", userName)
	if err != nil {
		log.Err(err).Str("userName", userName).Msg("couldn't select user")
		return
	}

	var userId int64
	rows.Next()
	err = rows.Scan(&userId)
	if err != nil {
		log.Err(err).Msg("couldn't scan row")
		return
	}
	rows.Close()

	_, err = db.Exec("UPDATE users SET phone = ? WHERE id = ?;", phoneNumber, userId)
	if err != nil {
		log.Err(err).Int64("userId", userId).Str("phoneNumber", phoneNumber).Msg("couldn't update user")
		return
	}

	return
}

// Update an organisation's users:
//  1. new users that don't exist:
//     1.1) for any org,
//     1.2) for only this org,
//  2. already exist but are now suspended
func RefreshOrgUsers(orgName string) (userCount int, err error) {
	log.Info().Str("orgName", orgName).Msg("Calling RefreshOrgUsers")
	orgId, err := getOrgIdByName(orgName)
	if err != nil {
		return
	}

	creds, err := getZdCredentials()
	if err != nil {
		return
	}

	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	// create a ZD client
	client, err := zendesk.NewClient(creds.instance, creds.email, creds.password)
	if err != nil {
		log.Err(err).Str("instance", creds.instance).Str("email", creds.email).Msg("couldn't create Zendesk client")
		return
	}

	orgUsers, err := getUsersByOrgId(client, orgId)
	if err != nil {
		return
	}

	for _, user := range orgUsers {
		// does this user exist in zd.users?
		cnt := 0
		rows, err2 := db.Query("SELECT COUNT(1) FROM users WHERE id = ?;", user.ID)
		if err2 != nil {
			log.Err(err2).Int64("userId", *user.ID).Msg("couldn't select user by id")
			return
		}
		defer rows.Close()

		rows.Next()
		err = rows.Scan(&cnt)
		if err != nil {
			log.Err(err).Msg("couldn't scan row")
			return
		}

		if cnt == 0 && !*(user.Suspended) { // 1.1 doesn't exist for any org
			ph := ""
			if user.Phone != nil {
				ph = *(user.Phone)
			}
			_, err = db.Exec("INSERT INTO users VALUES(?, ?, ?, ?);", *(user.ID), *(user.Name), *(user.Email), ph, *(user.Name), *(user.Email), ph)
			if err != nil {
				log.Err(err).Int64("userId", *(user.ID)).Str("userName", *(user.Name)).Str("email", *(user.Email)).Str("phoneNumber", ph).Msg("couldn't insert user")
				return
			}
			_, err = db.Exec("INSERT INTO orguser VALUES(null, ?, ?);", *(user.ID), orgId)
			if err != nil {
				log.Err(err).Int64("userId", *(user.ID)).Int64("orgId", orgId).Msg("couldn't insert orguser")
				return
			}
		} else {
			// does this user exist for this org?
			cnt := 0
			rows, err2 := db.Query("SELECT COUNT(1) FROM orguser WHERE user = ? AND org = ?;", user.ID, orgId)
			if err2 != nil {
				log.Err(err2).Int64("userId", *(user.ID)).Int64("orgId", orgId).Msg("couldn't select user count")
				return
			}
			defer rows.Close()

			rows.Next()
			err = rows.Scan(&cnt)
			if err != nil {
				log.Err(err).Msg("couldn't scan row")
				return
			}

			if cnt == 0 && !*(user.Suspended) { // 1.2 exists but not for this org
				_, err = db.Exec("INSERT INTO orguser VALUES(null, ?, ?);", *(user.ID), orgId)
				if err != nil {
					log.Err(err).Int64("userId", *(user.ID)).Int64("orgId", orgId).Msg("couldn't insert orguser")
					return
				}
			} else if cnt > 0 && *(user.Suspended) { // 2 exists but is suspended
				_, err = db.Exec("DELETE FROM orguser WHERE user = ? AND org = ?;", user.ID, orgId)
				if err != nil {
					log.Err(err).Int64("userId", *(user.ID)).Int64("orgId", orgId).Msg("couldn't delete orguser")
					return
				}
			}
		}
		userCount++
	}

	return
}
