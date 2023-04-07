package queryApi

import (
	"github.com/gin-gonic/gin"
	"github.com/potato/simple-restful-api/infra/query"
	"strconv"
)

func NewTodoRouter(group *gin.RouterGroup, ets query.EntityStore) TodoRouter {
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
	ets query.EntityStore
}

func (t *todoRouter) GetAllByUserNo(c *gin.Context) {
	// get param from path
	userNo := c.Param("userNo")
	userNoUint, err := strconv.ParseUint(userNo, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	//TODO implement me
	list, err := t.ets.GetTodosByUserNo(uint(userNoUint))
	if err != nil {
		return
	}

	c.JSON(200, list)

}

func (t *todoRouter) GetByUserNoAndID(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
