package resource_api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/syunkitada/goapp/pkg/base"
	"github.com/syunkitada/goapp/pkg/config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/resource_api_grpc_pb"
	"github.com/syunkitada/goapp/pkg/resource/resource_model/resource_model_api"
)

type ResourceApiServer struct {
	base.BaseApp
	conf             *config.Config
	resourceModelApi *resource_model_api.ResourceModelApi
}

func NewResourceApiServer(conf *config.Config) *ResourceApiServer {
	conf.Resource.ApiApp.Name = "resource.api"
	server := ResourceApiServer{
		BaseApp:          base.NewBaseApp(conf, &conf.Resource.ApiApp),
		conf:             conf,
		resourceModelApi: resource_model_api.NewResourceModelApi(conf, nil),
	}

	server.RegisterDriver(&server)

	return &server
}

func (cli *ResourceApiServer) newTraceContext(host string, app string, ctx context.Context, tctx *resource_api_grpc_pb.TraceContext) *logger.TraceContext {
	var client string
	if pr, ok := peer.FromContext(ctx); ok {
		client = pr.Addr.String()
	} else {
		client = ""
	}

	return &logger.TraceContext{
		TraceId: tctx.TraceId,
		Host:    host,
		App:     app,
		Metadata: map[string]string{
			"Client":          client,
			"ActionName":      tctx.ActionName,
			"UserName":        tctx.UserName,
			"RoleName":        tctx.RoleName,
			"ProjectName":     tctx.ProjectName,
			"ProjectRoleName": tctx.ProjectRoleName,
		},
	}
}

func (srv *ResourceApiServer) RegisterGrpcServer(grpcServer *grpc.Server) error {
	resource_api_grpc_pb.RegisterResourceApiServer(grpcServer, srv)
	return nil
}

func (srv *ResourceApiServer) Status(ctx context.Context, statusRequest *resource_api_grpc_pb.StatusRequest) (*resource_api_grpc_pb.StatusReply, error) {
	return &resource_api_grpc_pb.StatusReply{Msg: "Status"}, nil
}

//
// Action
//
func (srv *ResourceApiServer) Action(ctx context.Context, req *resource_api_grpc_pb.ActionRequest) (*resource_api_grpc_pb.ActionReply, error) {
	var err error
	rep := &resource_api_grpc_pb.ActionReply{Tctx: req.Tctx}
	tctx := srv.newTraceContext(srv.Host, srv.Name, ctx, req.Tctx)
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	switch req.Tctx.ActionName {
	case "GetCluster":
		srv.resourceModelApi.GetCluster(tctx, req, rep)
	}

	return rep, nil
}

//
// Node
//
func (srv *ResourceApiServer) UpdateNode(ctx context.Context, req *resource_api_grpc_pb.UpdateNodeRequest) (*resource_api_grpc_pb.UpdateNodeReply, error) {
	// TODO
	tctx := logger.NewGrpcTraceContext(srv.Host, srv.Name, ctx)
	// startTime := logger.StartTrace(tctx)
	rep := srv.resourceModelApi.UpdateNode(tctx, req)
	// logger.EndGrpcTrace(tctx, startTime, rep.StatusCode, rep.Err)
	return rep, nil
}
