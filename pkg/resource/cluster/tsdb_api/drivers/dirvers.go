package drivers

import (
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/cluster/tsdb_api/drivers/influxdb_driver"
	"github.com/syunkitada/goapp/pkg/resource/config"
	api_spec "github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

type TsdbDriver interface {
	Report(tctx *logger.TraceContext, input *api_spec.ReportNode) (err error)
	GetNode(tctx *logger.TraceContext, input *api_spec.GetNode) (data []api_spec.MetricsGroup, err error)
	GetLogParams(tctx *logger.TraceContext, input *api_spec.GetLogParams) (data *api_spec.GetLogParamsData, err error)
	GetLogs(tctx *logger.TraceContext, input *api_spec.GetLogs) (data *api_spec.GetLogsData, err error)
}

func Load(clusterConf *config.ResourceClusterConfig) TsdbDriver {
	switch clusterConf.TimeSeriesDatabase.Driver {
	case "influxdb":
		driver := influxdb_driver.New(clusterConf)
		return driver
	}

	return nil
}
