package main

import (
	"goweb/framework"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}
	panic(server.ListenAndServe())
}
