package schemas

type Message struct {
	ID             string `bson:"_id"`
	NotificationID string `bson:"notification_id"`
	Status         string `bson:"status"`
}
