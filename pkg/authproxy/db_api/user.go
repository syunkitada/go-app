package db_api

import (
	"encoding/hex"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"

	"github.com/syunkitada/goapp/pkg/authproxy/db_model"
	"github.com/syunkitada/goapp/pkg/lib/error_utils"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (api *Api) CreateUser(tctx *logger.TraceContext, db *gorm.DB, name string, password string) (err error) {
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	tx := db.Begin()
	defer func() {
		if tmpErr := recover(); tmpErr != nil {
			err = error_utils.NewRecoveredError(tmpErr)
		}
		api.Rollback(tctx, tx, err)
	}()

	var user db_model.User
	if err = tx.Where("name = ?", name).First(&user).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		var hashedPassword string
		hashedPassword, err = api.generateHashFromPassword(password)
		if err != nil {
			return err
		}

		user = db_model.User{
			Name:     name,
			Password: hashedPassword,
		}
		tx.Create(&user)
		err = tx.Commit().Error
	}
	return err
}

func (api *Api) generateHashFromPassword(password string) (string, error) {
	converted, err := scrypt.Key([]byte(password), []byte(api.appConf.Auth.Secrets[0]), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(converted[:]), nil
}