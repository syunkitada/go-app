package resource_api

import (
	"fmt"

	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/resource_api_grpc_pb"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (srv *ResourceApiServer) MainTask(traceId string) error {
	if err := srv.UpdateNodeTask(traceId); err != nil {
		return err
	}

	return nil
}

func (srv *ResourceApiServer) UpdateNodeTask(traceId string) error {
	var err error
	startTime := logger.StartTaskTrace(traceId, srv.Host, srv.Name)
	defer func() {
		logger.EndTaskTrace(traceId, srv.Host, srv.Name, startTime, err)
	}()

	req := &resource_api_grpc_pb.UpdateNodeRequest{
		Name:         srv.Host,
		Kind:         resource_model.KindResourceApi,
		Role:         resource_model.RoleMember,
		Status:       resource_model.StatusEnabled,
		StatusReason: "Default",
		State:        resource_model.StateUp,
		StateReason:  "UpdateNode",
	}

	rep := srv.resourceModelApi.UpdateNode(req)
	if rep.Err != "" {
		err = fmt.Errorf(rep.Err)
		return err
	}
	return nil
}
