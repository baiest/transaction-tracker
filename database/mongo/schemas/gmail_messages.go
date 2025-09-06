package schemas

type Message struct {
	ID             string `bson:"_id"`
	NotificationID string `bson:"notification_id"`
	Status         string `bson:"status"`
}

func (m *Message) LogProperties() map[string]string {
	return map[string]string{
		"id":              m.ID,
		"notification_id": m.NotificationID,
		"status":          m.Status,
	}
}
