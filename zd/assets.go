package zd

import (
	"database/sql"
	"encoding/csv"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

// Read asset records from file and insert them into local db
func ImportAssets(path string) (err error) {
	log.Info().Str("path", path).Msg("Calling ImportAssets")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.Err(err).Str("path", path).Msg("couldn't open file")
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
			return err2
		}

		_, err = db.Exec("INSERT INTO assets VALUES(null, ?, ?, ?) ON CONFLICT DO NOTHING;", rec[1], rec[0], rec[2])
		if err != nil {
			log.Err(err).Str("name", rec[0]).Str("org", rec[1]).Str("tag", rec[2]).Msg("couldn't insert asset")
			return
		}
	}
	return
}

// Return the names of all assets defined for the org
func GetOrgAssets(orgName string) (assets []string, err error) {
	log.Info().Str("orgName", orgName).Msg("Calling GetOrgAssets")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	orgId, err := getOrgIdByName(orgName)
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT name FROM assets WHERE org = ? ORDER BY name ASC", orgId)
	if err != nil {
		log.Err(err).Str("orgName", orgName).Msg("couldn't select org assets")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var asset string
		err = rows.Scan(&asset)
		if err != nil {
			log.Err(err).Msg("error scanning row")
			return
		}

		assets = append(assets, asset)
	}

	return
}

// Return the tag string for the supplied asset name
func getAssetTag(assetName string) (assetTag string, err error) {
	log.Info().Str("assetName", assetName).Msg("Calling getAssetTag")
	db, err := sql.Open("sqlite3", "data/zd.db")
	if err != nil {
		log.Err(err).Msg("couldn't connect to database")
		return
	}

	rows, err := db.Query("SELECT tag FROM assets WHERE name = ? LIMIT 1", assetName)
	if err != nil {
		log.Err(err).Str("AssetName", assetName).Msg("couldn't select tag")
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&assetTag)
		if err != nil {
			log.Err(err).Msg("error scanning row")
			return
		}

	}

	return
}
