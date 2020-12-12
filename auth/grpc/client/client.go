package client

//StartClient Start all gRPC clients
func StartClient() {
	startMoveClient()
}

//CloseClient Close all gRPC clients
func CloseClient() {
	closeMoveClient()
}
