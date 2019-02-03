package main

import "log"

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Panicf("failed to initialize app %v\n", err)
	}
	app.Migrate()
}
