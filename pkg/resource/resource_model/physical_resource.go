package resource_model

import (
	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/authproxy/index_model"
)

const PhysicalResourceKind = "PhysicalResource"

type PhysicalResource struct {
	gorm.Model
	Datacenter    string `gorm:"not null;size:50;index;"`
	Rack          string `gorm:"not null;size:50;index;"`
	Cluster       string `gorm:"not null;size:50;index;"`  // 仮想リソースを扱う場合はClusterに紐図かせる
	Name          string `gorm:"not null;size:200;index;"` // Datacenter内でユニーク
	Kind          string `gorm:"not null;size:25;"`        // Server, Pdu, L2Switch, L3Switch, RootSwitch
	PhysicalModel string `gorm:"not null;size:200;"`
	RackPosition  uint8  `gorm:"not null;"`
	Status        string `gorm:"not null;size:25;"`
	StatusReason  string `gorm:"not null;size:50;"`
	PowerLinks    string `gorm:"not null;size:5000;"`
	NetLinks      string `gorm:"not null;size:5000;"`
	Spec          string `gorm:"not null;size:5000;"`
}

type PhysicalResourceSpec struct {
	Kind         string `validate:"required"`
	Name         string `validate:"required"`
	Datacenter   string `validate:"required"`
	Cluster      string
	Rack         string
	Model        string
	RackPosition uint8
	NetLinks     []string
	PowerLinks   []string
	Spec         string
}

var PhysicalResourceCmd map[string]index_model.Cmd = map[string]index_model.Cmd{
	"CreatePhysicalResource": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeFile,
		ArgKind: PhysicalResourceKind,
		Help:    "helptext",
	},
	"UpdatePhysicalResource": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeFile,
		ArgKind: PhysicalResourceKind,
		Help:    "helptext",
	},
	"GetPhysicalResources": index_model.Cmd{
		Arg:         index_model.ArgOptional,
		ArgType:     index_model.ArgTypeString,
		ArgKind:     PhysicalResourceKind,
		Help:        "helptext",
		TableHeader: []string{"Name", "Kind", "Datacenter"},
		FlagMap: map[string]index_model.Flag{
			"datacenter": index_model.Flag{
				Flag:     index_model.ArgRequired,
				FlagType: index_model.ArgTypeString,
				Help:     "datacenter",
			},
		},
	},
	"GetPhysicalResource": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeString,
		ArgKind: PhysicalResourceKind,
		Help:    "helptext",
	},
	"DeletePhysicalResource": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeString,
		ArgKind: PhysicalResourceKind,
		Help:    "helptext",
	},
}

var PhysicalResourcesTable = index_model.Table{
	Name:    "Resources",
	Route:   "PhysicalResources",
	Kind:    "Table",
	DataKey: "PhysicalResources",
	Actions: []index_model.Action{
		index_model.Action{
			Name: "Create", Icon: "Create", Kind: "Form",
			DataKind: "PhysicalResource",
			Fields: []index_model.Field{
				index_model.Field{Name: "Name", Kind: "text",
					Require: true, Min: 5, Max: 200, RegExp: "^[0-9a-zA-Z]+$"},
				index_model.Field{Name: "Kind", Kind: "select", Require: true,
					Options: []string{
						"Server", "Pdu", "RackSpineRouter",
						"FloorLeafRouter", "FloorSpineRouter", "GatewayRouter",
					}},
				index_model.Field{Name: "Rack", Kind: "select", Require: true,
					DataKey: "Racks"},
				index_model.Field{Name: "Model", Kind: "select", Require: true,
					DataKey: "PhysicalModels"},
			},
		},
	},
	SelectActions: []index_model.Action{
		index_model.Action{Name: "Delete", Icon: "Delete",
			Kind:      "Form",
			DataKind:  "PhysicalResource",
			SelectKey: "Name",
		},
	},
	ColumnActions: []index_model.Action{
		index_model.Action{Name: "Detail", Icon: "Detail"},
		index_model.Action{Name: "Update", Icon: "Update"},
	},
	Columns: []index_model.TableColumn{
		index_model.TableColumn{
			Name: "Name", IsSearch: true,
			Link:           "Datacenters/:datacenter/Resources/Resources/Detail/:0/View",
			LinkParam:      "resource",
			LinkSync:       false,
			LinkGetQueries: []string{"GetPhysicalResource"},
		},
		index_model.TableColumn{Name: "Kind"},
		index_model.TableColumn{Name: "UpdatedAt", Kind: "Time"},
		index_model.TableColumn{Name: "CreatedAt", Kind: "Time"},
		index_model.TableColumn{Name: "Action", Kind: "Action"},
	},
}
