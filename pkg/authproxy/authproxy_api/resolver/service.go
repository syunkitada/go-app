package resolver

import (
	"github.com/syunkitada/goapp/pkg/base/base_const"
	"github.com/syunkitada/goapp/pkg/base/base_index_model"
	"github.com/syunkitada/goapp/pkg/base/base_spec"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (resolver *Resolver) GetServiceIndex(tctx *logger.TraceContext, input *base_spec.GetServiceIndex, user *base_spec.UserAuthority) (data *base_spec.GetServiceIndexData, code uint8, err error) {
	cmdMap := map[string]base_index_model.Cmd{}
	cmdMaps := []map[string]base_index_model.Cmd{
		base_spec.UserCmd,
	}
	for _, tmpCmdMap := range cmdMaps {
		for key, cmd := range tmpCmdMap {
			cmdMap[key] = cmd
		}
	}

	code = base_const.CodeOk
	data = &base_spec.GetServiceIndexData{
		Index: base_index_model.Index{
			CmdMap: cmdMap,
		},
	}

	return
}

func (resolver *Resolver) GetServiceDashboardIndex(tctx *logger.TraceContext, input *base_spec.GetServiceDashboardIndex, user *base_spec.UserAuthority) (data *base_spec.GetServiceDashboardIndexData, code uint8, err error) {
	switch input.Name {
	case "Home":
		data = &base_spec.GetServiceDashboardIndexData{
			Data: map[string]interface{}{
				"User": user,
			},
			Index: base_index_model.DashboardIndex{
				DefaultRoute: map[string]interface{}{
					"Path": []string{"User", "View"},
				},
				View: base_index_model.Panels{
					Name: "Root",
					Kind: "Panels",
					Children: []interface{}{
						map[string]interface{}{
							"Name": "User",
							"Kind": "Tabs",
							"Children": []interface{}{
								map[string]interface{}{
									"Name":    "View",
									"Kind":    "View",
									"DataKey": "User",
									"PanelsGroups": []interface{}{
										map[string]interface{}{
											"Name": "Detail",
											"Kind": "Cards",
											"Cards": []interface{}{
												map[string]interface{}{
													"Name": "Detail",
													"Kind": "Fields",
													"Fields": []base_index_model.Field{
														base_index_model.Field{Name: "Name"},
													},
												},
											},
										},
									},
								},
								map[string]interface{}{
									"Name": "Password Setting",
									"Kind": "Form",
								},
							},
						},
					},
				},
			},
		}

	default:
		code = base_const.CodeClientNotFound
	}

	return
}

func (resolver *Resolver) GetProjectServiceDashboardIndex(tctx *logger.TraceContext, input *base_spec.GetServiceDashboardIndex, user *base_spec.UserAuthority) (data *base_spec.GetServiceDashboardIndexData, code uint8, err error) {
	switch input.Name {
	case "HomeProject":
		data = &base_spec.GetServiceDashboardIndexData{
			Data: map[string]interface{}{
				"User": user,
			},
			Index: base_index_model.DashboardIndex{
				DefaultRoute: map[string]interface{}{
					"Path": []string{"User", "View"},
				},
				View: base_index_model.Panels{
					Name: "Root",
					Kind: "Panels",
					Children: []interface{}{
						map[string]interface{}{
							"Name": "User",
							"Kind": "Tabs",
							"Children": []interface{}{
								map[string]interface{}{
									"Name":    "View",
									"Kind":    "View",
									"DataKey": "User",
									"PanelsGroups": []interface{}{
										map[string]interface{}{
											"Name": "Detail",
											"Kind": "Cards",
											"Cards": []interface{}{
												map[string]interface{}{
													"Name": "Detail",
													"Kind": "Fields",
													"Fields": []base_index_model.Field{
														base_index_model.Field{Name: "Name"},
														base_index_model.Field{Name: "ProjectName"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	default:
		code = base_const.CodeClientNotFound
	}

	return
}
