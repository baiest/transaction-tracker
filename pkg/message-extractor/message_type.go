package messageextractor

type MessageType string

const (
	Extract  MessageType = "extract"
	Movement MessageType = "movement"
	Unknown  MessageType = "unknown"
)
