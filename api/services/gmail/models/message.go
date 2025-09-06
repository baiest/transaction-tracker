package models

type MessageType string

const (
	Unknown  MessageType = "unknown"
	Extract  MessageType = "extract"
	Movement MessageType = "movement"
)
