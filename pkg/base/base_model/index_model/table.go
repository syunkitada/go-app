package index_model

type Table struct {
	Name          string
	Kind          string
	Route         string
	Subname       string
	DataKey       string
	Columns       []TableColumn
	Actions       []Action
	SelectActions []Action
	ColumnActions []Action
}

type TableColumn struct {
	Name           string
	Kind           string
	IsSearch       bool
	Link           string
	LinkParam      string
	LinkSync       bool
	LinkGetQueries []string
	RowColoringMap map[string]string
	FilterValues   []map[string]string
	Sort           string
}

type Action struct {
	Name      string
	Kind      string
	Icon      string
	DataKind  string
	SelectKey string
	Query     string
	Fields    []Field
}

type Field struct {
	Name      string
	Kind      string
	Require   bool
	Min       uint
	Max       uint
	RegExp    string
	RegExpMsg string
	Options   []string
	DataKey   string
	Updatable bool
}
