package kernel

import "net/http"

const Key = "goweb:kernel"

type Kernel interface {
	HttpEngine() http.Handler
}
