package models

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(name string, email string) User {
	return User{Id: GenerateId(), Name: name, Email: email}
}

func GenerateId() string {
	return MathRandom().ToString(36).Substr(2, 9)
}

