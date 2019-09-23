package spec_model

type Spec struct {
	Meta     interface{}
	Name     string
	Apis     []Api
	QuerySet map[string]Query
}

type Api struct {
	Name            string
	Cmds            map[string]string
	RequiredAuth    bool
	RequiredProject bool
	RequiredService bool
	Queries         []Query
	QueryModels     []QueryModel
	ViewModels      []interface{}
}
