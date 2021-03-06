package resolver

import (
	"github.com/syunkitada/goapp/pkg/base/base_const"
	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
)

func (resolver *Resolver) GetNodes(tctx *logger.TraceContext, input *spec.GetNodes, user *base_spec.UserAuthority) (data *spec.GetNodesData, code uint8, err error) {
	var nodes []spec.Node
	if nodes, err = resolver.dbApi.GetNodes(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetNodesData{Nodes: nodes}
	return
}

func (resolver *Resolver) GetNode(tctx *logger.TraceContext, input *spec.GetNode, user *base_spec.UserAuthority) (data *spec.GetNodeData, code uint8, err error) {
	var node *spec.Node
	if node, err = resolver.dbApi.GetNode(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	data = &spec.GetNodeData{Node: *node}
	return
}

func (resolver *Resolver) GetNodeMetrics(tctx *logger.TraceContext, input *spec.GetNodeMetrics, user *base_spec.UserAuthority) (data *spec.GetNodeMetricsData, code uint8, err error) {
	if data, err = resolver.dbApi.GetNodeMetrics(tctx, input, user); err != nil {
		code = base_const.CodeServerInternalError
		return
	}
	code = base_const.CodeOk
	return
}
