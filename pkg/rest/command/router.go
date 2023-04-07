package commandApi

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/pkg/domain/command"
	"github.com/potato/simple-restful-api/pkg/domain/spec"
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

	// validate path params
	userNo, aggregateID, err := r.getPathParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check authorization
	if !r.checkOwnership(aggregateID, userNo) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not the owner",
		})
		return
	}

	// get title from request body
	var reqBody spec.UpdateTitle
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add command to command store
	resultEvent := spec.NewTitleUpdatedEvent(aggregateID, userNo, reqBody.Title)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return response with command
	c.JSON(http.StatusOK, &publishEvent)
	return
}

func (r *todoRouter) checkOwnership(aggregateID uuid.UUID, userNo uint) bool {
	lastOne, err := r.evs.GetLastEvent(aggregateID)
	if err != nil || lastOne.UserNo != userNo {
		return false
	}
	return true
}

func (r *todoRouter) UpdateStatus(c *gin.Context) {

	// validate path params
	userNo, aggregateID, err := r.getPathParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check authorization
	if !r.checkOwnership(aggregateID, userNo) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not the owner",
		})
		return
	}

	// get title from request body
	var reqBody spec.UpdateStatus
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add command to command store
	resultEvent := spec.NewStatusUpdatedEvent(aggregateID, userNo, reqBody.Status)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return response with command
	c.JSON(http.StatusOK, &publishEvent)
	return
}

func (r *todoRouter) Delete(c *gin.Context) {
	// validate path params
	userNo, aggregateID, err := r.getPathParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check authorization
	if !r.checkOwnership(aggregateID, userNo) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not the owner",
		})
		return
	}

	// get title from request body
	var reqBody spec.UpdateTitle
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add command to command store
	resultEvent := spec.NewTodoDeletedEvent(aggregateID)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return response with command
	c.JSON(http.StatusOK, &publishEvent)
	return
}

func (r *todoRouter) Create(c *gin.Context) {
	// validate path params
	// get userNo(uint) from path
	userNoUint64, err := strconv.ParseUint(c.Param("userNo"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	userNo := uint(userNoUint64)

	// get title, userNo from request body
	var reqBody spec.CreateTodo
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add command to command store
	resultEvent := spec.NewTodoCreatedEvent(uuid.New(), userNo, reqBody.Title)
	publishEvent, err := r.evs.AddAndPublishEvent(userNo, &resultEvent)
	if err != nil {
		return
	}

	// return response
	c.JSON(http.StatusCreated, publishEvent)

	return
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
