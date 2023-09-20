package graph

import (
	"context"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/pkg/validator"
)

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

type Resolver struct {
	*validator.CustomValidator
	peopleService peopleService
	logger        logger
}

func NewResolver(ps peopleService, l logger) *Resolver {
	return &Resolver{
		CustomValidator: validator.NewCustomValidator(),
		peopleService:   ps,
		logger:          l,
	}
}
