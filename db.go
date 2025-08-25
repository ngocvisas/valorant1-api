package main

import (
	"context"

	"encore.dev/storage/sqldb"
)

var db = sqldb.Named("valorant1")

var _ = sqldb.Migration("001_create_tables", `
	CREATE TABLE IF NOT EXISTS loadouts (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		agent TEXT NOT NULL,
		primary_weapon TEXT,
		sidearm TEXT,
		created TIMESTAMP DEFAULT NOW()
	);
`)

func createLoadout(ctx context.Context, userID, agent, primary, sidearm string) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO loadouts (user_id, agent, primary_weapon, sidearm)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, agent, primary, sidearm).Scan(&id)
	return id, err
}

func listLoadouts(ctx context.Context, userID string) ([]Loadout, error) {
	rows, err := db.Query(ctx, `
		SELECT id, agent, primary_weapon, sidearm, created
		FROM loadouts
		WHERE user_id = $1
		ORDER BY created DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Loadout
	for rows.Next() {
		var l Loadout
		var primary, sidearm *string
		if err := rows.Scan(&l.ID, &l.Agent, &primary, &sidearm, &l.Created); err != nil {
			return nil, err
		}
		l.UserID = userID
		if primary != nil {
			l.Primary = *primary
		}
		if sidearm != nil {
			l.Sidearm = *sidearm
		}
		out = append(out, l)
	}
	return out, nil
}

func countLoadouts(ctx context.Context) (int, error) {
	var n int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM loadouts").Scan(&n)
	return n, err
}

