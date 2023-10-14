package common

import (
	tele "gopkg.in/telebot.v3"

	"github.com/Dominux/fairy-tales-bot/internal/entities"
)

// / Deletes all the messages in interval from start to end inclusively
func DeleteMessagesInterval(bot *tele.Bot, chat_id int64, start_msg_id int, end_msg_id int) {
	// skipping if start message id is bigger than end one
	if start_msg_id > end_msg_id {
		return
	}

	current_id := start_msg_id

	for current_id <= end_msg_id {
		msg := entities.NewStoredMessage(current_id, chat_id)
		bot.Delete(&msg)

		current_id++
	}
}
