package api

//Start Start a new User server API
func Start() {
	app := routes()
	createServer(app)
}
