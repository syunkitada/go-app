package resolver

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/base/base_const"
	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

func (resolver *Resolver) GetFloor(tctx *logger.TraceContext, input *spec.GetFloor, user *base_spec.UserAuthority) (data *spec.GetFloorData, code uint8, err error) {
	var floor *spec.Floor
	if floor, err = resolver.dbApi.GetFloor(tctx, input, user); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			code = base_const.CodeOkNotFound
			return
		}
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetFloorData{Floor: *floor}
	return
}

func (resolver *Resolver) GetFloors(tctx *logger.TraceContext, input *spec.GetFloors, user *base_spec.UserAuthority) (data *spec.GetFloorsData, code uint8, err error) {
	var floors []spec.Floor
	if floors, err = resolver.dbApi.GetFloors(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetFloorsData{Floors: floors}
	return
}

func (resolver *Resolver) CreateFloor(tctx *logger.TraceContext, input *spec.CreateFloor, user *base_spec.UserAuthority) (data *spec.CreateFloorData, code uint8, err error) {
	var specs []spec.Floor
	if specs, err = resolver.ConvertToFloorSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.CreateFloors(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkCreated
	data = &spec.CreateFloorData{}
	return
}

func (resolver *Resolver) UpdateFloor(tctx *logger.TraceContext, input *spec.UpdateFloor, user *base_spec.UserAuthority) (data *spec.UpdateFloorData, code uint8, err error) {
	var specs []spec.Floor
	if specs, err = resolver.ConvertToFloorSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.UpdateFloors(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkUpdated
	data = &spec.UpdateFloorData{}
	return
}

func (resolver *Resolver) DeleteFloor(tctx *logger.TraceContext, input *spec.DeleteFloor, user *base_spec.UserAuthority) (data *spec.DeleteFloorData, code uint8, err error) {
	if err = resolver.dbApi.DeleteFloor(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkDeleted
	data = &spec.DeleteFloorData{}
	return
}

func (resolver *Resolver) DeleteFloors(tctx *logger.TraceContext, input *spec.DeleteFloors, user *base_spec.UserAuthority) (data *spec.DeleteFloorsData, code uint8, err error) {
	var specs []spec.Floor
	if specs, err = resolver.ConvertToFloorSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.DeleteFloors(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkDeleted
	data = &spec.DeleteFloorsData{}
	return
}

func (resolver *Resolver) ConvertToFloorSpecs(specStr string) (specs []spec.Floor, err error) {
	var baseSpecs []base_spec.Spec
	if err = json.Unmarshal([]byte(specStr), &baseSpecs); err != nil {
		return
	}

	for _, base := range baseSpecs {
		if base.Kind != "Floor" {
			continue
		}
		var specBytes []byte
		if specBytes, err = json.Marshal(base.Spec); err != nil {
			return
		}
		var specData spec.Floor
		if err = json.Unmarshal(specBytes, &specData); err != nil {
			return
		}
		if err = resolver.Validate.Struct(&specData); err != nil {
			return
		}
		specs = append(specs, specData)
	}
	return
}
