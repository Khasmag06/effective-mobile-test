//go:generate mockgen -source=$GOFILE -destination=mocks_test.go -package=$GOPACKAGE
package people

import (
	"context"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
)

type repository interface {
	CreatePerson(ctx context.Context, person entity.Person) error
	UpdatePersonData(ctx context.Context, personID int, person entity.Person) error
	DeletePersonData(ctx context.Context, personID int) error
	GetPeople(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error)
	CheckPersonExists(ctx context.Context, personID int) (bool, error)
}
