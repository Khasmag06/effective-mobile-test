package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
)

const (
	maxPaginationLimit         = 10
	dateSortType        string = "created_at"
	nationalitySortType string = "nationality"
	ageSortType         string = "age"
	genderSortType      string = "gender"

	sortAscending  string = "ASC"
	sortDescending string = "DESC"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repo {
	return &repo{
		pool: db,
	}
}

func (r *repo) CreatePerson(ctx context.Context, person entity.Person) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO people (name, surname, patronymic, age, gender, nationality)
			VALUES ($1, $2, $3, $4, $5, $6)`, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("personRepo - CreatePerson - r.pool.Exec: %v", err)
	}

	return nil
}

func (r *repo) UpdatePersonData(ctx context.Context, personID int, person entity.Person) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE people 
			SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
			WHERE id = $7`, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, personID)
	if err != nil {
		return fmt.Errorf("personRepo - UpdatePerson - r.pool.Exec: %v", err)
	}

	return nil
}

func (r *repo) DeletePersonData(ctx context.Context, fioID int) error {
	_, err := r.pool.Exec(ctx,
		`DELETE 
			FROM people
			WHERE id = $1`, fioID)
	if err != nil {
		return fmt.Errorf("personRepo - DeletePerson - r.pool.Exec: %v", err)
	}

	return nil
}

func (r *repo) GetPeople(ctx context.Context, page int, limit int, sortBy, sortOrder string) ([]entity.Person, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}

	var sortField string
	switch sortBy {
	case "nationality":
		sortField = nationalitySortType
	case "gender":
		sortField = genderSortType
	case "age":
		sortField = ageSortType
	default:
		sortField = dateSortType
	}

	var sortDir string
	switch sortOrder {
	case "desc":
		sortDir = sortDescending
	default:
		sortDir = sortAscending
	}

	offset := (page - 1) * limit
	orderBy := fmt.Sprintf("%s %s", sortField, sortDir)

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, surname, patronymic, age, gender, nationality
             FROM people
             ORDER BY `+orderBy+`
             LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("personRepo - GetPeople - r.pool.Query: %w", err)
	}
	defer rows.Close()

	var people []entity.Person

	for rows.Next() {
		var person entity.Person

		err := rows.Scan(&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality)
		if err != nil {
			return nil, fmt.Errorf("personRepo - GetPeople - rows.Scan: %w", err)
		}

		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("personRepo - GetPeople - rows.Err: %w", err)
	}

	return people, nil
}

func (r *repo) CheckPersonExists(ctx context.Context, personID int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM people WHERE id = $1)`, personID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("personRepo - CheckPersonExists - r.pool.QueryRow: %v", err)
	}
	return exists, nil
}
