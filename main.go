package main

import (
	"github.com/jonasleonhard/go-htmx-time/src/server"
)

func main() {
	server := server.New()
	defer server.Close()

	err := server.ListenAndServe()

	if err != nil {
		panic("cannot start server")
	}
}
