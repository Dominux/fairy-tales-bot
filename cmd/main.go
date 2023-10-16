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

	btnAdd     = menu.Text("➕ Добавить сказку")
	btnCancel  = menu.Text("❌ Отменить")
	btnGetList = menu.Text("📚 Список сказок")
	btnDelList = menu.Text("🗑 Список сказок на удаление")
)

func main() {
	var (
		db  *sqlx.DB
		bot *tele.Bot
	)

	// getting db conn
	{
		var (
			dbUser = os.Getenv("POSTGRES_USER")
			dbPswd = os.Getenv("POSTGRES_PASSWORD")
			dbName = os.Getenv("POSTGRES_DB")
		)

		db = common.NewDB(dbUser, dbPswd, dbName)
	}

	{

		timeoutStr := os.Getenv("LONG_POLLING_TIMEOUT")
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			log.Fatal(err)
			return
		}

		pref := tele.Settings{
			Token:  os.Getenv("TOKEN"),
			Poller: &tele.LongPoller{Timeout: time.Duration(timeout) * time.Second},
		}
		bot, err = tele.NewBot(pref)
		if err != nil {
			log.Fatal(err)
			return
		}
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

	ftHandler := handlers.NewFairyTalesHandler(db, menu, &btnAdd, &btnCancel, &btnGetList, &btnDelList)

	bot.Handle("/start", ftHandler.OnStart)
	bot.Handle(&btnAdd, ftHandler.OnBtnAdd)
	bot.Handle(&btnCancel, ftHandler.OnBtnCancel)
	bot.Handle(&btnGetList, ftHandler.OnGetList)
	bot.Handle(&btnDelList, ftHandler.OnDelList)
	bot.Handle(tele.OnText, ftHandler.OnText)
	bot.Handle(tele.OnVoice, ftHandler.OnVoice)

	log.Print("Starting bot")
	bot.Start()
}
