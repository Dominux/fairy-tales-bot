package main

import (
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	tele "gopkg.in/telebot.v3"

	"github.com/Dominux/fairy-tales-bot/internal/common"
	"github.com/Dominux/fairy-tales-bot/internal/entities"
	"github.com/Dominux/fairy-tales-bot/internal/services"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}

	btnAdd    = menu.Text("➕ Добавить сказку")
	btnCancel = menu.Text("❌ Отменить")
)

func main() {
	// getting db conn
	var db *sqlx.DB
	{
		var (
			dbUser = os.Getenv("POSTGRES_USER")
			dbPswd = os.Getenv("POSTGRES_PASSWORD")
			dbName = os.Getenv("POSTGRES_DB")
		)

		db = common.NewDB(dbUser, dbPswd, dbName)
	}

	// getting fairy tales service
	fairyTalesService := services.NewFairyTalesService(db)

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(c tele.Context) error {
		menu.Reply(
			menu.Row(btnAdd),
		)

		return c.Send("Здарова!", menu)
	})

	bot.Handle(&btnAdd, func(c tele.Context) error {
		// initing fairy tale creation
		msg_id := c.Message().ID
		fairyTalesService.InitCreating(msg_id)

		// replying with menu
		menu.Reply(
			menu.Row(btnCancel),
		)
		return c.Send("Как будет называться сказка?", menu)
	})

	bot.Handle(&btnCancel, func(c tele.Context) error {
		// doing so to send our answer ASAP and make our message id as just the next after user's one,
		// to prevent user to get the next message id faster than us
		// so at this moment we can send message that the fairy tail creation canceled
		// and to delete it too after a while (3 seconds, for example)
		defer func() {
			// deleting all the messages in creation interval
			ft, _ := fairyTalesService.GetUncompleted()

			current_msg_id := c.Message().ID
			chat_id := c.Chat().ID
			common.DeleteMessagesInterval(bot, chat_id, ft.Init_msg_id, current_msg_id)

			// initing fairy tale creation
			fairyTalesService.CancelCreation()
		}()

		// replying with menu
		menu.Reply(
			menu.Row(btnAdd),
		)
		return c.Send("Создание сказки отмененно", menu)
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		// getting fairy tale stage
		ft, err := fairyTalesService.GetUncompleted()
		if err != nil || ft.Stage != entities.Inited {
			return c.Delete()
		}

		// registering fairy tale name
		name := c.Message().Text
		fairyTalesService.RegisterName(name)

		return c.Send("Теперь запиши голосовое сообщение со сказкой")
	})

	bot.Handle(tele.OnVoice, func(c tele.Context) error {
		// getting fairy tale stage
		ft, err := fairyTalesService.GetUncompleted()
		if err != nil || ft.Stage != entities.Named {
			return c.Delete()
		}

		// getting file from telegram servers
		voiceFile := c.Message().Voice.File
		readCloser, err := bot.File(&voiceFile)
		if err != nil {
			log.Print(err)
		}

		file := tele.File{FileReader: readCloser}
		a := &tele.Audio{File: file, Title: *ft.Name, FileName: *ft.Name, Performer: "Fairy tales bot"}

		// sending audio
		c.Send(a)

		defer func() {
			// assuming the user didn't send anything after his message
			msg_id := c.Message().ID
			bot_msg_id := msg_id + 1

			// saving our audio message id
			fairyTalesService.RegisterAudio(bot_msg_id)

			// deleting unneeded messages
			chat_id := c.Chat().ID
			go common.DeleteMessagesInterval(bot, chat_id, ft.Init_msg_id, msg_id)
		}()

		// sending message about finishing fairy tale creation
		menu.Reply(
			menu.Row(btnAdd),
		)
		return c.Send("Сказка успешно сохранена", menu)
	})

	log.Print("Starting bot")
	bot.Start()
}
