package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/khasmag06/effective-mobile-test/internal/entity"
	"github.com/khasmag06/effective-mobile-test/internal/repo/people/repoerrs"
	"net/http"
	"strconv"
)

const (
	defaultPaginationLimit = 10
	defaultPageNumber      = 1
)

// @Tags People
// @Summary addPerson
// @Description create a new person
// @ID createPerson
// @Accept  json
// @Produce json
// @Param input body entity.Person true "person info"
// @Success 201 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /person/create [post]
func (h *Handler) addPerson(c *gin.Context) {
	ctx := context.Background()
	var personReq entity.Person
	if err := c.Bind(&personReq); err != nil {
		h.logger.Errorf("json body binding error: %v", err)
		writeErrorResponse(c, http.StatusBadRequest, "invalid request body format")
		return
	}
	if err := h.Validate(personReq); err != nil {
		h.logger.Errorf("validation err: %v", err)
		writeErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.peopleService.CreatePerson(ctx, personReq); err != nil {
		h.logger.Errorf("failed to create person data: %v", err.Error())
		writeErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	writeSuccessResponse(c, http.StatusCreated, "success")

}

// @Tags People
// @Summary get list of people
// @Description get a list of people with pagination and sorting
// @ID getPeople
// @Accept json
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Number of items per page (default is 10)"
// @Param sortBy query string false "Sorting field (default is 'date')"
// @Param sortOrder query string false "Sorting order (default is 'asc')"
// @Success 200 {array} entity.Person "List of people"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /people/get [get]
func (h *Handler) getPeople(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = defaultPageNumber
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = defaultPaginationLimit
	}
	sortBy := c.DefaultQuery("sortBy", "date")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	ctx := context.Background()

	people, err := h.peopleService.GetPeople(ctx, page, limit, sortBy, sortOrder)
	if err != nil {
		h.logger.Errorf("failed to fetch people data: %v", err.Error())
		writeErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, people)
}

// @Tags People
// @Summary updatePerson
// @Description update a person
// @ID updatePerson
// @Accept  json
// @Produce json
// @Param id path int64 true "ID of the person to update"
// @Param input body entity.Person true "person info"
// @Success 200 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /person/update/{id} [put]
func (h *Handler) updatePerson(c *gin.Context) {
	personID, err := parseID(c.Param("id"))
	if err != nil {
		h.logger.Error(err.Error())
		writeErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := context.Background()
	var personReq entity.Person
	if err := c.Bind(&personReq); err != nil {
		h.logger.Errorf("json body binding error: %v", err)
		writeErrorResponse(c, http.StatusBadRequest, "invalid request body format")
		return
	}
	if err := h.Validate(personReq); err != nil {
		h.logger.Errorf("validation err: %v", err)
		writeErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.peopleService.UpdatePersonData(ctx, personID, personReq); err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			h.logger.Errorf("error when receiving update data: %v", err.Error())
			writeErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		h.logger.Errorf("failed to update person data: %v", err.Error())
		writeErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	writeSuccessResponse(c, http.StatusOK, "success")
}

// @Tags People
// @Summary deletePerson
// @Description delete a person
// @ID deletePerson
// @Accept  json
// @Produce json
// @Param id path int64 true "ID of the person to delete"
// @Success 200 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /person/delete/{id} [delete]
func (h *Handler) deletePerson(c *gin.Context) {
	personID, err := parseID(c.Param("id"))
	if err != nil {
		h.logger.Error(err.Error())
		writeErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := context.Background()
	if err := h.peopleService.DeletePersonData(ctx, personID); err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			h.logger.Errorf("error when receiving data to delete: %v", err.Error())
			writeErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		h.logger.Errorf("failed to delete person data: %v", err.Error())
		writeErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	writeSuccessResponse(c, http.StatusOK, "success")
}
