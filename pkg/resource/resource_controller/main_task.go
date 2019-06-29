package resource_controller

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"

	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_model"
	"github.com/syunkitada/goapp/pkg/lib/json_utils"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (srv *ResourceControllerServer) MainTask(tctx *logger.TraceContext) error {
	if err := srv.UpdateNode(tctx); err != nil {
		return err
	}
	if err := srv.SyncRole(tctx); err != nil {
		return err
	}
	if srv.role == resource_model.RoleMember {
		return nil
	}

	if err := srv.resourceModelApi.CheckNodes(tctx); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go srv.SyncRegionService(tctx, &wg)
	wg.Wait()

	// TODO
	// implement with goroutine
	// check compute
	// check container
	// check image
	// check loadbalancer

	return nil
}

func (srv *ResourceControllerServer) UpdateNode(tctx *logger.TraceContext) error {
	nodes := []resource_model.NodeSpec{
		resource_model.NodeSpec{
			Name:         srv.conf.Default.Host,
			Kind:         resource_model.KindResourceApi,
			Role:         resource_model.RoleMember,
			Status:       resource_model.StatusEnabled,
			StatusReason: "Default",
			State:        resource_model.StateUp,
			StateReason:  "UpdateNode",
		},
	}
	specs, err := json_utils.Marshal(nodes)
	if err != nil {
		return err
	}
	queries := []authproxy_model.Query{
		authproxy_model.Query{
			Kind: "update_node",
			StrParams: map[string]string{
				"Specs": string(specs),
			},
		},
	}

	if _, err := srv.resourceApiClient.VirtualAction(
		logger.NewActionTraceContext(tctx, "system", "system", queries)); err != nil {
		return err
	}

	return nil
}

func (srv *ResourceControllerServer) SyncRole(tctx *logger.TraceContext) error {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	nodes, err := srv.resourceModelApi.SyncRole(tctx, resource_model.KindResourceController)
	if err != nil {
		return err
	}

	existsSelfNode := false
	existsActiveLeader := false
	for _, node := range nodes {
		if node.Kind != resource_model.KindResourceController {
			continue
		}
		if node.Name == srv.conf.Default.Host && node.Status == resource_model.StatusEnabled && node.State == resource_model.StateUp {
			existsSelfNode = true
			srv.role = node.Role
		}
		if node.Status == resource_model.StatusEnabled && node.State == resource_model.StateUp {
			if node.Role == resource_model.RoleLeader {
				existsActiveLeader = true
			}
		}
	}

	if !existsSelfNode {
		err = fmt.Errorf("This node is not activated")
		return err
	}

	if !existsActiveLeader {
		err = fmt.Errorf("Active Leader is not exists, after ReassignNode")
		return err
	}

	return nil
}

func (srv *ResourceControllerServer) SyncRegionService(tctx *logger.TraceContext, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	errChan := make(chan error)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, srv.syncResourceTimeout)
	defer cancel()

	go func() {
		errChan <- srv.resourceModelApi.SyncRegionService(tctx)
	}()

	select {
	case err = <-errChan:
		break
	case <-ctx.Done():
		err = ctx.Err()
	}
}
