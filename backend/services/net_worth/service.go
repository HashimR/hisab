package net_worth

import (
	"github.com/shopspring/decimal"
	"main/services/balances"
	"main/store"
	"time"
)

type Service struct {
	accountRepository  *store.AccountRepository
	balanceService     *balances.Service
	netWorthRepository *store.NetWorthRepository
}

const creditType = "CREDIT"

func NewNetWorthService(
	ar *store.AccountRepository,
	bs *balances.Service,
	nwr *store.NetWorthRepository) *Service {
	return &Service{
		accountRepository:  ar,
		balanceService:     bs,
		netWorthRepository: nwr,
	}
}

func (s *Service) CalculateCurrentNetWorth(userId int) (*NetWorth, error) {
	latestBalances, err := s.balanceService.GetLatestBalancesForUser(userId)
	if err != nil {
		return nil, err
	}

	//TODO: store new balance in accounts
	totalNetWorth := decimal.NewFromFloat(0.0)

	for _, balance := range latestBalances {
		if balance.Type != creditType {
			totalNetWorth = totalNetWorth.Add(balance.Amount)
		}
	}

	currentNetWorth := &NetWorth{
		Amount: totalNetWorth.String(),
		Time:   time.Now(),
	}

	err = s.netWorthRepository.StoreNetWorthForDate(userId, time.Now(), totalNetWorth, time.Now())
	if err != nil {
		return nil, err
	}

	return currentNetWorth, nil
}

func (s *Service) GetNetWorthForDate(userId int, date time.Time) (*NetWorth, error) {
	dailyNetWorth, err := s.netWorthRepository.GetNetWorthForDate(userId, date)
	if err != nil {
		return nil, err
	}

	// Convert the retrieved dailyNetWorth to the NetWorth model
	netWorth := &NetWorth{
		Amount: dailyNetWorth.NetWorth.String(), // Assuming NetWorth field in DailyNetWorth struct matches the type of Amount in NetWorth
		Time:   dailyNetWorth.Date,
	}

	return netWorth, nil
}

func (s *Service) GetNetWorthForDates(userId int, start time.Time, end time.Time) ([]*NetWorth, error) {
	dailyNetWorthList, err := s.netWorthRepository.GetNetWorthForDateRange(userId, start, end)
	if err != nil {
		return nil, err
	}

	netWorthList := make([]*NetWorth, len(dailyNetWorthList))
	for i, dailyNetWorth := range dailyNetWorthList {
		netWorthList[i] = &NetWorth{
			Amount: dailyNetWorth.NetWorth.String(),
			Time:   dailyNetWorth.Date,
		}
	}

	return netWorthList, nil
}

func (s *Service) GetLastXNetWorthRecords(userId int, x int) (*NetWorthGraph, error) {
	dailyNetWorthList, err := s.netWorthRepository.GetLastXNetWorths(userId, x)
	if err != nil {
		return nil, err
	}

	// Convert the retrieved dailyNetWorthList to a slice of NetWorth models
	netWorthList := make([]*NetWorth, len(dailyNetWorthList))
	for i, dailyNetWorth := range dailyNetWorthList {
		netWorthList[i] = &NetWorth{
			Amount: dailyNetWorth.NetWorth.String(),
			Time:   dailyNetWorth.Date,
		}
	}

	netWorthGraph := &NetWorthGraph{Points: netWorthList}

	return netWorthGraph, nil
}
