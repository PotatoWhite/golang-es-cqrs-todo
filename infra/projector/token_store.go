package projector

import (
	"github.com/potato/simple-restful-api/infra/command"
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
	GetLastTokenByTargetGroup(targetGroup string) (token *command.Token, err error)
}

type tokenStore struct {
	db *gorm.DB
}

func (t *tokenStore) GetLastTokenByTargetGroup(targetGroup string) (token *command.Token, err error) {
	return token, t.db.Where("target_group = ?", targetGroup).Last(&token).Error
}

func (t *tokenStore) CreateToken(targetGroup string, eventId uint) (token string, err error) {
	aToken := command.Token{
		TargetGroup: targetGroup,
		EventId:     eventId,
		CreatedAt:   time.Now(),
	}
	return token, t.db.Create(&aToken).Error
}
