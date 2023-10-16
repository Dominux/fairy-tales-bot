package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	tele "gopkg.in/telebot.v3"

	"github.com/Dominux/fairy-tales-bot/internal/common"
	"github.com/Dominux/fairy-tales-bot/internal/entities"
	"github.com/Dominux/fairy-tales-bot/internal/services"
)

const (
	ftGetPrefix = "/get_ft"
	ftDelPrefix = "/del_ft"
)

type FairyTalesHandler struct {
	service *services.FairyTalesService

	menu       *tele.ReplyMarkup
	btnAdd     *tele.Btn
	btnCancel  *tele.Btn
	btnGetList *tele.Btn
	btnDelList *tele.Btn
}

func NewFairyTalesHandler(db *sqlx.DB, menu *tele.ReplyMarkup, btnAdd *tele.Btn, btnCancel *tele.Btn, btnGetList *tele.Btn, btnDelList *tele.Btn) FairyTalesHandler {
	service := services.NewFairyTalesService(db)
	return FairyTalesHandler{service: &service, menu: menu, btnAdd: btnAdd, btnCancel: btnCancel, btnGetList: btnGetList, btnDelList: btnDelList}
}

func (h *FairyTalesHandler) OnStart(c tele.Context) error {
	h.buildDefaultMenu()

	return c.Send("Здарова!", h.menu)
}

func (h *FairyTalesHandler) OnBtnAdd(c tele.Context) error {
	// initing fairy tale creation
	msg_id := c.Message().ID
	h.service.InitCreating(msg_id)

	// replying with menu
	h.menu.Reply(
		h.menu.Row(*h.btnCancel),
	)
	return c.Send("Как будет называться сказка?", h.menu)
}

func (h *FairyTalesHandler) OnBtnCancel(c tele.Context) error {
	// doing so to send our answer ASAP and make our message id as just the next after user's one,
	// to prevent user to get the next message id faster than us
	// so at this moment we can send message that the fairy tail creation canceled
	// and to delete it too after a while (3 seconds, for example)
	defer func() {
		// deleting all the messages in creation interval
		ft, _ := h.service.GetUncompleted()

		current_msg_id := c.Message().ID
		chat_id := c.Chat().ID
		common.DeleteMessagesInterval(c.Bot(), chat_id, ft.Init_msg_id, current_msg_id)

		// canceling fairy tale creation
		h.service.CancelCreation()
	}()

	// replying with menu
	h.buildDefaultMenu()
	return c.Send("Создание сказки отмененно", h.menu)
}

func (h *FairyTalesHandler) OnText(c tele.Context) error {
	msgText := c.Message().Text
	if strings.HasPrefix(msgText, ftGetPrefix) {
		return h.OnGet(c)
	} else if strings.HasPrefix(msgText, ftDelPrefix) {
		return h.OnDel(c)
	}

	// getting fairy tale stage
	ft, err := h.service.GetUncompleted()
	if err != nil || ft.Stage != entities.Inited {
		return c.Delete()
	}

	// registering fairy tale name
	h.service.RegisterName(msgText)

	return c.Send("Теперь запиши голосовое сообщение со сказкой")
}

func (h *FairyTalesHandler) OnVoice(c tele.Context) error {
	// getting fairy tale stage
	ft, err := h.service.GetUncompleted()
	if err != nil || ft.Stage != entities.Named {
		return c.Delete()
	}

	// getting file from telegram servers
	voiceFile := c.Message().Voice.File
	readCloser, err := c.Bot().File(&voiceFile)
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
		h.service.RegisterAudio(bot_msg_id)

		// deleting unneeded messages
		chat_id := c.Chat().ID
		go common.DeleteMessagesInterval(c.Bot(), chat_id, ft.Init_msg_id, msg_id)
	}()

	// sending message about finishing fairy tale creation
	h.buildDefaultMenu()
	return c.Send("Сказка успешно сохранена", h.menu)
}

func (h *FairyTalesHandler) OnGetList(c tele.Context) error {
	return h.onList(c, ftGetPrefix)
}

func (h *FairyTalesHandler) OnDelList(c tele.Context) error {
	return h.onList(c, ftDelPrefix)
}

func (h *FairyTalesHandler) onList(c tele.Context, symbol string) error {
	fts, _ := h.service.List()

	msg := ""
	for _, ft := range fts {
		cmd := formatCmd(ft.ID, symbol)
		line := fmt.Sprintf("• %s — %s\n", *ft.Name, cmd)
		msg += line
	}

	return c.Send(msg)
}

func (h *FairyTalesHandler) OnGet(c tele.Context) error {
	// getting and parsing fairy tale id
	msgText := c.Message().Text
	ftID, err := getIDFromGetCmd(msgText)
	if err != nil {
		return c.Delete()
	}

	// getting info from db
	ft, err := h.service.GetByID(ftID)
	if err != nil {
		return c.Delete()
	}

	// forwarding message with audio
	chatID := c.Chat().ID
	msg := entities.NewStoredMessage(*ft.Audio_msg_id, chatID)
	return c.Forward(msg)
}

func (h *FairyTalesHandler) OnDel(c tele.Context) error {
	// getting and parsing fairy tale id
	msgText := c.Message().Text
	ftID, err := getIDFromDelCmd(msgText)
	if err != nil {
		return c.Delete()
	}

	// getting from db
	ft, err := h.service.GetByID(ftID)
	if err != nil {
		return c.Delete()
	}

	defer func() {
		// deleting audio message
		chatID := c.Chat().ID
		msg := entities.NewStoredMessage(*ft.Audio_msg_id, chatID)
		c.Bot().Delete(msg)
	}()

	// deleting in db
	if err := h.service.Delete(ftID); err != nil {
		return c.Delete()
	}
	return c.Send("Сказка успешно удалена")
}

func (h *FairyTalesHandler) buildDefaultMenu() {
	h.menu.Reply(
		h.menu.Row(*h.btnAdd),
		h.menu.Row(*h.btnGetList),
		h.menu.Row(*h.btnDelList),
	)
}

////////////////////////////////////////////////////////////////
////	HELPERS
////////////////////////////////////////////////////////////////

func getIDFromGetCmd(cmd string) (uuid.UUID, error) {
	return getIDFromCmd(cmd, ftGetPrefix)
}

func getIDFromDelCmd(cmd string) (uuid.UUID, error) {
	return getIDFromCmd(cmd, ftDelPrefix)
}

func formatCmd(id uuid.UUID, symbol string) string {
	idStr := id.String()
	idStr = strings.ReplaceAll(idStr, "-", "_")
	return symbol + idStr
}

func getIDFromCmd(cmd string, symbol string) (uuid.UUID, error) {
	ftIDRaw := cmd[len(symbol):]
	ftIDRaw = strings.ReplaceAll(ftIDRaw, "_", "-")
	return uuid.Parse(ftIDRaw)
}
