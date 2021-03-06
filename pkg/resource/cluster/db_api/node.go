package db_api

import (
	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/lib/error_utils"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/db_model"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
	api_spec "github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

func (api *Api) GetNodes(tctx *logger.TraceContext, input *spec.GetNodes) (data []spec.Node, err error) {
	err = api.DB.Find(&data).Error
	return
}

func (api *Api) GetNode(tctx *logger.TraceContext, input *spec.GetNode) (data spec.Node, err error) {
	var nodes []spec.Node
	err = api.DB.Where("name = ?", input.Name).Find(&nodes).Error
	if len(nodes) == 0 {
		err = error_utils.NewNotFoundError(input.Name)
	}
	data = nodes[0]
	return
}

func (api *Api) ReportNode(tctx *logger.TraceContext, input *api_spec.ReportNode) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		var tmpNode db_model.Node
		if err = tx.Table("nodes").Where(
			"name = ?", input.Name).First(&tmpNode).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return
			}
			tmpNode = db_model.Node{
				Name:  input.Name,
				State: input.State,
			}
			if err = tx.Create(&tmpNode).Error; err != nil {
				return
			}
		} else {
			tmpNode.State = input.State
			if err = tx.Save(&tmpNode).Error; err != nil {
				return
			}
		}
		return
	})

	return
}
