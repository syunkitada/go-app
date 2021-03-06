package resource_model

const (
	ResourceKindDatacenter       = "Datacenter"
	ResourceKindCluster          = "Cluster"
	ResourceKindFloor            = "Floor"
	ResourceKindRack             = "Rack"
	ResourceKindPhysicalModel    = "PhysicalModel"
	ResourceKindPhysicalResource = "PhysicalResource"
	ResourceKindCompute          = "Compute"
)

const (
	KindResourceController        = "ResourceController"
	KindResourceApi               = "ResourceApi"
	KindResourceClusterController = "ResourceClusterController"
	KindResourceClusterApi        = "ResourceClusterApi"
	KindResourceClusterAgent      = "ResourceClusterAgent"
)

const (
	StatusEnabled   = "Enabled"
	StatusEnabling  = "Enabling"
	StatusDisabled  = "Disabled"
	StatusDisabling = "Disabling"

	StatusInitializing          = "Initializing"
	StatusActive                = "Active"
	StatusCreating              = "Creating"
	StatusCreatingInitialized   = "Creating:Initialized"
	StatusCreatingScheduled     = "Creating:Scheduled"
	StatusUpdating              = "Updating"
	StatusUpdatingScheduled     = "Updating:Scheduled"
	StatusUnknownActivating     = "Unknown:Activating"
	StatusRebalancingUnassigned = "Rebalancing:Unassigned"
	StatusDeleting              = "Deleting"
	StatusDeletingScheduled     = "Deleting:Scheduled"
	StatusDeleted               = "Deleted"
	StatusError                 = "Error"
)

const (
	StatusMsgInitializing                  = "Initializing"
	StatusMsgInitializeErrorNoValidCluster = "InitializeError: NoValidCluster"
	StatusMsgInitializeSuccess             = "InitializeSuccess"
	StatusMsgUpdating                      = "Updating"
	StatusMsgCreating                      = "Creating"
)

const (
	StateUp      = "Up"
	StateDown    = "Down"
	StateUnknown = "Unknown"
)

const (
	RoleLeader = "Leader"
	RoleMember = "Member"
)

const (
	SpecNetworkV4    = "NetworkV4"
	SpecCompute      = "Compute"
	SpecContainer    = "Container"
	SpecImage        = "Image"
	SpecVolume       = "Volume"
	SpecLoadbalancer = "Loadbalancer"
)

const (
	SpecKindNetworkV4Local  = "Local"
	SpecKindComputeLibvirt  = "Libvirt"
	SpecKindContainerDocker = "Docker"
	SpecKindImageUrl        = "Url"
	SpecKindVolumeNfs       = "Nfs"
	SpecKindVolumeIscsi     = "Iscsi"
	SpecKindLoadbalancerVpp = "Vpp"
)
