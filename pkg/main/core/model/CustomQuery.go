package model

type CustomQuery struct {
	Id        string `bson:"_id"`
	Name      string `json:"name"`
	Tipo      string `json:"tipo"`
	CampoId   string `json:"campoId"`
	CampoDesc string `json:"campoDesc"`
	Query     string `json:"query"`
}
