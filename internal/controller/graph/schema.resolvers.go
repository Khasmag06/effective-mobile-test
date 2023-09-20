package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"fmt"

	"github.com/khasmag06/effective-mobile-test/internal/controller/graph/model"
)

// CreatePerson is the resolver for the createPerson field.
func (r *mutationResolver) CreatePerson(ctx context.Context, input model.PersonInput) (*model.Person, error) {
	panic(fmt.Errorf("not implemented: CreatePerson - createPerson"))
}

// UpdatePerson is the resolver for the updatePerson field.
func (r *mutationResolver) UpdatePerson(ctx context.Context, id int, input model.PersonInput) (*model.Person, error) {
	panic(fmt.Errorf("not implemented: UpdatePerson - updatePerson"))
}

// DeletePerson is the resolver for the deletePerson field.
func (r *mutationResolver) DeletePerson(ctx context.Context, id int) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeletePerson - deletePerson"))
}

// GetPeople is the resolver for the getPeople field.
func (r *queryResolver) GetPeople(ctx context.Context, page *int, limit *int, sortBy *string, sortOrder *string) ([]*model.Person, error) {
	panic(fmt.Errorf("not implemented: GetPeople - getPeople"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
