package user

import "time"

type Account struct {
	Name           string
	IsAdmin        bool
	Forename       string
	Surname        string
	HashedPassword string
	SessionExpires time.Time
}
