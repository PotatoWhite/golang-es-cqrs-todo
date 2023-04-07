package command

import (
	"gorm.io/gorm"
	"time"
)

func NewTokenStore(db *gorm.DB) TokenStore {
	return &tokenStore{
		db: db,
	}
}

type TokenStore interface {
	CreateToken(targetGroup string, eventId uint) (token string, err error)
	GetLastTokenByTargetGroup(targetGroup string) (token *Token, err error)
}

type tokenStore struct {
	db *gorm.DB
}

func (t *tokenStore) GetLastTokenByTargetGroup(targetGroup string) (token *Token, err error) {
	return token, t.db.Where("target_group = ?", targetGroup).Last(&token).Error
}

func (t *tokenStore) CreateToken(targetGroup string, eventId uint) (token string, err error) {
	aToken := Token{
		TargetGroup: targetGroup,
		EventId:     eventId,
		CreatedAt:   time.Now(),
	}
	return token, t.db.Create(&aToken).Error
}
