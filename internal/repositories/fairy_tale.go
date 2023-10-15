package repositories

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Dominux/fairy-tales-bot/internal/entities"
)

type FairyTalesRepository struct {
	db *sqlx.DB
}

func NewFairyTalesRepo(db *sqlx.DB) FairyTalesRepository {
	return FairyTalesRepository{db}
}

func (repo FairyTalesRepository) Create(init_msg_id int) error {
	id := uuid.New()
	query := `INSERT INTO fairy_tales (id, init_msg_id, stage) VALUES ($1, $2, $3)`
	_, err := repo.db.Exec(query, id, init_msg_id, entities.Inited)
	return err
}

func (repo FairyTalesRepository) List() ([]entities.FairyTale, error) {
	query := `SELECT * FROM fairy_tales`
	fts := []entities.FairyTale{}
	err := repo.db.Select(&fts, query)
	return fts, err
}

func (repo FairyTalesRepository) GetUncompleted() (entities.FairyTale, error) {
	query := `SELECT * FROM fairy_tales WHERE NOT stage = $1`
	var ft entities.FairyTale
	err := repo.db.QueryRowx(query, entities.Created).StructScan(&ft)
	return ft, err
}

func (repo FairyTalesRepository) GetByID(id uuid.UUID) (entities.FairyTale, error) {
	query := `SELECT * FROM fairy_tales WHERE id = $1`
	var ft entities.FairyTale
	err := repo.db.QueryRowx(query, id).StructScan(&ft)
	return ft, err
}

func (repo FairyTalesRepository) RegisterName(name string) error {
	query := `UPDATE fairy_tales SET name = $1, stage = $2 WHERE NOT stage = $3`
	_, err := repo.db.Exec(query, name, entities.Named, entities.Created)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (repo FairyTalesRepository) RegisterAudio(audio_msg_id int) error {
	query := `UPDATE fairy_tales SET audio_msg_id = $1, stage = $2 WHERE NOT stage = $3`
	_, err := repo.db.Exec(query, audio_msg_id, entities.Created, entities.Created)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (repo FairyTalesRepository) DeleteUncompleted() error {
	query := `DELETE FROM fairy_tales WHERE NOT stage = $1`
	_, err := repo.db.Exec(query, entities.Created)
	return err
}
