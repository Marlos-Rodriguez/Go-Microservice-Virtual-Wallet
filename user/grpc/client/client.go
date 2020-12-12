package client

//StartClient Start the client for all gRPC clients
func StartClient() {
	go startMoveClient()
	startAuthClient()
}

//CloseClient Close the client of all gRPC clients
func CloseClient() {
	closeAuthClient()
	closeMoveClient()
}
