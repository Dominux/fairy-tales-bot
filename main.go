package main

import (
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	tele "gopkg.in/telebot.v3"

	"github.com/Dominux/fairy-tales-bot/internal/common"
	"github.com/Dominux/fairy-tales-bot/internal/handlers"
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

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	ftHandler := handlers.NewFairyTalesHandler(db, menu, &btnAdd, &btnCancel)

	bot.Handle("/start", ftHandler.OnStart)
	bot.Handle(&btnAdd, ftHandler.OnBtnAdd)
	bot.Handle(&btnCancel, ftHandler.OnBtnCancel)
	bot.Handle(tele.OnText, ftHandler.OnText)
	bot.Handle(tele.OnVoice, ftHandler.OnVoice)

	log.Print("Starting bot")
	bot.Start()
}
