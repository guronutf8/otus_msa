package api

import (
	"context"
	"fmt"
	swagger "github.com/davidebianchi/gswagger"
	"github.com/davidebianchi/gswagger/support/gorilla"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"usercrud/internal/db"
	"usercrud/internal/entity"
)

type Server struct {
	db *db.DB
}

func NewServer() *Server {
	mongoEndpoint := os.Getenv("mongo_endpoint")
	mongoUser := os.Getenv("mongo_user")
	mongoPassword := os.Getenv("mongo_password")
	fmt.Println(mongoEndpoint)
	fmt.Println(mongoUser)
	fmt.Println(mongoPassword)
	return &Server{db: db.New(mongoEndpoint, mongoUser, mongoPassword)}
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

	router.AddRoute(http.MethodGet, "/mm", s.IndexHandlerMM, swagger.Definitions{})
	router.AddRoute(http.MethodGet, "/s1", s.IndexHandlerS1, swagger.Definitions{})
	router.AddRoute(http.MethodGet, "/s2", s.IndexHandlerS2, swagger.Definitions{})
	router.AddRoute(http.MethodGet, "/s3", s.IndexHandlerS3, swagger.Definitions{})
	router.AddRoute(http.MethodGet, "/s25", s.IndexHandlerS25, swagger.Definitions{})

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
	router.AddRoute(http.MethodDelete, "/user/{userId}", s.userDeleteHandler, swagger.Definitions{
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

	router.AddRoute(http.MethodPost, "/percent200", s.PercentPostHandler, swagger.Definitions{
		RequestBody: &swagger.ContentValue{
			Content: swagger.Content{
				"application/json": {Value: entity.Percent200{}},
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
	router.AddRoute(http.MethodGet, "/percent200", s.PercentGetHandler, swagger.Definitions{
		Querystring: swagger.ParameterValue{
			"query": {
				Schema: &swagger.Schema{Value: ""},
			},
		},
	})

	router.AddRoute(http.MethodGet, "/switch", s.SwitchApiHandler, swagger.Definitions{
		Querystring: swagger.ParameterValue{
			"query": {
				Schema: &swagger.Schema{Value: ""},
			},
		},
	})

	router.GenerateAndExposeOpenapi()

}
