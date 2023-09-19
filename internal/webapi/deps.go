package webapi

import (
	"context"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
)

type peopleService interface {
	CreatePerson(ctx context.Context, person entity.Person) error
}

type logger interface {
	Error(text ...any)
	Errorf(format string, args ...any)
}
