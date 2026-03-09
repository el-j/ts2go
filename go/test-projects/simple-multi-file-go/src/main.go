package main

import (
	"github.com/test/simple-multi-file/services"
)

func Main() interface{} {
	fmt.Println("Starting application...")
	userInfo := GetUserService()
	fmt.Println(userInfo)
	users := GetAllUsers()
	fmt.Println( /* unsupported expression */ )
}

func main() {
	Main()
}
