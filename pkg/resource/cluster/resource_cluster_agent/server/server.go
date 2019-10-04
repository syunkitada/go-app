package server

import (
	"github.com/syunkitada/goapp/pkg/base/base_app"
	"github.com/syunkitada/goapp/pkg/base/base_config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_agent/resolver"
	"github.com/syunkitada/goapp/pkg/resource/cluster/resource_cluster_agent/spec/genpkg"
	"github.com/syunkitada/goapp/pkg/resource/config"
)

type Server struct {
	base_app.BaseApp
	baseConf     *base_config.Config
	clusterConf  *config.ResourceClusterConfig
	queryHandler *genpkg.QueryHandler
}

func New(baseConf *base_config.Config, mainConf *config.Config) *Server {
	clusterConf, ok := mainConf.Resource.ClusterMap[mainConf.Resource.ClusterName]
	if !ok {
		logger.StdoutFatalf("cluster config is not found: cluster=%s", mainConf.Resource.ClusterName)
	}

	resolver := resolver.New(baseConf, &clusterConf)
	queryHandler := genpkg.NewQueryHandler(baseConf, &clusterConf.Api, resolver)
	baseApp := base_app.New(baseConf, &clusterConf.Api, nil, queryHandler)

	srv := &Server{
		BaseApp:      baseApp,
		baseConf:     baseConf,
		clusterConf:  &clusterConf,
		queryHandler: queryHandler,
	}
	srv.SetDriver(srv)
	return srv
}