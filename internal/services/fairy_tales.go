package services

import (
	"github.com/jmoiron/sqlx"

	"github.com/Dominux/fairy-tales-bot/internal/entities"
	"github.com/Dominux/fairy-tales-bot/internal/repositories"
)

type FairyTalesService struct {
	repo *repositories.FairyTalesRepository
}

func NewFairyTalesService(db *sqlx.DB) FairyTalesService {
	repo := repositories.NewFairyTalesRepo(db)
	return FairyTalesService{repo: &repo}
}

func (s *FairyTalesService) InitCreating(init_msg_id int) error {
	return s.repo.Create(init_msg_id)
}

func (s *FairyTalesService) GetUncompleted() (entities.FairyTale, error) {
	return s.repo.GetUncompleted()
}

func (s *FairyTalesService) RegisterName(name string) error {
	return s.repo.RegisterName(name)
}

func (s *FairyTalesService) RegisterAudio(audio_msg_id int) error {
	return s.repo.RegisterAudio(audio_msg_id)
}

func (s *FairyTalesService) CancelCreation() error {
	return s.repo.DeleteUncompleted()
}
