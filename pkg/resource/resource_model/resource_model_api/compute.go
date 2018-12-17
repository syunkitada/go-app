package resource_model_api

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/lib/codes"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_api/resource_cluster_api_grpc_pb"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/resource_api_grpc_pb"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (modelApi *ResourceModelApi) GetCompute(req *resource_api_grpc_pb.GetComputeRequest) *resource_api_grpc_pb.GetComputeReply {
	rep := &resource_api_grpc_pb.GetComputeReply{}

	db, err := gorm.Open("mysql", modelApi.conf.Resource.Database.Connection)
	defer db.Close()
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}
	db.LogMode(modelApi.conf.Default.EnableDatabaseLog)

	var computes []resource_model.Compute
	if err = db.Where("name like ?", req.Target).Find(&computes).Error; err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}

	rep.Computes = modelApi.convertComputes(req.TraceId, computes)
	rep.StatusCode = codes.Ok
	return rep
}

func (modelApi *ResourceModelApi) CreateCompute(req *resource_api_grpc_pb.CreateComputeRequest) *resource_api_grpc_pb.CreateComputeReply {
	rep := &resource_api_grpc_pb.CreateComputeReply{}

	db, err := gorm.Open("mysql", modelApi.conf.Resource.Database.Connection)
	defer db.Close()
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}
	db.LogMode(modelApi.conf.Default.EnableDatabaseLog)

	spec, statusCode, err := modelApi.validateComputeSpec(db, req.Spec)
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = statusCode
		return rep
	}

	var compute resource_model.Compute
	tx := db.Begin()
	defer tx.Rollback()
	if err = tx.Where("name = ? and cluster = ?", spec.Name, spec.Cluster).First(&compute).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			rep.Err = err.Error()
			rep.StatusCode = codes.RemoteDbError
			return rep
		}

		compute = resource_model.Compute{
			Cluster:      spec.Cluster,
			Kind:         spec.Kind,
			Name:         spec.Name,
			Spec:         req.Spec,
			Status:       resource_model.StatusActive,
			StatusReason: fmt.Sprintf("CreateCompute: user=%v, project=%v", req.UserName, req.ProjectName),
		}
		if err = tx.Create(&compute).Error; err != nil {
			rep.Err = err.Error()
			rep.StatusCode = codes.RemoteDbError
			return rep
		}
	} else {
		rep.Err = fmt.Sprintf("Already Exists: cluster=%v, name=%v",
			spec.Cluster, spec.Name)
		rep.StatusCode = codes.ClientAlreadyExists
		return rep
	}
	tx.Commit()

	computePb, err := modelApi.convertCompute(&compute)
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.ServerInternalError
		return rep
	}

	rep.Compute = computePb
	rep.StatusCode = codes.Ok
	return rep
}

func (modelApi *ResourceModelApi) UpdateCompute(req *resource_api_grpc_pb.UpdateComputeRequest) *resource_api_grpc_pb.UpdateComputeReply {
	rep := &resource_api_grpc_pb.UpdateComputeReply{}

	db, err := gorm.Open("mysql", modelApi.conf.Resource.Database.Connection)
	defer db.Close()
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}
	db.LogMode(modelApi.conf.Default.EnableDatabaseLog)

	spec, statusCode, err := modelApi.validateComputeSpec(db, req.Spec)
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = statusCode
		return rep
	}

	tx := db.Begin()
	defer tx.Rollback()
	var compute resource_model.Compute
	if err = tx.Where("name = ? and cluster = ?", spec.Name, spec.Cluster).First(&compute).Error; err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}

	compute.Spec = req.Spec
	compute.Status = resource_model.StatusActive
	compute.StatusReason = fmt.Sprintf("UpdateCompute: user=%v, project=%v", req.UserName, req.ProjectName)
	tx.Save(compute)
	tx.Commit()

	computePb, err := modelApi.convertCompute(&compute)
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.ServerInternalError
		return rep
	}

	rep.Compute = computePb
	rep.StatusCode = codes.Ok
	return rep
}

func (modelApi *ResourceModelApi) DeleteCompute(req *resource_api_grpc_pb.DeleteComputeRequest) *resource_api_grpc_pb.DeleteComputeReply {
	rep := &resource_api_grpc_pb.DeleteComputeReply{}

	db, err := gorm.Open("mysql", modelApi.conf.Resource.Database.Connection)
	defer db.Close()
	if err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}
	db.LogMode(modelApi.conf.Default.EnableDatabaseLog)

	tx := db.Begin()
	defer tx.Rollback()
	var compute resource_model.Compute
	if err = tx.Where("name = ?", req.Target).Delete(&compute).Error; err != nil {
		rep.Err = err.Error()
		rep.StatusCode = codes.RemoteDbError
		return rep
	}
	tx.Commit()

	rep.StatusCode = codes.Ok
	return rep
}

