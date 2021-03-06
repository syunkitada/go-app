package server

import (
	"github.com/syunkitada/goapp/pkg/base/base_app"
	"github.com/syunkitada/goapp/pkg/base/base_config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/cluster/db_api"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_controller/resolver"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_controller/spec/genpkg"
	"github.com/syunkitada/goapp/pkg/resource/config"
	"github.com/syunkitada/goapp/pkg/resource/consts"
)

type Server struct {
	base_app.BaseApp
	baseConf     *base_config.Config
	clusterConf  *config.ResourceClusterConfig
	queryHandler *genpkg.QueryHandler
	dbApi        *db_api.Api
}

func New(baseConf *base_config.Config, mainConf *config.Config) *Server {
	clusterConf, ok := mainConf.Resource.ClusterMap[mainConf.Resource.ClusterName]
	if !ok {
		logger.StdoutFatalf("cluster config is not found: cluster=%s", mainConf.Resource.ClusterName)
	}
	clusterConf.Controller.Name = consts.KindResourceClusterController

	dbApi := db_api.New(baseConf, &clusterConf)
	resolver := resolver.New(baseConf, &clusterConf, dbApi)
	queryHandler := genpkg.NewQueryHandler(baseConf, &clusterConf.Controller, resolver)
	baseApp := base_app.New(baseConf, &clusterConf.Controller, dbApi, queryHandler)

	srv := &Server{
		BaseApp:      baseApp,
		baseConf:     baseConf,
		clusterConf:  &clusterConf,
		queryHandler: queryHandler,
		dbApi:        dbApi,
	}
	srv.SetDriver(srv)
	return srv
}
