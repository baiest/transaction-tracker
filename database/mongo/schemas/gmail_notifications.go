package schemas

type GamilNotification struct {
	ID       string `bson:"_id"`
	Email    string `bson:"email" validate:"required"`
	Status   string
	Messages []*Message `bson:"-"`
}
