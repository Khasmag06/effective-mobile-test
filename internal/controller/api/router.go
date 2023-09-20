package api

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/khasmag06/effective-mobile-test/internal/controller/graph"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/pkg/validator"
	"strconv"

	_ "github.com/khasmag06/effective-mobile-test/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ErrInvalidID = errors.New("invalid person id")

type peopleService interface {
	CreatePerson(ctx context.Context, person entity.Person) error
	UpdatePersonData(ctx context.Context, personID int, person entity.Person) error
	DeletePersonData(ctx context.Context, personID int) error
	GetPeople(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error)
}

type logger interface {
	Info(text ...any)
	Error(text ...any)
	Errorf(format string, args ...any)
}

type Handler struct {
	*gin.Engine
	*validator.CustomValidator
	peopleService peopleService
	logger        logger
}

func NewHandler(ps peopleService, l logger) *Handler {
	h := &Handler{
		Engine:          gin.New(),
		CustomValidator: validator.NewCustomValidator(),
		peopleService:   ps,
		logger:          l,
	}

	h.Use(gin.Recovery())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(ps, l)}))

	// GraphQL
	h.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/query")))

	h.POST("/query", gin.WrapH(srv))

	// Swagger
	h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := h.Group("/api")

	api.GET("people/get", h.getPeople)
	api.POST("person/create", h.addPerson)
	api.DELETE("person/delete/:id", h.deletePerson)
	api.PUT("person/update/:id", h.updatePerson)

	return h

}

func parseID(idQuery string) (int, error) {
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		return 0, ErrInvalidID
	}
	if id <= 0 {
		return 0, ErrInvalidID
	}

	return id, nil
}
