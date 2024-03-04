package api

import (
	"context"
	swagger "github.com/davidebianchi/gswagger"
	"github.com/davidebianchi/gswagger/support/gorilla"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"net/http"
	"usercrud/internal/entity"
)

const v = 1

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) Init(ctx context.Context, muxRouter *mux.Router) {
	router, _ := swagger.NewRouter(gorilla.NewRouter(muxRouter), swagger.Options{
		Context: ctx,
		Openapi: &openapi3.T{
			Info: &openapi3.Info{
				Title:   "my title",
				Version: "1.0.0",
			},
		},
	})
	//muxRouter.Handle("/metrics", promhttp.Handler())
	//promhttp.HandlerFor()
	//promhttp.Handler()
	//router.AddRoute(http.MethodGet, "/metrics", promhttp.Handler, swagger.Definitions{})

	router.AddRoute(http.MethodGet, "/", s.IndexHandler, swagger.Definitions{})
	router.AddRoute(http.MethodGet, "/list", s.IndexHandler, swagger.Definitions{})

	router.AddRoute(http.MethodGet, "/health", s.healthHandler, swagger.Definitions{})

	router.AddRoute(http.MethodPost, "/user", s.userPostHandler, swagger.Definitions{
		RequestBody: &swagger.ContentValue{
			Content: swagger.Content{
				"application/json": {Value: entity.User{}},
			},
		},
		Responses: map[int]swagger.ContentValue{
			201: {
				Content: swagger.Content{
					"text/html": {Value: ""},
				},
			},
			401: {
				Content: swagger.Content{
					"application/json": {Value: &errorResponse{}},
				},
				Description: "invalid request",
			},
		},
	})
	router.AddRoute(http.MethodGet, "/user/{userId}", s.userGetHandler, swagger.Definitions{
		Querystring: swagger.ParameterValue{
			"query": {
				Schema: &swagger.Schema{Value: ""},
			},
		},
	})

	router.AddRoute(http.MethodPut, "/user/{userId}", s.userPutHandler, swagger.Definitions{
		Querystring: swagger.ParameterValue{
			"query": {
				Schema: &swagger.Schema{Value: ""},
			},
		},
		RequestBody: &swagger.ContentValue{
			Content: swagger.Content{
				"application/json": {Value: entity.UserShort{}},
			},
		},
	})

	router.AddRoute(http.MethodGet, "/signin", s.signinHandler, swagger.Definitions{
		Querystring: swagger.ParameterValue{
			"query": {
				Schema: &swagger.Schema{Value: ""},
			},
		},
	})

	router.GenerateAndExposeOpenapi()
}
