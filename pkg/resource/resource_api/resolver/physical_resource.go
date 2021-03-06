package resolver

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/base/base_const"
	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

func (resolver *Resolver) GetPhysicalResource(tctx *logger.TraceContext, input *spec.GetPhysicalResource, user *base_spec.UserAuthority) (data *spec.GetPhysicalResourceData, code uint8, err error) {
	var physicalResource *spec.PhysicalResource
	if physicalResource, err = resolver.dbApi.GetPhysicalResource(tctx, input, user); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			code = base_const.CodeOkNotFound
			return
		}
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetPhysicalResourceData{PhysicalResource: *physicalResource}
	return
}

func (resolver *Resolver) GetPhysicalResources(tctx *logger.TraceContext, input *spec.GetPhysicalResources, user *base_spec.UserAuthority) (data *spec.GetPhysicalResourcesData, code uint8, err error) {
	var physicalResources []spec.PhysicalResource
	if physicalResources, err = resolver.dbApi.GetPhysicalResources(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetPhysicalResourcesData{PhysicalResources: physicalResources}
	return
}

func (resolver *Resolver) CreatePhysicalResource(tctx *logger.TraceContext, input *spec.CreatePhysicalResource, user *base_spec.UserAuthority) (data *spec.CreatePhysicalResourceData, code uint8, err error) {
	var specs []spec.PhysicalResource
	if specs, err = resolver.ConvertToPhysicalResourceSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.CreatePhysicalResources(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkCreated
	data = &spec.CreatePhysicalResourceData{}
	return
}

func (resolver *Resolver) UpdatePhysicalResource(tctx *logger.TraceContext, input *spec.UpdatePhysicalResource, user *base_spec.UserAuthority) (data *spec.UpdatePhysicalResourceData, code uint8, err error) {
	var specs []spec.PhysicalResource
	if specs, err = resolver.ConvertToPhysicalResourceSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.UpdatePhysicalResources(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkUpdated
	data = &spec.UpdatePhysicalResourceData{}
	return
}

func (resolver *Resolver) DeletePhysicalResource(tctx *logger.TraceContext, input *spec.DeletePhysicalResource, user *base_spec.UserAuthority) (data *spec.DeletePhysicalResourceData, code uint8, err error) {
	if err = resolver.dbApi.DeletePhysicalResource(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkDeleted
	data = &spec.DeletePhysicalResourceData{}
	return
}

func (resolver *Resolver) DeletePhysicalResources(tctx *logger.TraceContext, input *spec.DeletePhysicalResources, user *base_spec.UserAuthority) (data *spec.DeletePhysicalResourcesData, code uint8, err error) {
	var specs []spec.PhysicalResource
	if specs, err = resolver.ConvertToPhysicalResourceSpecs(input.Spec); err != nil {
		code = base_const.CodeClientBadRequest
		return
	}
	if err = resolver.dbApi.DeletePhysicalResources(tctx, specs, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOkDeleted
	data = &spec.DeletePhysicalResourcesData{}
	return
}

func (resolver *Resolver) ConvertToPhysicalResourceSpecs(specStr string) (specs []spec.PhysicalResource, err error) {
	var baseSpecs []base_spec.Spec
	if err = json.Unmarshal([]byte(specStr), &baseSpecs); err != nil {
		return
	}

	for _, base := range baseSpecs {
		if base.Kind != "PhysicalResource" {
			continue
		}
		var specBytes []byte
		if specBytes, err = json.Marshal(base.Spec); err != nil {
			return
		}
		var specData spec.PhysicalResource
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
