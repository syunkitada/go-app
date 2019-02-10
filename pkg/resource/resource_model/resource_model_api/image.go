package resource_model_api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/lib/codes"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/resource_api_grpc_pb"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (modelApi *ResourceModelApi) GetImage(tctx *logger.TraceContext, req *resource_api_grpc_pb.ActionRequest, rep *resource_api_grpc_pb.ActionReply) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	var db *gorm.DB
	if db, err = modelApi.open(tctx); err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}
	defer db.Close()

	var images []resource_model.Image
	if err = db.Where("name like ?", req.Target).Find(&images).Error; err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}

	rep.Images = modelApi.convertImages(tctx, images)
	rep.Tctx.StatusCode = codes.Ok
	return
}

func (modelApi *ResourceModelApi) CreateImage(tctx *logger.TraceContext, req *resource_api_grpc_pb.ActionRequest, rep *resource_api_grpc_pb.ActionReply) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	var db *gorm.DB
	if db, err = modelApi.open(tctx); err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}
	defer db.Close()

	spec, statusCode, err := modelApi.validateImageSpec(db, req.Spec)
	if err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = statusCode
		return
	}

	var image resource_model.Image
	tx := db.Begin()
	defer tx.Rollback()
	if err = tx.Where("name = ? and cluster = ?", spec.Name, spec.Cluster).First(&image).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			rep.Tctx.Err = err.Error()
			rep.Tctx.StatusCode = codes.RemoteDbError
			return
		}

		image = resource_model.Image{
			Cluster:      spec.Cluster,
			Kind:         spec.Kind,
			Name:         spec.Name,
			Spec:         req.Spec,
			Status:       resource_model.StatusActive,
			StatusReason: fmt.Sprintf("CreateImage: user=%v, project=%v", req.Tctx.UserName, req.Tctx.ProjectName),
		}
		if err = tx.Create(&image).Error; err != nil {
			rep.Tctx.Err = err.Error()
			rep.Tctx.StatusCode = codes.RemoteDbError
			return
		}
	} else {
		rep.Tctx.Err = fmt.Sprintf("Already Exists: cluster=%v, name=%v", spec.Cluster, spec.Name)
		rep.Tctx.StatusCode = codes.ClientAlreadyExists
		return
	}
	tx.Commit()

	imagePb, err := modelApi.convertImage(&image)
	if err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.ServerInternalError
		return
	}

	rep.Images = []*resource_api_grpc_pb.Image{imagePb}
	rep.Tctx.StatusCode = codes.Ok
	return
}

func (modelApi *ResourceModelApi) UpdateImage(tctx *logger.TraceContext, req *resource_api_grpc_pb.ActionRequest, rep *resource_api_grpc_pb.ActionReply) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	var db *gorm.DB
	if db, err = modelApi.open(tctx); err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}
	defer db.Close()

	spec, statusCode, err := modelApi.validateImageSpec(db, req.Spec)
	if err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = statusCode
		return
	}

	tx := db.Begin()
	defer tx.Rollback()
	var image resource_model.Image
	if err = tx.Where("name = ? and cluster = ?", spec.Name, spec.Cluster).First(&image).Error; err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}

	image.Spec = req.Spec
	image.Status = resource_model.StatusActive
	image.StatusReason = fmt.Sprintf("UpdateImage: user=%v, project=%v", req.Tctx.UserName, req.Tctx.ProjectName)
	tx.Save(image)
	tx.Commit()

	imagePb, err := modelApi.convertImage(&image)
	if err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.ServerInternalError
		return
	}

	rep.Images = []*resource_api_grpc_pb.Image{imagePb}
	rep.Tctx.StatusCode = codes.Ok
}

func (modelApi *ResourceModelApi) DeleteImage(tctx *logger.TraceContext, req *resource_api_grpc_pb.ActionRequest, rep *resource_api_grpc_pb.ActionReply) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	var db *gorm.DB
	if db, err = modelApi.open(tctx); err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}
	defer db.Close()

	tx := db.Begin()
	defer tx.Rollback()
	var image resource_model.Image
	if err = tx.Where("name = ?", req.Target).Delete(&image).Error; err != nil {
		rep.Tctx.Err = err.Error()
		rep.Tctx.StatusCode = codes.RemoteDbError
		return
	}
	tx.Commit()

	rep.Tctx.StatusCode = codes.Ok
}

func (modelApi *ResourceModelApi) convertImages(tctx *logger.TraceContext, images []resource_model.Image) []*resource_api_grpc_pb.Image {
	pbImages := make([]*resource_api_grpc_pb.Image, len(images))
	for i, image := range images {
		updatedAt, err := ptypes.TimestampProto(image.Model.UpdatedAt)
		if err != nil {
			logger.Warningf(tctx, err, "Failed ptypes.TimestampProto: %v", image.Model.UpdatedAt)
			continue
		}
		createdAt, err := ptypes.TimestampProto(image.Model.CreatedAt)
		if err != nil {
			logger.Warningf(tctx, err, "Failed ptypes.TimestampProto: %v", image.Model.CreatedAt)
			continue
		}

		pbImages[i] = &resource_api_grpc_pb.Image{
			Cluster:      image.Cluster,
			Name:         image.Name,
			Kind:         image.Kind,
			Labels:       image.Labels,
			Status:       image.Status,
			StatusReason: image.StatusReason,
			UpdatedAt:    updatedAt,
			CreatedAt:    createdAt,
		}
	}

	return pbImages
}

func (modelApi *ResourceModelApi) convertImage(image *resource_model.Image) (*resource_api_grpc_pb.Image, error) {
	updatedAt, err := ptypes.TimestampProto(image.Model.UpdatedAt)
	createdAt, err := ptypes.TimestampProto(image.Model.CreatedAt)
	if err != nil {
		return nil, err
	}

	imagePb := &resource_api_grpc_pb.Image{
		Cluster:      image.Cluster,
		Name:         image.Name,
		Kind:         image.Kind,
		Labels:       image.Labels,
		Status:       image.Status,
		StatusReason: image.StatusReason,
		UpdatedAt:    updatedAt,
		CreatedAt:    createdAt,
	}

	return imagePb, nil
}

func (modelApi *ResourceModelApi) validateImageSpec(db *gorm.DB, specStr string) (resource_model.ImageSpec, int64, error) {
	var spec resource_model.ImageSpec
	var err error
	if err = json.Unmarshal([]byte(specStr), &spec); err != nil {
		return spec, codes.ClientBadRequest, err
	}
	if err = modelApi.validate.Struct(spec); err != nil {
		return spec, codes.ClientInvalidRequest, err
	}

	ok, err := modelApi.ValidateClusterName(db, spec.Cluster)
	if err != nil {
		return spec, codes.RemoteDbError, err
	}
	if !ok {
		return spec, codes.ClientInvalidRequest, fmt.Errorf("Invalid cluster: %v", spec.Cluster)
	}

	errors := []string{}
	switch spec.Spec.Kind {
	case resource_model.SpecKindImageUrl:
		_, err := url.Parse(spec.Spec.Url)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Invalid ImageUrl: %v", err.Error()))
		}

	default:
		errors = append(errors, fmt.Sprintf("Invalid kind: %v", spec.Spec.Kind))
	}

	if len(errors) > 0 {
		return spec, codes.ClientInvalidRequest, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return spec, codes.Ok, nil
}
