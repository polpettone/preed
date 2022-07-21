package core

import (
	"github.com/polpettone/preed/cmd/adapter/persistence"
	"github.com/polpettone/preed/cmd/core/models"
)

type LedgerService struct {
	Repo persistence.Repo
}

func NewLedgerService(repo persistence.Repo) *LedgerService {
	return &LedgerService{Repo: repo}
}

func (service *LedgerService) SaveLedgerEntry(entry *models.LedgerEntry) error {
	return service.Repo.SaveLedgerEntry(entry)
}

func (service *LedgerService) DeleteLedgerEntry(entry *models.LedgerEntry) error {
	return service.Repo.DeleteLedgerEntry(entry)
}

func (service *LedgerService) FindLedgerEntryById(id string) (*models.LedgerEntry, error) {
	return service.Repo.FindLedgerEntryByID(id)
}

func (service *LedgerService) GetAllLedgerEntries() ([]models.LedgerEntry, error) {
	entries, err := service.Repo.FindAllLedgerEntries()
	if err != nil {
		return nil, err
	}
	return entries, nil
}
