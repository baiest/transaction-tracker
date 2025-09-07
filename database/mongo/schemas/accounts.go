package schemas

type Account struct {
	ID           string `bson:"_id"`
	Email        string `bson:"email"`
	RefreshToken string `bson:"refresh_token"`
}
