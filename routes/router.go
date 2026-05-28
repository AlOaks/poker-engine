package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"poker-engine/db"
	"poker-engine/engine"
	"poker-engine/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func NewAppRouter(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handler := NewHandler(db)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	router.Post("/games/new", handler.CreateNewGame)
	router.Post("/games/{id}/next", handler.NextGameStep)
	router.Post("/games/{id}/fold", handler.FoldStage)
	router.Get("/games/all", handler.GetGames)

	return router
}

func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := db.GetGames(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(*games)
}

func (h *Handler) CreateNewGame(w http.ResponseWriter, r *http.Request) {
	newGame, err := engine.NewGame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	equities := h.runSimulator(newGame)

	gameID, err := db.SaveGame(h.DB, *newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGame.ID = *gameID

	w.Header().Set("Content-Type", "application/json")

	response := EquityPredictionResponse{
		GameState: *newGame,
		Equities:  equities,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) NextGameStep(w http.ResponseWriter, r *http.Request) {
	existingGame, err := h.retrieveGame(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingGame = engine.NextGameStep(*existingGame)
	var equities []models.EquityPrediction

	if existingGame.Street != models.River {
		equities = h.runSimulator(existingGame)
	} else {
		for _, player := range existingGame.Players {
			isWinner := false
			for _, winner := range existingGame.Winners {
				if player.ID == winner.PlayerID {
					isWinner = true
					break
				}
			}

			if isWinner {
				equities = append(equities, models.EquityPrediction{ID: player.ID, Equity: 100.0})
				continue
			}

			equities = append(equities, models.EquityPrediction{ID: player.ID, Equity: 0.0})
		}
	}

	err = db.UpdateGame(h.DB, *existingGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := EquityPredictionResponse{
		GameState: *existingGame,
		Equities:  equities,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) FoldStage(w http.ResponseWriter, r *http.Request) {
	existingGame, err := h.retrieveGame(r)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingGame, err = engine.Folder(existingGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newEquities := h.runSimulator(existingGame)

	err = db.UpdateGame(h.DB, *existingGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := EquityPredictionResponse{
		GameState: *existingGame,
		Equities:  newEquities,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) runSimulator(game *models.GameState) []models.EquityPrediction {
	simulator := engine.NewSimulator()
	simulator.RunPredictions(*game)
	simChannel := simulator.GetChannel()

	equities := []models.EquityPrediction{}

	for event := range simChannel {
		equities = append(equities, *event)

		for i := range game.Players {
			if game.Players[i].ID == event.ID {
				game.Players[i].Equity = event.Equity
			}
		}
	}

	return equities
}

func (h *Handler) retrieveGame(r *http.Request) (*models.GameState, error) {
	gameId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, err
	}

	existingGame, err := db.GetGame(h.DB, gameId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, err
	}

	return existingGame, nil
}
