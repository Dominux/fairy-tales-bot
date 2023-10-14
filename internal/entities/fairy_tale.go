package entities

import "github.com/google/uuid"

type FairyTaleStage string

const (
	Inited  FairyTaleStage = "inited"
	Named   FairyTaleStage = "named"
	Created FairyTaleStage = "created"
)

type FairyTale struct {
	ID           uuid.UUID      `db:"id"`
	Name         *string        `db:"name"`
	Init_msg_id  int            `db:"init_msg_id"`
	Audio_msg_id *int           `db:"audio_msg_id"`
	Stage        FairyTaleStage `db:"stage"`
}
