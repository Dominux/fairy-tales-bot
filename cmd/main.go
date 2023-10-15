package main

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"

	"github.com/Dominux/fairy-tales-bot/internal/common"
	"github.com/Dominux/fairy-tales-bot/internal/handlers"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}

	btnAdd    = menu.Text("➕ Добавить сказку")
	btnCancel = menu.Text("❌ Отменить")
	btnList   = menu.Text("📚 Список сказок")
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

	// allowing only 1 user to use our bot
	{
		allowedUserStr := os.Getenv("ALLOWED_USER_ID")
		allowedUser, err := strconv.ParseInt(allowedUserStr, 10, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		bot.Use(middleware.Whitelist(allowedUser))
	}

	ftHandler := handlers.NewFairyTalesHandler(db, menu, &btnAdd, &btnCancel, &btnList)

	bot.Handle("/start", ftHandler.OnStart)
	bot.Handle(&btnAdd, ftHandler.OnBtnAdd)
	bot.Handle(&btnCancel, ftHandler.OnBtnCancel)
	bot.Handle(&btnList, ftHandler.OnList)
	bot.Handle(tele.OnText, ftHandler.OnText)
	bot.Handle(tele.OnVoice, ftHandler.OnVoice)

	log.Print("Starting bot")
	bot.Start()
}
