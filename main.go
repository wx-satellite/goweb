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

//func main() {
//	http.HandleFunc("/haha", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Println(request.RequestURI)
//		fmt.Println(request.URL.Path)
//		panic("123")
//	})
//	http.HandleFunc("/name", func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Println(request.RequestURI)
//		fmt.Println(request.URL.Path)
//	})
//	http.ListenAndServe(":8080", nil)
//}
