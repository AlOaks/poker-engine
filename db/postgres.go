package db

import (
	"fmt"
	"log"
	"poker-engine/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DBConnection(url string) (*sqlx.DB, error) {
	log.Printf("Starting DB Connection...")
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("Error opening DB connection => %s", err.Error())

	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging DB => %s", err.Error())
	}

	log.Printf("Connected to DB")
	return db, nil
}

func GetGame(db *sqlx.DB, gameID int) (*models.GameState, error) {
	game := models.GameState{}

	err := db.Get(&game, "SELECT * FROM games WHERE id=$1", gameID)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func GetPlayersByGameID(db *sqlx.DB, gameID int) (*[]models.Player, error) {
	var players []models.Player
	err := db.Select(&players, "SELECT * FROM players WHERE game_id=$1", gameID)
	if err != nil {
		return nil, err
	}

	return &players, nil
}

func GetGames(db *sqlx.DB) (*[]models.GameState, error) {
	var games []models.GameState
	err := db.Select(&games, "SELECT * FROM games WHERE done=false")
	if err != nil {
		return nil, err
	}

	return &games, nil
}

func SaveGame(db *sqlx.DB, game models.GameState) (*int, error) {
	tx := db.MustBegin()
	defer tx.Rollback()

	row := tx.QueryRow("INSERT INTO games (street, done, winners, deck, board) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		game.Street,
		game.Done,
		game.Winners,
		game.RemainingDeck,
		game.Board,
	)

	var gameId int
	err := row.Scan(&gameId)
	if err != nil {
		return nil, err
	}

	for _, player := range game.Players {
		_, err = tx.Exec("INSERT into players (hole_cards, folded, game_id, player_name) VALUES ($1, $2, $3, $4)",
			player.HoleCards,
			player.Folded,
			gameId,
			player.PlayerName,
		)

		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()

	return &gameId, err
}

func UpdateGame(db *sqlx.DB, game models.GameState) error {
	tx := db.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec("UPDATE games SET street=$1, board=$2, deck=$3, done=$4, winners=$5 WHERE id=$6",
		game.Street,
		game.Board,
		game.RemainingDeck,
		game.Done,
		game.Winners,
		game.ID,
	)
	if err != nil {
		return err
	}

	for _, player := range game.Players {
		_, err := tx.Exec("UPDATE players SET folded=$1, hole_cards=$2 WHERE id=$3",
			player.Folded,
			player.HoleCards,
			player.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
