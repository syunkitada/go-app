package base_db_api

import (
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/base/base_db_model"
	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (api *Api) CreateOrUpdateService(tctx *logger.TraceContext, db *gorm.DB, input *base_spec.UpdateService) (err error) {
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	err = api.Transact(tctx, db, func(tx *gorm.DB) (err error) {
		var service base_db_model.Service
		if err = tx.Where("name = ?", input.Name).First(&service).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return
			}

			service = base_db_model.Service{
				Name:      input.Name,
				Scope:     input.Scope,
				Endpoints: strings.Join(input.Endpoints, ","),
			}
			if err = tx.Create(&service).Error; err != nil {
				return
			}
		} else {
			service.Scope = input.Scope
			service.Endpoints = ""
			if err = tx.Save(&service).Error; err != nil {
				return
			}
		}

		for _, projectRoleName := range input.ProjectRoles {
			var projectRole base_db_model.ProjectRole
			if err = db.Where("name = ?", projectRoleName).First(&projectRole).Error; err != nil {
				return
			}

			if err = tx.Model(&projectRole).Association("Services").Append(&service).Error; err != nil {
				return
			}
		}
		return
	})
	return
}
