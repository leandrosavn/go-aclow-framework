package modeltype

type User struct {
	Id        string `bson:"_id"`
	User      string `json:"user"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Profile   string `json:"profile"`
	Type      string `json:"type"`
}
