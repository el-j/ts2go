package main

func FetchUrl(url string) chan Promise {
	resultCh := make(chan interface{}, 1)
	go func() {
		return
	}()
	return resultCh
}

func main() {
	program := NewCommand()
	program.Name("mytool").Description("A sample CLI tool").Version("1.0.0")
	program.Command("greet <name>").Description("Greet someone").Option("-l, --loud", "Greet loudly").Action(func(name interface{}, options interface{}) interface{} {
		greeting := "Hello, " + fmt.Sprint(name) + "!"
		if options.Loud {
			fmt.Println(greeting)
		}
	})
	program.Command("fetch <url>").Description("Fetch data from URL").Action(func(url interface{}) interface{} {
		func() {
			defer func() {
				if error := recover(); error != nil {
					console.Error("Error:", error)
				}
			}()
			data := (<-fetchUrl(url))
			fmt.Println("Data:", data)
		}()
	})
	program.Parse(process.Argv)
}
