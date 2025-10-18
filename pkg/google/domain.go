package google

import "golang.org/x/oauth2"

type GoogleAccount struct {
	Token           *oauth2.Token `bson:"token" validate:"required"`
	IsWatchingGmail bool          `bson:"is_watching_gmail"`
}
