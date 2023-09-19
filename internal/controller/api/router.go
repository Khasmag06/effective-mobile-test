package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/pkg/validator"
	"strconv"
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
	Infof(format string, args ...any)
	Warn(text ...any)
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

	api := h.Group("/api")

	api.GET("people/get", h.GetPeople)
	api.POST("person/create", h.AddPerson)
	api.DELETE("person/delete/:id", h.DeletePerson)
	api.PUT("person/update/:id", h.UpdatePerson)

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