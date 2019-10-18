package db_api

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/base/base_const"
	"github.com/syunkitada/goapp/pkg/base/base_db_model"
	"github.com/syunkitada/goapp/pkg/lib/error_utils"
	"github.com/syunkitada/goapp/pkg/lib/json_utils"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/consts"
	"github.com/syunkitada/goapp/pkg/resource/db_model"
	"github.com/syunkitada/goapp/pkg/resource/resource_api/spec"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

func (api *Api) GetCompute(tctx *logger.TraceContext, input *spec.GetCompute) (data *spec.Compute, err error) {
	data = &spec.Compute{}
	err = api.DB.Where("name = ?", input.Name).First(data).Error
	return
}

func (api *Api) GetComputes(tctx *logger.TraceContext, input *spec.GetComputes) (data []spec.Compute, err error) {
	err = api.DB.Find(&data).Error
	return
}

func (api *Api) CreateComputes(tctx *logger.TraceContext, specs []spec.RegionServiceComputeSpec) (err error) {
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()
	fmt.Println("CreateComputes", specs)

	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, spec := range specs {
			var specBytes []byte
			specBytes, err = json_utils.Marshal(spec)
			data := db_model.Compute{
				Name:         spec.Name,
				Spec:         string(specBytes),
				Status:       base_const.StatusCreating,
				StatusReason: "CreateComputes",
			}
			if err = tx.Create(&data).Error; err != nil {
				return
			}
		}
		return
	})
	return
}

func (api *Api) UpdateComputes(tctx *logger.TraceContext, computes []spec.RegionServiceComputeSpec) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, compute := range computes {
			if err = tx.Model(&db_model.Compute{}).Where("name = ?", compute.Name).Updates(&db_model.Compute{
				Kind: compute.Kind,
			}).Error; err != nil {
				return
			}
		}
		return
	})
	return
}

func (api *Api) DeleteCompute(tctx *logger.TraceContext, input *spec.DeleteCompute) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		err = tx.Where("name = ?", input.Name).Unscoped().Delete(&db_model.Compute{}).Error
		return
	})
	return
}

func (api *Api) DeleteComputes(tctx *logger.TraceContext, input []spec.RegionServiceComputeSpec) (err error) {
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		return
	})
	return
}

func (api *Api) SyncCompute(tctx *logger.TraceContext) (err error) {
	fmt.Println("SyncCompute")

	var computes []db_model.Compute
	var nodes []db_model.ClusterNode
	var computeAssignments []db_model.ComputeAssignmentWithComputeAndNode
	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		if err = tx.Find(&computes).Error; err != nil {
			return
		}

		// TODO filter by resource driver
		if err = tx.Where(&db_model.ClusterNode{
			Node: base_db_model.Node{
				Kind: consts.KindResourceClusterAgent,
			},
		}).Find(&nodes).Error; err != nil {
			return
		}

		if computeAssignments, err = api.GetComputeAssignments(tctx, tx, ""); err != nil {
			return
		}
		return
	})

	nodeMap := map[uint]*db_model.ClusterNode{}
	nodeAssignmentsMap := map[uint][]db_model.ComputeAssignmentWithComputeAndNode{}
	for _, node := range nodes {
		nodeAssignmentsMap[node.ID] = []db_model.ComputeAssignmentWithComputeAndNode{}
		nodeMap[node.ID] = &node
	}

	computeAssignmentsMap := map[string][]db_model.ComputeAssignmentWithComputeAndNode{}
	for _, assignment := range computeAssignments {
		assignments, ok := computeAssignmentsMap[assignment.ComputeName]
		if !ok {
			assignments = []db_model.ComputeAssignmentWithComputeAndNode{}
		}
		assignments = append(assignments, assignment)
		computeAssignmentsMap[assignment.ComputeName] = assignments

		nodeAssignments := nodeAssignmentsMap[assignment.NodeID]
		nodeAssignments = append(nodeAssignments, assignment)
		nodeAssignmentsMap[assignment.NodeID] = nodeAssignments
	}

	for _, compute := range computes {
		switch compute.Status {
		case base_const.StatusInitializing:
			api.AssignCompute(tctx, &compute, nodeMap, nodeAssignmentsMap, computeAssignmentsMap, false)
		case base_const.StatusCreatingScheduled:
			api.ConfirmCreatingScheduledCompute(tctx, &compute, computeAssignmentsMap)
		}
	}

	return
}

