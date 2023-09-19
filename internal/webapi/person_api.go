package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/khasmag06/effective-mobile-test/config"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/pkg/validator"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type PersonInfoApi struct {
	client *http.Client
	apiCfg config.PersonApiConfig
	ps     peopleService
	logger logger
	*validator.CustomValidator
}

func New(cfg config.PersonApiConfig, ps peopleService, l logger) *PersonInfoApi {
	return &PersonInfoApi{
		client:          &http.Client{},
		apiCfg:          cfg,
		ps:              ps,
		logger:          l,
		CustomValidator: validator.NewCustomValidator(),
	}
}

func (p *PersonInfoApi) AddFioData(msg []byte) error {
	var person *entity.Person
	if err := json.Unmarshal(msg, &person); err != nil {
		return fmt.Errorf("error decoding age response: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		if err := p.getAge(person); err != nil {
			p.logger.Error(err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := p.getGender(person); err != nil {
			p.logger.Error(err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := p.getNationality(person); err != nil {
			p.logger.Error(err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if err := p.Validate(person); err != nil {
		p.logger.Errorf("validation error: %v", err)
		return err
	}
	if err := p.ps.CreatePerson(context.Background(), *person); err != nil {
		p.logger.Errorf("error adding person to database: %v", err)
		return err
	}
	return nil
}

func (p *PersonInfoApi) getAge(person *entity.Person) error {
	reqURL := fmt.Sprintf("%s?name=%s", p.apiCfg.AgeURL, person.Name)
	resp, err := p.client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("error getting age response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed age response: status not ok")
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error decoding age response: %w", err)
	}
	age := int(data["age"].(float64))
	person.Age = age

	return nil
}

func (p *PersonInfoApi) getGender(person *entity.Person) error {
	reqURL := fmt.Sprintf("%s?name=%s", p.apiCfg.GenderURL, person.Name)
	resp, err := p.client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("error getting gender response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed gender response: status not ok")
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error decoding gender response: %w", err)
	}
	gender := data["gender"].(string)
	person.Gender = gender

	return nil
}

type NationalityResponse struct {
	Countries []NationalityInfo `json:"country"`
}

type NationalityInfo struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func (p *PersonInfoApi) getNationality(person *entity.Person) error {
	reqURL := fmt.Sprintf("%s?name=%s", p.apiCfg.NationalityURL, person.Name)
	resp, err := p.client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("error getting nationality response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed nationality response: status not ok")
	}

	defer resp.Body.Close()

	var response NationalityResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error decoding nationality response: %w", err)
	}

	if len(response.Countries) == 0 {
		return fmt.Errorf("no nationality data available")
	}

	var maxProbability float64
	var mostProbableCountry string
	for _, info := range response.Countries {
		if info.Probability > maxProbability {
			maxProbability = info.Probability
			mostProbableCountry = info.CountryID
		}
	}
	person.Nationality = mostProbableCountry

	return nil
}
