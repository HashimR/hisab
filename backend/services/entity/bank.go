package entity

import (
	"main/models"
	"main/store"
	"main/store/dbmodels"
)

type BankService struct {
	bankRepository *store.BankRepository
}

func NewBankService(br *store.BankRepository) *BankService {
	return &BankService{
		bankRepository: br,
	}
}

func (bs *BankService) StoreBank(bank models.BankResponse) (int, error) {
	bankData := &dbmodels.BankData{
		Name:            bank.Name,
		LeanIdentifier:  bank.LeanBankIdentifier,
		Logo:            bank.Logo,
		MainColor:       bank.MainColor,
		BackgroundColor: bank.BackgroundColor,
	}

	// Call the UserRepository to create the user.
	id, err := bs.bankRepository.CreateBankIfDoesNotExist(bankData)
	if err != nil {
		return 0, err
	}

	return id, nil
}
