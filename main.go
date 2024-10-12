package main

func init() {
	EnsureDB()
}

func main() {
	RunApp()
	defer store.db.Close()
}
