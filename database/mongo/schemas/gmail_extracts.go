package schemas

import (
	"time"
)

type GmailExtract struct {
	ID       string    `bson:"_id,omitempty"`
	Email    string    `bson:"email"`
	FilePath string    `bson:"file_path"`
	Date     time.Time `bson:"date"`
}
