// Package app ties together application resources and handlers.
package app

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/UNWomenTeam/base-backend/database"
	"github.com/UNWomenTeam/base-backend/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
	ctxProfile
)

// API provides application resources and handlers.
type API struct {
	Account *AccountResource
	Profile *ProfileResource
}

// NewAPI configures and returns application API.
func NewAPI(db *pg.DB) (*API, error) {
	accountStore := database.NewAccountStore(db)
	account := NewAccountResource(accountStore)

	profileStore := database.NewProfileStore(db)
	profile := NewProfileResource(profileStore)

	api := &API{
		Account: account,
		Profile: profile,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/account", a.Account.router())
	r.Mount("/profile", a.Profile.router())

	return r
}

func log(r *http.Request) zap.Logger {
	return logging.GetLogEntry(r)
}
