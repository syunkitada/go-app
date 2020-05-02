package spec

import (
	"time"

	"github.com/syunkitada/goapp/pkg/base/base_index_model"
)

type Compute struct {
	Region        string
	Cluster       string
	RegionService string
	Name          string
	Kind          string
	Labels        string
	Status        string
	StatusReason  string
	Project       string
	Spec          interface{}
	LinkSpec      string
	Image         string
	Vcpus         uint
	Memory        uint
	Disk          uint
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

type GetCompute struct {
	Name   string `validate:"required"`
	Region string `validate:"required"`
}

type GetComputeData struct {
	Compute Compute
}

type GetComputeConsole struct {
	Name   string `validate:"required"`
	Region string `validate:"required"`
}

type GetComputeConsoleData struct {
	Compute Compute
}

type WsComputeConsoleInput struct {
	Bytes []byte
}

type WsComputeConsoleOutput struct {
	Bytes []byte
}

type GetComputes struct {
	Region string `validate:"required"`
}

type GetComputesData struct {
	Computes []Compute
}

type CreateCompute struct {
	Spec string `validate:"required" flagKind:"file"`
}

type CreateComputeData struct{}

type UpdateCompute struct {
	Spec string `validate:"required" flagKind:"file"`
}

type UpdateComputeData struct{}

type DeleteCompute struct {
	Name string `validate:"required"`
}

type DeleteComputeData struct{}

type DeleteComputes struct {
	Spec string `validate:"required" flagKind:"file"`
}

type DeleteComputesData struct{}

var ComputesTable = base_index_model.Table{
	Name:        "Computes",
	Kind:        "Table",
	DataQueries: []string{"GetComputes"},
	DataKey:     "Computes",
	SelectActions: []base_index_model.Action{
		base_index_model.Action{
			Name:      "Delete",
			Icon:      "Delete",
			Kind:      "Form",
			DataKind:  "Compute",
			SelectKey: "Name",
		},
	},
	Columns: []base_index_model.TableColumn{
		base_index_model.TableColumn{
			Name: "Name", IsSearch: true,
			Align:    "left",
			LinkPath: []string{"Resource", "Compute", "View"},
			LinkKey:  "Name",
		},
		base_index_model.TableColumn{Name: "RegionService"},
		base_index_model.TableColumn{Name: "Kind"},
		base_index_model.TableColumn{Name: "Status"},
		base_index_model.TableColumn{Name: "StatusReason"},
		base_index_model.TableColumn{Name: "UpdatedAt", Kind: "Time"},
		base_index_model.TableColumn{Name: "CreatedAt", Kind: "Time"},
	},
}

var ComputesDetail = base_index_model.Tabs{
	Name:   "Compute",
	Kind:   "RouteTabs",
	IsSync: true,
	Children: []interface{}{
		base_index_model.View{
			Name:        "View",
			Kind:        "View",
			DataKey:     "Compute",
			DataQueries: []string{"GetCompute"},
			PanelsGroups: []interface{}{
				map[string]interface{}{
					"Name": "Detail",
					"Kind": "Cards",
					"Cards": []interface{}{
						map[string]interface{}{
							"Name": "Detail",
							"Kind": "Fields",
							"Fields": []base_index_model.Field{
								base_index_model.Field{Name: "Name", Kind: "text"},
								base_index_model.Field{Name: "Kind", Kind: "select"},
							},
						},
					},
				},
			},
		},
		base_index_model.Form{
			Name:         "Edit",
			Kind:         "Form",
			DataKey:      "Compute",
			DataQueries:  []string{"GetCompute"},
			SubmitAction: "update compute",
			Icon:         "Update",
			Fields: []base_index_model.Field{
				base_index_model.Field{Name: "Name", Kind: "text", Require: true,
					Updatable: false,
					Min:       5, Max: 200, RegExp: "^[0-9a-zA-Z]+$",
					RegExpMsg: "Please enter alphanumeric characters."},
				base_index_model.Field{Name: "Kind", Kind: "select", Require: true,
					Updatable: true,
					Options: []string{
						"Server", "Pdu", "RackSpineRouter",
						"FloorLeafRouter", "FloorSpineRouter", "GatewayRouter",
					}},
			},
		},
		base_index_model.View{
			Name:            "Console",
			Route:           "/Console",
			Kind:            "View",
			DataKey:         "Compute",
			DataQueries:     []string{"GetComputeConsole"},
			EnableWebSocket: true,
			WebSocketKey:    "ComputeConsole",
			PanelsGroups: []interface{}{
				map[string]interface{}{
					"Name": "Console",
					"Kind": "Cards",
					"Cards": []interface{}{
						map[string]interface{}{
							"Name": "Console",
							"Kind": "Console",
						},
					},
				},
			},
		},
	},
}
