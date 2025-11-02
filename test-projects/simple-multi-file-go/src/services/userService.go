package services

import (
	"github.com/test/simple-multi-file/models"
)

func GetUserService() string {
	user := CreateUser("Alice", "alice@example.com")
	return /* unsupported expression */
}

func GetAllUsers() []User {
	return []interface{}{CreateUser("Alice", "alice@example.com"), CreateUser("Bob", "bob@example.com")}
}

func main() {
}
