package common

import (
	"article-service/api"
	"article-service/pkg/gateway"
	"article-service/third_party"
	"io/fs"
	"mime"
	"net/http"
)

func EnableSwagger(g gateway.GatewayInterface) {
	fileServer := http.FileServer(http.FS(api.FS))

	mux := g.GetMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", swaggerHandler())
}

func swaggerHandler() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		panic(err)
	}

	subFS, err := fs.Sub(third_party.SwaggerUI, "swagger-ui")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}

	return http.FileServer(http.FS(subFS))
}
