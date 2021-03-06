package server

import (
	"github.com/syunkitada/goapp/pkg/base/base_app"
	"github.com/syunkitada/goapp/pkg/base/base_config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/cluster/db_api"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_api/resolver"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_api/spec/genpkg"
	"github.com/syunkitada/goapp/pkg/resource/config"
	resource_api "github.com/syunkitada/goapp/pkg/resource/resource_api/spec/genpkg"
)

type Server struct {
	base_app.BaseApp
	baseConf     *base_config.Config
	clusterConf  *config.ResourceClusterConfig
	queryHandler *genpkg.QueryHandler
	dbApi        *db_api.Api
	resolver     *resolver.Resolver
	rootClient   *resource_api.Client
	clusterName  string
}

func New(baseConf *base_config.Config, mainConf *config.Config) *Server {
	clusterConf, ok := mainConf.Resource.ClusterMap[mainConf.Resource.ClusterName]
	if !ok {
		logger.StdoutFatalf("cluster config is not found: cluster=%s", mainConf.Resource.ClusterName)
	}
	clusterConf.Api.Name = "ReosurceClusterApi"

	dbApi := db_api.New(baseConf, &clusterConf)
	resolver := resolver.New(baseConf, &clusterConf, dbApi)
	queryHandler := genpkg.NewQueryHandler(baseConf, &clusterConf.Api, resolver)
	baseApp := base_app.New(baseConf, &clusterConf.Api, dbApi, queryHandler)

	srv := &Server{
		BaseApp:      baseApp,
		baseConf:     baseConf,
		clusterConf:  &clusterConf,
		queryHandler: queryHandler,
		dbApi:        dbApi,
		resolver:     resolver,
		rootClient:   resource_api.NewClient(&clusterConf.Api.RootClient),
		clusterName:  mainConf.Resource.ClusterName,
	}
	srv.SetDriver(srv)
	return srv
}
