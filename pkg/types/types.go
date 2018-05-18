package types

type ModeReplicated struct {
	Replicas int
}

type Mode struct {
	Replicated ModeReplicated
}

type UpdateConfig struct {
	Parallelism     int
	Delay           int
	FailureAction   string
	Monitor         int
	MaxFailureRatio int
	Order           string
}

type RollbackConfig struct {
	Parallelism     int
	FailureAction   string
	Monitor         int
	MaxFailureRatio int
	Order           string
}

type Network struct {
	Target string
}

type EndpointSpec struct {
	Mode string
}

type ContainerSpec struct {
	Image string
}

type TaskTemplate struct {
	ContainerSpec ContainerSpec
	ForceUpdate   int
	Runtime       string
}

// Docker service spec struct
type Spec struct {
	Name           string
	Labels         map[string]string
	TaskTemplate   TaskTemplate
	Mode           Mode
	UpdateConfig   UpdateConfig
	RollbackConfig RollbackConfig
	Networks       []Network
	EndpointSpec   EndpointSpec
}

type ServiceVersion struct {
	Index int
}

// Docker service struct
type Service struct {
	Spec    Spec
	Version ServiceVersion
}

type ServiceUpdateQueryString struct {
	version int
}

// CreateServiceResponse format of docker api response when create a service
type CreateServiceResponse struct {
	message string
	ID      string
}
