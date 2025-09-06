package schemas

import (
	"golang.org/x/oauth2"
)

type GoogleAccount struct {
	ID    string        `bson:"_id"`
	Token *oauth2.Token `bson:"token" validate:"required"`
}
