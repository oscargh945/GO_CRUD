package entities

type User struct {
	Id          string `json:"_id" bson:"_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	SoftDeleted bool   `json:"soft_deleted"`
}

type CreateUserInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	SoftDeleted bool   `json:"softDeleted" default:"false"`
}
