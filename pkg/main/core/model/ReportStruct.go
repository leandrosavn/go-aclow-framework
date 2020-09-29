package model

type ValidateDate struct {
	GroupName      string `json:"groupName"`
	MaxDaysBetween int    `json:"maxDaysBetween"`
}

type ParamOptions struct {
	Descricao string `json:"descricao"`
	Valor     string `json:"valor"`
}

type ParamsReport struct {
	Name              string         `json:"name"`
	Tipo              string         `json:"tipo"`
	Rules             string         `json:"rules"`
	Options           []ParamOptions `json:"options"`
	QueryId           string         `json:"queryId"`
	Query             string         `json:"query"`
	CampoId           string         `json:"campoId"`
	CampoDesc         string         `json:"campoDesc"`
	Entidade          string         `json:"entidade"`
	DateGroupValidate string         `json:"dateGroupValidate"`
	DateOrder         int            `json:"dateOrder"`
}

type Report struct {
	Id            interface{}    `json:"id"`
	Name          string         `json:"name"`
	Schema        string         `json:"schema"`
	Database      string         `json:"database"`
	ProcName      string         `json:"procName"`
	ExecutionType string         `json:"executionType"`
	Params        []ParamsReport `json:"params"`
	ParamsValues  []interface{}  `json:"paramsValues"`
	Validate      []ValidateDate `json:"validate"`
}
