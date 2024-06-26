package main

import "log"

func main() {
	app := App{}
	err := app.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize app: ", err)
	}
	app.Run("localhost:3000")
}
