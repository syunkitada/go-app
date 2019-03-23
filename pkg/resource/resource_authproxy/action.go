package resource_authproxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (resource *Resource) PhysicalAction(c *gin.Context) {
	tctx, err := logger.NewAuthproxyActionTraceContext(resource.host, resource.name, c)
	startTime := logger.StartTrace(&tctx.TraceContext)
	defer func() { logger.EndTrace(&tctx.TraceContext, startTime, err, 1) }()

	if err != nil {
		c.JSON(500, gin.H{
			"err": "InvalidRequest",
		})
		return
	}

	resp, err := resource.resourceApiClient.Action(tctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"TraceId": tctx.TraceId,
			"Err":     err,
		})
		return
	}

	c.JSON(200, gin.H{
		"Data": resp,
		"Index": gin.H{
			"Name": "Root",
			"Kind": "RoutePanels",
			"Panels": []interface{}{
				gin.H{
					"Name":    "Datacenters",
					"Route":   "",
					"Subname": "datacenter",
					"Kind":    "Table",
					"DataKey": "Datacenters",
					"Columns": []interface{}{
						gin.H{"Name": "Name", "IsSearch": true, "Link": "/Datacenters"},
						gin.H{"Name": "Region", "IsSearch": true},
						gin.H{"Name": "UpdatedAt", "Type": "Time"},
						gin.H{"Name": "CreatedAt", "Type": "Time"},
					},
				},
				gin.H{
					"Name":       "Resources",
					"Route":      "/Datacenters/:datacenter",
					"Kind":       "RouteTabs",
					"GetQueries": []string{"GetPhysicalResources", "GetRacks", "GetFloors"},
					"Tabs": []interface{}{
						gin.H{
							"Name":    "Resources",
							"Route":   "",
							"Kind":    "Table",
							"DataKey": "PhysicalResources",
							"Columns": []interface{}{
								gin.H{"Name": "Name", "IsSearch": true},
								gin.H{"Name": "Kind"},
								gin.H{"Name": "UpdatedAt", "Type": "Time"},
								gin.H{"Name": "CreatedAt", "Type": "Time"},
							},
						},
						gin.H{
							"Name":    "Racks",
							"Route":   "/Racks",
							"Kind":    "Table",
							"DataKey": "Racks",
							"Columns": []interface{}{
								gin.H{"Name": "Name", "IsSearch": true},
								gin.H{"Name": "Kind"},
								gin.H{"Name": "UpdatedAt", "Type": "Time"},
								gin.H{"Name": "CreatedAt", "Type": "Time"},
							},
						},
						gin.H{
							"Name":    "Floors",
							"Route":   "/Floors",
							"Kind":    "Table",
							"DataKey": "Floors",
							"Columns": []interface{}{
								gin.H{"Name": "Name", "IsSearch": true},
								gin.H{"Name": "Kind"},
								gin.H{"Name": "UpdatedAt", "Type": "Time"},
								gin.H{"Name": "CreatedAt", "Type": "Time"},
							},
						},
					},
				},
			},
		},
	})
}

func (resource *Resource) VirtualAction(c *gin.Context) {
	tctx, err := logger.NewAuthproxyActionTraceContext(resource.host, resource.name, c)
	startTime := logger.StartTrace(&tctx.TraceContext)
	defer func() { logger.EndTrace(&tctx.TraceContext, startTime, err, 1) }()

	if err != nil {
		c.JSON(500, gin.H{
			"err": "InvalidRequest",
		})
		return
	}

	resp, err := resource.resourceApiClient.Action(tctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"TraceId": tctx.TraceId,
			"Err":     err,
		})
		return
	}

	c.JSON(200, resp)
}
