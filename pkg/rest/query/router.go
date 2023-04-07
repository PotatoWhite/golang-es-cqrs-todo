package queryApi

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/potato/simple-restful-api/pkg/repository"
	"strconv"
)

func NewTodoRouter(group *gin.RouterGroup, ets repository.TodoStore) TodoRouter {
	router := &todoRouter{
		ets: ets,
	}

	group.GET(":userNo/todos", router.GetAllByUserNo)
	group.GET(":userNo/todos/:id", router.GetByUserNoAndID)

	return router
}

type TodoRouter interface {
	GetAllByUserNo(c *gin.Context)
	GetByUserNoAndID(c *gin.Context)
}

type todoRouter struct {
	ets repository.TodoStore
}

func (t *todoRouter) GetAllByUserNo(c *gin.Context) {
	// get param from path
	userNo, err := t.parsePathUserNo(c)
	if err != nil && userNo == 0 {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	//TODO implement me
	list, err := t.ets.GetTodosByUserNo(userNo)
	if err != nil {
		return
	}

	c.JSON(200, list)

}

func (t *todoRouter) GetByUserNoAndID(c *gin.Context) {
	userNo, id, err := t.parsePathParams(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	todo, err := t.ets.GetTodoByUserNoAndId(userNo, id)
	if err != nil {
		return
	}

	c.JSON(200, todo)
}

func (t *todoRouter) parsePathUserNo(c *gin.Context) (userNo uint, err error) {
	userNo, err = t.parseUserNo(c)
	if err != nil {
		return 0, err
	}

	return userNo, nil
}

func (t *todoRouter) parsePathParams(c *gin.Context) (userNo uint, aggregateID uuid.UUID, err error) {
	userNo, err = t.parseUserNo(c)
	if err != nil {
		return 0, uuid.Nil, err
	}
	aggregateID, err = uuid.Parse(c.Param("id"))
	if err != nil {
		return userNo, uuid.Nil, err
	}

	return userNo, aggregateID, nil
}

func (t *todoRouter) parseUserNo(c *gin.Context) (uint, error) {
	userNoUint64, err := strconv.ParseUint(c.Param("userNo"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userNoUint64), nil
}
