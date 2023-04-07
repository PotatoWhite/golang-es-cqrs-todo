package commandApi

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/infra/command"
	"github.com/potato/simple-restful-api/pkg/domain/todospec"
	"net/http"
	"strconv"
)

func NewTodoRouter(group *gin.RouterGroup, evs command.EventStore) TodoRouter {

	router := &todoRouter{
		evs: evs,
	}

	group.POST(":userNo/todos", router.Create)
	group.PATCH(":userNo/todos/:id/title", router.UpdateTitle)
	group.PATCH(":userNo/todos/:id/status", router.UpdateStatus)
	group.DELETE(":userNo/todos/:id", router.Delete)
	return router
}

type TodoRouter interface {
	Create(c *gin.Context)
	UpdateTitle(c *gin.Context)
	UpdateStatus(c *gin.Context)
	Delete(c *gin.Context)
}

type todoRouter struct {
	evs command.EventStore
}

func (r *todoRouter) UpdateTitle(c *gin.Context) {
	userNo, aggregateID, err := r.parsePathParams(c)
	if err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	if !r.checkOwnership(aggregateID, userNo) {
		r.respondWithError(c, http.StatusForbidden, "you are not the owner")
		return
	}

	var reqBody todospec.UpdateTitle
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	resultEvent := todospec.NewTitleUpdatedEvent(aggregateID, userNo, reqBody.Title)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		r.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &publishEvent)
}

func (r *todoRouter) UpdateStatus(c *gin.Context) {
	userNo, aggregateID, err := r.parsePathParams(c)
	if err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	if !r.checkOwnership(aggregateID, userNo) {
		r.respondWithError(c, http.StatusForbidden, "you are not the owner")
		return
	}

	var reqBody todospec.UpdateStatus
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	resultEvent := todospec.NewStatusUpdatedEvent(aggregateID, userNo, reqBody.Status)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		r.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &publishEvent)
}

func (r *todoRouter) Delete(c *gin.Context) {
	userNo, aggregateID, err := r.parsePathParams(c)
	if err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	if !r.checkOwnership(aggregateID, userNo) {
		r.respondWithError(c, http.StatusForbidden, "you are not the owner")
		return
	}

	resultEvent := todospec.NewTodoDeletedEvent(aggregateID)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		r.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &publishEvent)
}

func (r *todoRouter) Create(c *gin.Context) {
	userNo, err := r.parseUserNo(c)
	if err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	var reqBody todospec.CreateTodo
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		r.respondWithError(c, http.StatusBadRequest, err)
		return
	}

	resultEvent := todospec.NewTodoCreatedEvent(uuid.New(), userNo, reqBody.Title)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		r.respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, publishEvent)
}

func (r *todoRouter) getPathParams(c *gin.Context) (userNo uint, aggregateId uuid.UUID, error error) {
	// get userNo(uint) from path
	userNoUint64, err := strconv.ParseUint(c.Param("userNo"), 10, 64)
	if err != nil {
		return userNo, aggregateId, err
	}
	userNo = uint(userNoUint64)

	// string to uuid
	if aggregateId, err = uuid.Parse(c.Param("id")); err != nil {
		return userNo, aggregateId, err
	}
	return userNo, aggregateId, error
}

func (r *todoRouter) respondWithError(c *gin.Context, status int, err interface{}) {
	c.JSON(status, gin.H{
		"error": err,
	})
}

func (r *todoRouter) parsePathParams(c *gin.Context) (userNo uint, aggregateID uuid.UUID, err error) {
	userNo, err = r.parseUserNo(c)
	if err != nil {
		return 0, uuid.Nil, err
	}
	aggregateID, err = uuid.Parse(c.Param("id"))
	if err != nil {
		return userNo, uuid.Nil, err
	}

	return userNo, aggregateID, nil
}

func (r *todoRouter) parseUserNo(c *gin.Context) (uint, error) {
	userNoUint64, err := strconv.ParseUint(c.Param("userNo"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userNoUint64), nil
}

func (r *todoRouter) checkOwnership(aggregateID uuid.UUID, userNo uint) bool {
	lastOne, err := r.evs.GetLastEvent(aggregateID)
	if err != nil || lastOne.UserNo != userNo {
		return false
	}
	return true
}
