package entity

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"main/models"
	"main/store"
	"main/store/dbmodels"
	"strings"
)

type EntityService struct {
	entityRepository *store.EntityRepository
	bankService      *BankService
}

func NewEntityService(er *store.EntityRepository, bs *BankService) *EntityService {
	return &EntityService{
		entityRepository: er,
		bankService:      bs,
	}
}

func (es *EntityService) StoreEntity(entity models.EntityResponse, userId int) error {
	bankId, err := es.bankService.StoreBank(entity.Payload.BankDetails)
	if err != nil {
		return err
	}

	permissionsJson, err := json.Marshal(entity.Payload.Permissions)
	if err != nil {
		return err
	}

	// Add my db id
	entityData := &dbmodels.EntityData{
		LeanId:      entity.Payload.Id,
		LeanUserId:  entity.Payload.LeanCustomerId,
		BankId:      bankId,
		Permissions: string(permissionsJson),
		UserId:      userId,
	}

	if err := es.entityRepository.CreateEntity(entityData); err != nil {
		return err
	}

	return nil
}

func (es *EntityService) GetEntitiesForUser(userID int) ([]*models.Entity, error) {
	entitiesData, err := es.entityRepository.GetEntitiesForUser(userID)
	if err != nil {
		log.Error().Err(err).Msgf("Could not get entities for user id: %d", userID)
		return nil, err
	}

	entities := make([]*models.Entity, len(entitiesData))
	for i, entityData := range entitiesData {
		entities[i] = mapEntityDataToEntity(entityData)
	}

	return entities, nil
}

func mapEntityDataToEntity(entityData *dbmodels.EntityData) *models.Entity {
	entity := &models.Entity{
		Id:             entityData.ID,
		LeanId:         entityData.LeanId,
		LeanCustomerId: entityData.LeanUserId,                      // Assuming this field corresponds to LeanUserId
		Permissions:    strings.Split(entityData.Permissions, ","), // Assuming permissions are comma-separated strings
		BankId:         entityData.BankId,
	}
	return entity
}
