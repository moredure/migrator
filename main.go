package main

func main() {
	migrator, err := initializeApp()
	if err != nil {
		panic(err)
	}
	migrator.Migrate()
}



