package main

func main() {

	server := NewServer(":3000")
	server.Handle("GET", "/", HandleRoot)
	server.Handle("POST", "/api", server.AddMiddleware(HandleHome, CheckAuth(), logging()))
	server.Handle("POST", "/create", PostRequest)
	server.Handle("POST", "/trx", TrxPostRequest)
	server.Listen()
}

func Add(a, b int) int {
	return a + b
}
