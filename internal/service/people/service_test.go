package people_test

import (
	"context"
	"errors"
	"github.com/khasmag06/effective-mobile-test/internal/repo/people/repoerrs"
	"github.com/khasmag06/effective-mobile-test/internal/service/people"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
)

func TestService_CreatePerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := people.NewMockrepository(ctrl)
	svc := people.New(mockRepo)

	tests := []struct {
		name        string
		inputPerson entity.Person
		mockResult  error
		expectedErr error
	}{
		{
			name: "valid person",
			inputPerson: entity.Person{
				ID:          1,
				Name:        "John",
				Surname:     "Doe",
				Patronymic:  "Smith",
				Age:         30,
				Gender:      "male",
				Nationality: "American",
			},
			mockResult:  nil,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().CreatePerson(gomock.Any(), test.inputPerson).Return(test.mockResult)

			err := svc.CreatePerson(context.Background(), test.inputPerson)

			assert.Equal(t, test.expectedErr, err, "Test case %s failed", test.name)
		})
	}
}

func TestService_UpdatePersonData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := people.NewMockrepository(ctrl)
	svc := people.New(mockRepo)

	tests := []struct {
		name           string
		personID       int
		inputPerson    entity.Person
		existsInRepo   bool
		checkExistsErr error
		updateErr      error
		expectedErr    error
	}{
		{
			name:     "valid update",
			personID: 1,
			inputPerson: entity.Person{
				ID:          1,
				Name:        "John",
				Surname:     "Doe",
				Patronymic:  "Smith",
				Age:         30,
				Gender:      "male",
				Nationality: "American",
			},
			existsInRepo:   true,
			checkExistsErr: nil,
			updateErr:      nil,
			expectedErr:    nil,
		},
		{
			name:     "person not found",
			personID: 2,
			inputPerson: entity.Person{
				ID:          2,
				Name:        "Alice",
				Surname:     "Smith",
				Patronymic:  "Johnson",
				Age:         25,
				Gender:      "female",
				Nationality: "Canadian",
			},
			existsInRepo:   false,
			checkExistsErr: nil,
			updateErr:      nil,
			expectedErr:    repoerrs.ErrNotFound,
		},
		{
			name:     "check exits person",
			personID: 3,
			inputPerson: entity.Person{
				ID:          3,
				Name:        "Bob",
				Surname:     "Brown",
				Patronymic:  "Williams",
				Age:         35,
				Gender:      "male",
				Nationality: "British",
			},
			existsInRepo:   false,
			checkExistsErr: errors.New("check exists error"),
			updateErr:      nil,
			expectedErr:    errors.New("check exists error"),
		},
		{
			name:     "Update Error",
			personID: 4,
			inputPerson: entity.Person{
				ID:          4,
				Name:        "Eva",
				Surname:     "Johnson",
				Patronymic:  "Davis",
				Age:         28,
				Gender:      "female",
				Nationality: "Australian",
			},
			existsInRepo:   true,
			checkExistsErr: nil,
			updateErr:      errors.New("update error"),
			expectedErr:    errors.New("update error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().CheckPersonExists(gomock.Any(), test.personID).Return(test.existsInRepo, test.checkExistsErr)
			if test.existsInRepo {
				mockRepo.EXPECT().UpdatePersonData(gomock.Any(), test.personID, test.inputPerson).Return(test.updateErr)
			}
			err := svc.UpdatePersonData(context.Background(), test.personID, test.inputPerson)

			assert.Equal(t, test.expectedErr, err, "Test case %s failed", test.name)
		})
	}
}

func TestService_DeletePersonData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := people.NewMockrepository(ctrl)
	svc := people.New(mockRepo)

	tests := []struct {
		name           string
		personID       int
		existsInRepo   bool
		checkExistsErr error
		deleteErr      error
		expectedErr    error
	}{
		{
			name:           "valid deletion",
			personID:       1,
			existsInRepo:   true,
			checkExistsErr: nil,
			deleteErr:      nil,
			expectedErr:    nil,
		},
		{
			name:           "person not found",
			personID:       2,
			existsInRepo:   false,
			checkExistsErr: nil,
			deleteErr:      nil,
			expectedErr:    repoerrs.ErrNotFound,
		},
		{
			name:           "check exists person",
			personID:       3,
			existsInRepo:   false,
			checkExistsErr: errors.New("check exists error"),
			deleteErr:      nil,
			expectedErr:    errors.New("check exists error"),
		},
		{
			name:           "delete error",
			personID:       4,
			existsInRepo:   true,
			checkExistsErr: nil,
			deleteErr:      errors.New("delete error"),
			expectedErr:    errors.New("delete error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().CheckPersonExists(gomock.Any(), test.personID).Return(test.existsInRepo, test.checkExistsErr)

			if test.existsInRepo {
				mockRepo.EXPECT().DeletePersonData(gomock.Any(), test.personID).Return(test.deleteErr)
			}
			err := svc.DeletePersonData(context.Background(), test.personID)

			assert.Equal(t, test.expectedErr, err, "Test case %s failed", test.name)
		})
	}
}

func TestService_GetPeople(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := people.NewMockrepository(ctrl)
	svc := people.New(mockRepo)

	tests := []struct {
		name           string
		page           int
		limit          int
		sortBy         string
		sortOrder      string
		repoResult     []entity.Person
		repoError      error
		expectedPeople []entity.Person
		expectedError  error
	}{
		{
			name:           "valid result",
			page:           1,
			limit:          10,
			sortBy:         "name",
			sortOrder:      "asc",
			repoResult:     []entity.Person{{ID: 1, Name: "John"}, {ID: 2, Name: "Alice"}},
			repoError:      nil,
			expectedPeople: []entity.Person{{ID: 1, Name: "John"}, {ID: 2, Name: "Alice"}},
			expectedError:  nil,
		},
		{
			name:           "empty result",
			page:           1,
			limit:          10,
			sortBy:         "name",
			sortOrder:      "asc",
			repoResult:     nil,
			repoError:      nil,
			expectedPeople: []entity.Person{},
			expectedError:  nil,
		},
		{
			name:           "repo error",
			page:           1,
			limit:          10,
			sortBy:         "name",
			sortOrder:      "asc",
			repoResult:     nil,
			repoError:      errors.New("repository error"),
			expectedPeople: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo.EXPECT().GetPeople(gomock.Any(), test.page, test.limit, test.sortBy, test.sortOrder).Return(test.repoResult, test.repoError)

			people, err := svc.GetPeople(context.Background(), test.page, test.limit, test.sortBy, test.sortOrder)

			assert.Equal(t, test.expectedPeople, people, "Test case %s failed: People not as expected", test.name)
			assert.Equal(t, test.expectedError, err, "Test case %s failed: Error not as expected", test.name)
		})
	}
}
