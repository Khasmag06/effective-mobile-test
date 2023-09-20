package graph

import (
	"context"
	"errors"
	"github.com/khasmag06/effective-mobile-test/internal/controller/graph/model"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/internal/repo/people/repoerrs"
)

const (
	defaultPaginationLimit = 10
	defaultPageNumber      = 1
)

// CreatePerson is the resolver for the createPerson field.
func (r *Resolver) CreatePerson(ctx context.Context, input model.PersonInput) (*model.Person, error) {
	newPerson := entity.Person{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  *input.Patronymic,
		Age:         input.Age,
		Gender:      input.Gender,
		Nationality: input.Nationality,
	}

	if err := r.Validate(newPerson); err != nil {
		r.logger.Errorf("validation err: %v", err)
		return nil, err
	}
	if err := r.peopleService.CreatePerson(ctx, newPerson); err != nil {
		r.logger.Errorf("failed to create person data: %v", err.Error())
		return nil, err
	}

	return &model.Person{
		Name:        newPerson.Name,
		Surname:     newPerson.Surname,
		Patronymic:  &newPerson.Patronymic,
		Age:         newPerson.Age,
		Gender:      newPerson.Gender,
		Nationality: newPerson.Nationality}, nil
}

// UpdatePerson is the resolver for the updatePerson field.
func (r *mutationResolver) UpdatePerson(ctx context.Context, id int, input model.PersonInput) (*model.Person, error) {
	newPerson := entity.Person{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  *input.Patronymic,
		Age:         input.Age,
		Gender:      input.Gender,
		Nationality: input.Nationality,
	}

	if err := r.Validate(newPerson); err != nil {
		r.logger.Errorf("validation err: %v", err)
		return nil, err
	}
	if err := r.peopleService.UpdatePersonData(ctx, id, newPerson); err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			r.logger.Errorf("person not found: %v", err)
			return nil, err
		}
		r.logger.Errorf("failed to update person: %v", err)
		return nil, err
	}
	return &model.Person{
		Name:        newPerson.Name,
		Surname:     newPerson.Surname,
		Patronymic:  &newPerson.Patronymic,
		Age:         newPerson.Age,
		Gender:      newPerson.Gender,
		Nationality: newPerson.Nationality}, nil
}

// DeletePerson is the resolver for the deletePerson field.
func (r *mutationResolver) DeletePerson(ctx context.Context, id int) (*bool, error) {
	if err := r.peopleService.DeletePersonData(ctx, id); err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			r.logger.Errorf("person not found: %v", err)
			return nil, err
		}
		r.logger.Errorf("failed to delete person: %v", err)
		return nil, err
	}

	result := true
	return &result, nil
}

// GetPeople is the resolver for the getPeople field.
func (r *queryResolver) GetPeople(ctx context.Context, page *int, limit *int, sortBy *string, sortOrder *string) ([]*model.Person, error) {

	if page == nil || *page <= 0 {
		defaultPage := defaultPageNumber
		page = &defaultPage
	}
	if limit == nil || *limit <= 0 {
		defaultLimit := defaultPaginationLimit
		limit = &defaultLimit
	}
	if sortBy == nil {
		defaultSortBy := "date"
		sortBy = &defaultSortBy
	}
	if sortOrder == nil {
		defaultSortOrder := "asc"
		sortOrder = &defaultSortOrder
	}

	people, err := r.peopleService.GetPeople(ctx, *page, *limit, *sortBy, *sortOrder)
	if err != nil {
		r.logger.Errorf("failed to fetch people data: %v", err)
		return nil, err
	}

	var result []*model.Person
	for _, person := range people {
		graphqlPerson := &model.Person{
			ID:          &person.ID,
			Name:        person.Name,
			Surname:     person.Surname,
			Patronymic:  &person.Patronymic,
			Age:         person.Age,
			Gender:      person.Gender,
			Nationality: person.Nationality,
		}
		result = append(result, graphqlPerson)
	}
	return result, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
