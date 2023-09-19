package people

import (
	"context"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/internal/repo/people/repoerrs"
)

type service struct {
	repo repository
}

func New(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) CreatePerson(ctx context.Context, person entity.Person) error {
	return s.repo.CreatePerson(ctx, person)
}

func (s *service) UpdatePersonData(ctx context.Context, personID int, person entity.Person) error {
	exists, err := s.repo.CheckPersonExists(ctx, personID)
	if err != nil {
		return err
	}

	if !exists {
		return repoerrs.ErrNotFound
	}

	return s.repo.UpdatePersonData(ctx, personID, person)
}

func (s *service) DeletePersonData(ctx context.Context, personID int) error {
	exists, err := s.repo.CheckPersonExists(ctx, personID)
	if err != nil {
		return err
	}

	if !exists {
		return repoerrs.ErrNotFound
	}

	return s.repo.DeletePersonData(ctx, personID)
}

func (s *service) GetPeople(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error) {
	people, err := s.repo.GetPeople(ctx, page, limit, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}
	if people == nil {
		return []entity.Person{}, nil
	}
	return people, nil
}
