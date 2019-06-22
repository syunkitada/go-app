package resource_model

import (
	"github.com/jinzhu/gorm"
	"github.com/syunkitada/goapp/pkg/authproxy/index_model"
)

const RegionServiceKind = "RegionService"

type RegionService struct {
	gorm.Model
	Region string `gorm:"not null;size:50;"`
	Name   string `gorm:"not null;size:255;"` // Vip Domain
	Kind   string `gorm:"not null;size:25;"`
}

type RegionServiceSpec struct {
	Name   string `validate:"required"`
	Region string `validate:"required"`
	Kind   string `validate:"required"`
}

var RegionServiceCmd map[string]index_model.Cmd = map[string]index_model.Cmd{
	"create_region-service": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeFile,
		ArgKind: RegionServiceKind,
		Help:    "helptext",
	},
	"update_region-service": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeFile,
		ArgKind: RegionServiceKind,
		Help:    "helptext",
	},
	"get_region-services": index_model.Cmd{
		Arg:         index_model.ArgOptional,
		ArgType:     index_model.ArgTypeString,
		ArgKind:     RegionServiceKind,
		Help:        "helptext",
		TableHeader: []string{"Name", "Kind", "Region"},
	},
	"get_region-service": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeString,
		ArgKind: RegionServiceKind,
		Help:    "helptext",
	},
	"delete_region-service": index_model.Cmd{
		Arg:     index_model.ArgRequired,
		ArgType: index_model.ArgTypeString,
		ArgKind: RegionServiceKind,
		Help:    "helptext",
	},
}