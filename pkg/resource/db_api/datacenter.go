package db_api

import (
	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/db_model"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

func (api *Api) GetDatacenter(tctx *logger.TraceContext, input *spec.GetDatacenter, user *base_spec.UserAuthority) (data *spec.Datacenter, err error) {
	data = &spec.Datacenter{}
	err = api.DB.Where("name = ? AND deleted_at IS NULL", input.Name).First(data).Error
	return
}

func (api *Api) GetDatacenters(tctx *logger.TraceContext, input *spec.GetDatacenters, user *base_spec.UserAuthority) (data []spec.Datacenter, err error) {
	err = api.DB.Where("deleted_at IS NULL").Find(&data).Error
	return
}

func (api *Api) CreateDatacenters(tctx *logger.TraceContext, input []spec.Datacenter, user *base_spec.UserAuthority) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, val := range input {
			var tmp db_model.Datacenter
			if err = tx.Where("name = ?", val.Name).First(&tmp).Error; err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					return
				}
				tmp = db_model.Datacenter{
					Name:         val.Name,
					Kind:         val.Kind,
					Description:  val.Description,
					Region:       val.Region,
					DomainSuffix: val.DomainSuffix,
				}
				if err = tx.Create(&tmp).Error; err != nil {
					return
				}
			}
		}
		return
	})
	return
}

func (api *Api) UpdateDatacenters(tctx *logger.TraceContext, input []spec.Datacenter, user *base_spec.UserAuthority) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, val := range input {
			if err = tx.Model(&db_model.Datacenter{}).Where("name = ?", val.Name).Updates(&db_model.Datacenter{
				Kind:        val.Kind,
				Description: val.Description,
			}).Error; err != nil {
				return
			}
		}
		return
	})
	return
}

func (api *Api) DeleteDatacenter(tctx *logger.TraceContext, input *spec.DeleteDatacenter, user *base_spec.UserAuthority) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		err = tx.Where("name = ?", input.Name).Delete(&db_model.Datacenter{}).Error
		return
	})
	return
}

func (api *Api) DeleteDatacenters(tctx *logger.TraceContext, input []spec.Datacenter, user *base_spec.UserAuthority) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, data := range input {
			if err = tx.Where("name = ?", data.Name).
				Delete(&db_model.Datacenter{}).Error; err != nil {
				return
			}
		}
		return
	})
	return
}
