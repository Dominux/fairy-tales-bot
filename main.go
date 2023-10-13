package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		crudMenu = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}

		btnAdd    = crudMenu.Text("➕ Добавить сказку")
		btnCancel = crudMenu.Text("X Отменить")
	)

	b.Handle("/start", func(c tele.Context) error {
		crudMenu.Reply(
			crudMenu.Row(btnAdd),
		)

		return c.Send("Hello!", crudMenu)
	})

	b.Handle(&btnAdd, func(c tele.Context) error {
		crudMenu.Reply(
			crudMenu.Row(btnCancel),
		)

		return c.Send("Как сказка будет называться?", crudMenu)
	})

	b.Handle(&btnCancel, func(c tele.Context) error {
		return c.Send("Создание сказки отмененно")
	})

	b.Start()
}