func (modelApi *ResourceModelApi) convertComputes(traceId string, computes []resource_model.Compute) []*resource_api_grpc_pb.Compute {
	pbComputes := make([]*resource_api_grpc_pb.Compute, len(computes))
	for i, compute := range computes {
		updatedAt, err := ptypes.TimestampProto(compute.Model.UpdatedAt)
		if err != nil {
			logger.TraceError(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg":    fmt.Sprintf("Failed ptypes.TimestampProto: %v", compute.Model.UpdatedAt),
				"Err":    err.Error(),
				"Method": "CreateCompute",
			})
			continue
		}
		createdAt, err := ptypes.TimestampProto(compute.Model.CreatedAt)
		if err != nil {
			logger.TraceError(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg":    fmt.Sprintf("Failed ptypes.TimestampProto: %v", compute.Model.CreatedAt),
				"Err":    err.Error(),
				"Method": "CreateCompute",
			})
			continue
		}

		pbComputes[i] = &resource_api_grpc_pb.Compute{
			Cluster:      compute.Cluster,
			Name:         compute.Name,
			Kind:         compute.Kind,
			Labels:       compute.Labels,
			Status:       compute.Status,
			StatusReason: compute.StatusReason,
			UpdatedAt:    updatedAt,
			CreatedAt:    createdAt,
		}
	}

	return pbComputes
}

func (modelApi *ResourceModelApi) convertCompute(compute *resource_model.Compute) (*resource_api_grpc_pb.Compute, error) {
	updatedAt, err := ptypes.TimestampProto(compute.Model.UpdatedAt)
	createdAt, err := ptypes.TimestampProto(compute.Model.CreatedAt)
	if err != nil {
		return nil, err
	}

	computePb := &resource_api_grpc_pb.Compute{
		Cluster:      compute.Cluster,
		Name:         compute.Name,
		Kind:         compute.Kind,
		Labels:       compute.Labels,
		Status:       compute.Status,
		StatusReason: compute.StatusReason,
		UpdatedAt:    updatedAt,
		CreatedAt:    createdAt,
	}

	return computePb, nil
}

func (modelApi *ResourceModelApi) validateComputeSpec(db *gorm.DB, specStr string) (resource_model.ComputeSpec, int64, error) {
	var spec resource_model.ComputeSpec
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
	case resource_model.SpecKindComputeLibvirt:
		// TODO Implement Validate SpecKindComputeLibvirt
		logger.Warning(modelApi.host, modelApi.name, "Validate SpecKindComputeLibvirt is not implemented")

	default:
		errors = append(errors, fmt.Sprintf("Invalid kind: %v", spec.Spec.Kind))
	}

	if len(errors) > 0 {
		return spec, codes.ClientInvalidRequest, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return spec, codes.Ok, nil
}

func (modelApi *ResourceModelApi) SyncCompute(traceId string) error {
	var err error
	db, err := gorm.Open("mysql", modelApi.conf.Resource.Database.Connection)
	defer db.Close()
	if err != nil {
		return err
	}
	db.LogMode(modelApi.conf.Default.EnableDatabaseLog)

	var computes []resource_model.Compute
	if err = db.Find(&computes).Error; err != nil {
		return err
	}

	computeMap := map[string]resource_cluster_api_grpc_pb.Compute{}
	req := resource_cluster_api_grpc_pb.GetComputeRequest{Target: "%"}
	for clusterName, clusterClient := range modelApi.clusterClientMap {
		result, err := clusterClient.GetCompute(&req)
		if err != nil {
			logger.TraceError(traceId, modelApi.host, modelApi.name, map[string]string{
				"Err": fmt.Sprintf("Failed GetCompute from %v: %v", clusterName, err),
			})
		}
		for _, compute := range result.Computes {
			computeMap[compute.FullName] = *compute
		}
	}

	for _, compute := range computes {
		switch compute.Status {
		case resource_model.StatusCreating:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
			modelApi.InitializeCompute(db, &compute, computeMap)
		case resource_model.StatusCreatingInitialized:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		case resource_model.StatusCreatingScheduled:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		case resource_model.StatusUpdating:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		case resource_model.StatusUpdatingScheduled:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		case resource_model.StatusDeleting:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		case resource_model.StatusDeletingScheduled:
			logger.TraceInfo(traceId, modelApi.host, modelApi.name, map[string]string{
				"Msg": fmt.Sprintf("Found %v resource: %v", compute.Status, compute.Name),
			})
		}
	}

	return nil
}

func (modelApi *ResourceModelApi) InitializeCompute(db *gorm.DB, compute *resource_model.Compute, computeMap map[string]resource_cluster_api_grpc_pb.Compute) error {
	// TODO
	// Assgin IP address
	return nil
}