func (api *Api) GetComputeAssignments(tctx *logger.TraceContext, db *gorm.DB,
	nodeName string) (assignments []db_model.ComputeAssignmentWithComputeAndNode, err error) {
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	query := db.Table("compute_assignments as ca").
		Select("ca.id, ca.status, ca.updated_at, ca.compute_id, c.name as compute_name, c.spec as compute_spec, ca.node_id, n.name as node_name").
		Joins("INNER JOIN computes AS c ON c.id = ca.compute_id").
		Joins("INNER JOIN nodes AS n ON n.id = ca.node_id")
	if nodeName != "" {
		query = query.Where("n.name = ?", nodeName)
	}

	err = query.Find(&assignments).Error
	return
}

func (api *Api) AssignCompute(tctx *logger.TraceContext,
	compute *db_model.Compute, nodeMap map[uint]*db_model.ClusterNode,
	nodeAssignmentsMap map[uint][]db_model.ComputeAssignmentWithComputeAndNode,
	assignmentsMap map[string][]db_model.ComputeAssignmentWithComputeAndNode,
	isReschedule bool) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	var rspec spec.RegionServiceComputeSpec
	if err = json_utils.Unmarshal(compute.Spec, &rspec); err != nil {
		return
	}

	policy := rspec.SchedulePolicy
	assignNodes := []uint{}
	updateNodes := []uint{}
	unassignNodes := []uint{}

	currentAssignments, ok := assignmentsMap[compute.Name]
	if ok {
		infoMsg := []string{}
		for _, currentAssignment := range currentAssignments {
			infoMsg = append(infoMsg, currentAssignment.NodeName)
		}
		logger.Infof(tctx, "currentAssignments: %v", infoMsg)
	}

	fmt.Println("DEBUG nodeAssignments: ", nodeAssignmentsMap)

	// filtering node
	enableNodeFilters := false
	if len(policy.NodeFilters) > 0 {
		enableNodeFilters = true
	}
	enableLabelFilters := false
	if len(policy.NodeLabelFilters) > 0 {
		enableLabelFilters = true
	}
	enableHardAffinites := false
	if len(policy.NodeLabelHardAffinities) > 0 {
		enableHardAffinites = true
	}
	enableHardAntiAffinites := false
	if len(policy.NodeLabelHardAntiAffinities) > 0 {
		enableHardAntiAffinites = true
	}
	enableSoftAffinites := false
	if len(policy.NodeLabelSoftAffinities) > 0 {
		enableSoftAffinites = true
	}
	enableSoftAntiAffinites := false
	if len(policy.NodeLabelSoftAntiAffinities) > 0 {
		enableSoftAntiAffinites = true
	}

	labelFilterNodeMap := map[uint]*db_model.ClusterNode{}
	filteredNodes := []*db_model.ClusterNode{}
	labelNodesMap := map[string][]*db_model.ClusterNode{} // LabelごとのNode候補
	for _, node := range nodeMap {
		labels := []string{}
		ok := true
		if enableNodeFilters {
			ok = false
			for _, nodeName := range policy.NodeFilters {
				if node.Name == nodeName {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		if enableLabelFilters {
			ok = false
			for _, label := range policy.NodeLabelFilters {
				if strings.Index(node.Labels, label) >= 0 {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		if enableHardAffinites {
			ok = false
			for _, label := range policy.NodeLabelHardAffinities {
				if strings.Index(node.Labels, label) >= 0 {
					ok = true
					labels = append(labels, label)
					break
				}
			}
			if !ok {
				continue
			}
		}

		if enableHardAntiAffinites {
			ok = false
			for _, label := range policy.NodeLabelHardAntiAffinities {
				if strings.Index(node.Labels, label) >= 0 {
					ok = true
					labels = append(labels, label)
					break
				}
			}
			if !ok {
				continue
			}
		}

		if enableSoftAffinites {
			ok = false
			for _, label := range policy.NodeLabelSoftAffinities {
				if strings.Index(node.Labels, label) >= 0 {
					ok = true
					labels = append(labels, label)
					break
				}
			}
			if !ok {
				continue
			}
		}

		if enableSoftAntiAffinites {
			ok = false
			for _, label := range policy.NodeLabelSoftAntiAffinities {
				if strings.Index(node.Labels, label) >= 0 {
					ok = true
					labels = append(labels, label)
					break
				}
			}
			if !ok {
				continue
			}
		}

		// labelFilterNodeMapには、LabelのみによるNodeのフィルタリング結果を格納する
		labelFilterNodeMap[node.ID] = node

		// Filter node by status, state
		if node.Status != base_const.StatusEnabled {
			continue
		}

		if node.State != base_const.StateUp {
			continue
		}

		// TODO
		// Filter node by cpu, memory, disk

		filteredNodes = append(filteredNodes, node)

		for _, label := range labels {
			nodes, lok := labelNodesMap[label]
			if !lok {
				nodes = []*db_model.ClusterNode{}
			}
			nodes = append(nodes, node)
			labelNodesMap[label] = nodes
		}
	}

	replicas := policy.Replicas
	if !isReschedule {
		for _, assignment := range currentAssignments {
			// labelFilterNodeMapには、LabelのみによるNodeのフィルタリング結果が格納されている
			// label変更されてNodeが候補から外された場合は、unassignNodesに追加する
			_, ok := labelFilterNodeMap[assignment.NodeID]
			if ok {
				updateNodes = append(updateNodes, assignment.NodeID)
			} else {
				unassignNodes = append(unassignNodes, assignment.NodeID)
			}
		}
		replicas = replicas - len(currentAssignments) + len(unassignNodes)
	}

	if replicas != 0 {
		for i := 0; i < replicas; i++ {
			candidates := []*db_model.ClusterNode{}
			for _, label := range policy.NodeLabelHardAntiAffinities {
				tmpCandidates := []*db_model.ClusterNode{}
				nodes := labelNodesMap[label]
				for _, node := range nodes {
					existsNode := false
					for _, n := range assignNodes {
						if node.ID == n {
							existsNode = true
							break
						}
					}
					if existsNode {
						continue
					}
					for _, n := range updateNodes {
						if node.ID == n {
							existsNode = true
							break
						}
					}
					if existsNode {
						continue
					}
					tmpCandidates = append(candidates, node)
				}
				if len(candidates) == 0 {
					candidates = tmpCandidates
				} else {
					newCandidates := []*db_model.ClusterNode{}
					for _, c := range candidates {
						for _, tc := range tmpCandidates {
							if c == tc {
								newCandidates = append(newCandidates, c)
								break
							}
						}
					}
					candidates = newCandidates
				}
			}

			for _, label := range policy.NodeLabelHardAffinities {
				tmpCandidates := []*db_model.ClusterNode{}
				nodes := labelNodesMap[label]
				if len(candidates) == 0 && len(assignNodes) == 0 && len(updateNodes) == 0 {
					for _, node := range nodes {
						tmpCandidates = append(tmpCandidates, node)
					}
					candidates = tmpCandidates
					break
				} else if len(assignNodes) > 0 {
					for _, node := range nodes {
						for _, assignNodeID := range assignNodes {
							if node.ID == assignNodeID {
								candidates = append(candidates, node)
								break
							}
						}
					}
					break
				} else if len(updateNodes) > 0 {
					for _, node := range nodes {
						for _, updateNodeID := range updateNodes {
							if node.ID == updateNodeID {
								candidates = append(candidates, node)
								break
							}
						}
					}
					break
				}
			}

			if !enableNodeFilters && !enableLabelFilters && !enableHardAffinites && !enableHardAntiAffinites {
				if len(candidates) == 0 {
					for _, node := range filteredNodes {
						candidates = append(candidates, node)
					}
				}
			}

			// candidatesのweightを調整する
			for _, label := range policy.NodeLabelSoftAffinities {
				nodes := labelNodesMap[label]
				for _, node := range nodes {
					for _, assignNodeId := range assignNodes {
						if node.ID == assignNodeId {
							node.Weight += 1000
							break
						}
					}
					for _, updateNodeId := range updateNodes {
						if node.ID == updateNodeId {
							node.Weight += 1000
							break
						}
					}
				}
			}

			for _, label := range policy.NodeLabelSoftAntiAffinities {
				nodes := labelNodesMap[label]
				for _, node := range nodes {
					for _, assignNodeId := range assignNodes {
						if node.ID == assignNodeId {
							node.Weight -= 1000
							break
						}
					}
					for _, updateNodeId := range updateNodes {
						if node.ID == updateNodeId {
							node.Weight -= 1000
							break
						}
					}
				}
			}

			if len(candidates) == 0 {
				break
			}

			for _, candidate := range candidates {
				assignments, ok := nodeAssignmentsMap[candidate.ID]
				if ok {
					for _, assignment := range assignments {
						// TODO calucurate assignment.Cost before scheduling
						candidate.Weight -= (10 + assignment.Cost)
					}
				}
			}

			// TODO Sort candidates by weight

			assignNodes = append(assignNodes, candidates[0].ID)
		}
	}

	if policy.Replicas != len(assignNodes)+len(updateNodes)-len(unassignNodes) {
		logger.Warningf(tctx, "Failed assign: compute=%v", compute.Name)
		return
	}

	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		for _, nodeID := range updateNodes {
			switch compute.Status {
			case base_const.StatusUpdating:
				tx.Create(&db_model.ComputeAssignment{
					ComputeID:    compute.ID,
					NodeID:       nodeID,
					Status:       base_const.StatusUpdating,
					StatusReason: "Updating",
				})
			}
		}

		for _, nodeID := range assignNodes {
			switch compute.Status {
			case base_const.StatusInitializing:
				tx.Create(&db_model.ComputeAssignment{
					ComputeID:    compute.ID,
					NodeID:       nodeID,
					Status:       base_const.StatusCreating,
					StatusReason: "Creating",
				})
			}
		}

		switch compute.Status {
		case base_const.StatusInitializing:
			compute.Status = base_const.StatusCreatingScheduled
			compute.StatusReason = "CreatingScheduled"
		}
		return
	})
}

func (api *Api) ConfirmCreatingScheduledCompute(tctx *logger.TraceContext,
	compute *db_model.Compute,
	assignmentsMap map[string][]db_model.ComputeAssignmentWithComputeAndNode) {
	var err error
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 1) }()

	assignments, ok := assignmentsMap[compute.Name]
	if !ok {
		err = error_utils.NewConflictNotFoundError(compute.Name)
		return
	}

	existsNonActiveAssignments := false
	for _, assignment := range assignments {
		if assignment.Status != base_const.StatusActive {
			existsNonActiveAssignments = true
			break
		}
	}

	if existsNonActiveAssignments {
		logger.Info(tctx, "Waiting: exists non active assignments")
		return
	}

	err = api.Transact(tctx, func(tx *gorm.DB) (err error) {
		tmpCompute := resource_model.Compute{
			Status:       resource_model.StatusActive,
			StatusReason: "ConfirmedCreagingScheduled",
		}
		err = tx.Model(&tmpCompute).Where("id = ?", compute.ID).Updates(&tmpCompute).Error
		return
	})
}