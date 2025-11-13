package main

func FetchData() chan Promise {
	resultCh := make(chan interface{}, 1)
	go func() {
		return
	}()
	return resultCh
}

func main() {
	app := express()
	PORT := 3000
	app.Use(express.Json())
	app.Get("/", func(req interface{}, res interface{}) interface{} { res.Send("Hello, World!") })
	app.Get("/user/:id", func(req interface{}, res interface{}) interface{} {
		userId := req.Params.Id
		res.Json(map[string]interface{}{"UserId": userId, "Name": "John Doe"})
	})
	app.Post("/api/data", func(req interface{}, res interface{}) interface{} {
		data := req.Body
		res.Json(map[string]interface{}{"Success": true, "Received": data})
	})
	app.Get("/api/async", func(req interface{}, res interface{}) interface{} {
		func() {
			defer func() {
				if error := recover(); error != nil {
					res.Status(500).Json(map[string]interface{}{"Error": error.Message})
				}
			}()
			result := (<-fetchData())
			res.Json(map[string]interface{}{"Data": result})
		}()
	})
	app.Listen(PORT, func() interface{} { fmt.Println("Server running on port " + fmt.Sprint(PORT)) })
}
