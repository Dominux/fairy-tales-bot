package entities

import "fmt"

// StoredMessage is an example struct suitable for being
// stored in the database as-is or being embedded into
// a larger struct, which is often the case (you might
// want to store some metadata alongside, or might not.)
type StoredMessage struct {
	MessageID int   `sql:"message_id" json:"message_id"`
	ChatID    int64 `sql:"chat_id" json:"chat_id"`
}

func NewStoredMessage(msg_id int, chat_id int64) StoredMessage {
	return StoredMessage{MessageID: msg_id, ChatID: chat_id}
}

func (x StoredMessage) MessageSig() (string, int64) {
	return fmt.Sprint(x.MessageID), x.ChatID
}
