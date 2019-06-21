package resource_model_api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_grpc_pb"
	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_model"
	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_utils"
	"github.com/syunkitada/goapp/pkg/authproxy/index_model"
	"github.com/syunkitada/goapp/pkg/lib/codes"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (modelApi *ResourceModelApi) PhysicalAction(tctx *logger.TraceContext,
	req *authproxy_grpc_pb.ActionRequest, rep *authproxy_grpc_pb.ActionReply) {
	var err error
	var statusCode int64
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	data := map[string]interface{}{}
	response := authproxy_model.ActionResponse{
		Tctx: *req.Tctx,
	}

	var db *gorm.DB
	if db, err = modelApi.open(tctx); err != nil {
		authproxy_utils.MergeResponse(rep, &response, data, err, codes.RemoteDbError)
		return
	}
	defer func() {
		tmpErr := db.Close()
		if tmpErr != nil {
			logger.Error(tctx, tmpErr)
		}
	}()

	statusCode = codes.Unknown
	for _, query := range req.Queries {
		switch query.Kind {
		case "GetIndex":
			response.Index = *modelApi.getPhysicalIndex()
		case "GetDashboardIndex":
			response.Index = *modelApi.getPhysicalIndex()
			statusCode, err = modelApi.GetDatacenters(tctx, db, query, data)

		case "GetDatacenter":
			statusCode, err = modelApi.GetDatacenter(tctx, db, query, data)
		case "GetDatacenters":
			statusCode, err = modelApi.GetDatacenters(tctx, db, query, data)
		case "CreateDatacenter":
			statusCode, err = modelApi.CreateDatacenter(tctx, db, query)
		case "UpdateDatacenter":
			statusCode, err = modelApi.UpdateDatacenter(tctx, db, query)
		case "DeleteDatacenter":
			statusCode, err = modelApi.DeleteDatacenter(tctx, db, query)

		case "GetFloor":
			statusCode, err = modelApi.GetFloor(tctx, db, query, data)
		case "GetFloors":
			statusCode, err = modelApi.GetFloors(tctx, db, query, data)
		case "CreateFloor":
			statusCode, err = modelApi.CreateFloor(tctx, db, query)
		case "UpdateFloor":
			statusCode, err = modelApi.UpdateFloor(tctx, db, query)
		case "DeleteFloor":
			statusCode, err = modelApi.DeleteFloor(tctx, db, query)

		case "GetRack":
			statusCode, err = modelApi.GetRack(tctx, db, query, data)
		case "GetRacks":
			statusCode, err = modelApi.GetRacks(tctx, db, query, data)
		case "CreateRack":
			statusCode, err = modelApi.CreateRack(tctx, db, query)
		case "UpdateRack":
			statusCode, err = modelApi.UpdateRack(tctx, db, query)
		case "DeleteRack":
			statusCode, err = modelApi.DeleteRack(tctx, db, query)

		case "GetPhysicalResource":
			statusCode, err = modelApi.GetPhysicalResource(tctx, db, query, data)
		case "GetPhysicalResources":
			statusCode, err = modelApi.GetPhysicalResources(tctx, db, query, data)
		case "CreatePhysicalResource":
			statusCode, err = modelApi.CreatePhysicalResource(tctx, db, query)
		case "UpdatePhysicalResource":
			statusCode, err = modelApi.UpdatePhysicalResource(tctx, db, query)
		case "DeletePhysicalResource":
			statusCode, err = modelApi.DeletePhysicalResource(tctx, db, query)

		case "GetPhysicalModel":
			statusCode, err = modelApi.GetPhysicalModel(tctx, db, query, data)
		case "GetPhysicalModels":
			statusCode, err = modelApi.GetPhysicalModels(tctx, db, query, data)
		case "CreatePhysicalModel":
			statusCode, err = modelApi.CreatePhysicalModel(tctx, db, query)
		case "UpdatePhysicalModel":
			statusCode, err = modelApi.UpdatePhysicalModel(tctx, db, query)
		case "DeletePhysicalModel":
			statusCode, err = modelApi.DeletePhysicalModel(tctx, db, query)
		}

		if err != nil {
			break
		}
	}

	authproxy_utils.MergeResponse(rep, &response, data, err, statusCode)
}

func (modelApi *ResourceModelApi) getPhysicalIndex() *index_model.Index {
	cmdMap := map[string]index_model.Cmd{}
	cmdMaps := []map[string]index_model.Cmd{
		resource_model.DatacenterCmd,
		resource_model.RackCmd,
		resource_model.FloorCmd,
		resource_model.PhysicalModelCmd,
		resource_model.PhysicalResourceCmd,
	}
	for _, tmpCmdMap := range cmdMaps {
		for key, cmd := range tmpCmdMap {
			cmdMap[key] = cmd
		}
	}

	return &index_model.Index{
		SyncDelay: 20000,
		CmdMap:    cmdMap,
		View: index_model.Panels{
			Name: "Root",
			Kind: "RoutePanels",
			Panels: []interface{}{
				resource_model.DatacentersTable,
				index_model.Tabs{
					Name:             "Resources",
					Kind:             "RouteTabs",
					Subname:          "kind",
					Route:            "/Datacenters/:datacenter/Resources/:kind",
					TabParam:         "kind",
					GetQueries:       []string{"GetPhysicalResources", "GetRacks", "GetFloors", "GetPhysicalModels"},
					ExpectedDataKeys: []string{"PhysicalResources", "Racks", "Floors", "PhysicalModels"},
					IsSync:           true,
					Tabs: []interface{}{
						resource_model.PhysicalResourcesTable,
						resource_model.RacksTable,
						resource_model.FloorsTable,
						resource_model.PhysicalModelsTable,
					}, // Tabs
				},
				gin.H{
					"Name":      "Resource",
					"Subname":   "resource",
					"Route":     "/Datacenters/:datacenter/Resources/:kind/Detail/:resource/:subkind",
					"Kind":      "RoutePanes",
					"PaneParam": "kind",
					"Panes": []interface{}{
						resource_model.PhysicalModelsDetail,
						resource_model.PhysicalResourcesDetail,
					},
				},
			},
		},
	}
}
