package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type repo struct {
	repository
	redis  *redis.Client
	logger logger
}

func New(rdb *redis.Client, peopleRepo repository, logger logger) *repo {
	return &repo{
		repository: peopleRepo,
		redis:      rdb,
		logger:     logger,
	}
}

const expiration = 48 * time.Hour // two days

func (r *repo) GetPeople(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error) {
	peopleDataCache, err := r.GetPeopleFromCache(ctx, page, limit, sortBy, sortOrder)
	if err != nil && !errors.Is(err, redis.Nil) {
		r.logger.Error(err)
	}
	if peopleDataCache != nil {
		return peopleDataCache, nil
	}

	peopleData, err := r.repository.GetPeople(ctx, page, limit, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	if err := r.SavePeopleToCache(ctx, page, limit, sortBy, sortOrder, peopleData); err != nil {
		r.logger.Error(err)
	}
	return peopleData, nil
}

func (r *repo) CreatePerson(ctx context.Context, person entity.Person) error {
	if err := r.repository.CreatePerson(ctx, person); err != nil {
		return err
	}
	if err := r.DeletePeopleFromCache(ctx); err != nil {
		r.logger.Error(err)
	}
	return nil
}

func (r *repo) UpdatePersonData(ctx context.Context, personID int, person entity.Person) error {
	if err := r.repository.UpdatePersonData(ctx, personID, person); err != nil {
		return err
	}
	if err := r.DeletePeopleFromCache(ctx); err != nil {
		r.logger.Error(err)
	}

	return nil
}

func (r *repo) DeletePersonData(ctx context.Context, personID int) error {
	if err := r.repository.DeletePersonData(ctx, personID); err != nil {
		return err
	}
	if err := r.DeletePeopleFromCache(ctx); err != nil {
		r.logger.Error(err)
	}
	return nil
}

func (r *repo) SavePeopleToCache(ctx context.Context, page int, limit int, sortBy, sortOrder string, peopleData []entity.Person) error {
	key := fmt.Sprintf("p:%d:%d:%s:%s", page, limit, sortBy, string(sortOrder[0])) // p - people
	peopleJSON, err := json.Marshal(peopleData)
	if err != nil {
		return err
	}
	if err := r.redis.Set(ctx, key, peopleJSON, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetPeopleFromCache(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error) {
	key := fmt.Sprintf("p:%d:%d:%s:%s", page, limit, sortBy, string(sortOrder[0])) // p - people
	peopleJSON, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var peopleDataCache []entity.Person
	if err := json.Unmarshal([]byte(peopleJSON), &peopleDataCache); err != nil {
		return nil, err
	}
	return peopleDataCache, nil
}

func (r *repo) DeletePeopleFromCache(ctx context.Context) error {
	pattern := "p:*" // p - people
	keysToDelete, err := r.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keysToDelete) > 0 {
		if err := r.redis.Del(ctx, keysToDelete...).Err(); err != nil {
			return err
		}
	}

	return nil
}
