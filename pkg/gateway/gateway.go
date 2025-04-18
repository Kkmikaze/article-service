package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Gateway struct {
	*runtime.ServeMux
	Addr           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	Http           *http.ServeMux
}

//go:generate mockery --name GatewayInterface --config .mockery.yaml
type GatewayInterface interface {
	Run(ctx context.Context, handler http.Handler) error
	GetMux() *http.ServeMux
	GetRuntimeMux() *runtime.ServeMux
}

type ErrorResponse struct {
	Code    uint32            `json:"code"`
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (g *Gateway) Run(ctx context.Context, handler http.Handler) error {
	host := g.Addr
	if !strings.Contains(host, ":") {
		host = fmt.Sprintf(":%v", host)
	}

	gwServer := &http.Server{
		Addr:           host,
		Handler:        handler,
		ReadTimeout:    g.ReadTimeout,
		WriteTimeout:   g.WriteTimeout,
		MaxHeaderBytes: g.MaxHeaderBytes,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down the HTTP gateway server...")
		if err := gwServer.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown HTTP gateway server: %v\n", err)
		}
	}()

	err := gwServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("gRPC-Gateway server exited with error: %v\n", err)
	}

	return nil
}

func (g *Gateway) GetMux() *http.ServeMux {
	return g.Http
}

func (g *Gateway) GetRuntimeMux() *runtime.ServeMux {
	return g.ServeMux
}

func NewGateway(
	addr string,
	maxHeaderBytes int,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	opts ...runtime.ServeMuxOption,
) GatewayInterface {
	runtimeMux := runtime.NewServeMux(opts...)
	httpMux := http.NewServeMux()

	return &Gateway{
		ServeMux:       runtimeMux,
		Addr:           addr,
		Http:           httpMux,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}
}

func ExceptionHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s := status.Convert(err)
	code := runtime.HTTPStatusFromCode(s.Code())
	fallback := []byte(`{"code":500,"status":"Internal Server Error","message":"Internal Server Error","errors":{}}`)
	w.Header().Set("Content-Type", m.ContentType("application/json"))
	w.WriteHeader(code)

	objectMapper := make(map[string]string)
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.FieldViolations {
				objectMapper[strings.ToLower(violation.GetField())] = violation.GetDescription()
			}
		}
	}

	response := ErrorResponse{
		Code:    uint32(code),
		Status:  http.StatusText(code),
		Message: s.Message(),
		Errors:  objectMapper,
	}

	marshaled, marshalErr := json.Marshal(&response)
	if marshalErr != nil {
		_, _ = w.Write(fallback)
		return
	}

	_, _ = w.Write(marshaled)
}
