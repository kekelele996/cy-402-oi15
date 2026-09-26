package router

import (
	"cylawcase/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerCaseTodoRoutes 案件待办路由。
func (r *Router) registerCaseTodoRoutes(g *gin.RouterGroup) {
	cases := g.Group("/cases")
	cases.Use(middleware.AuthRequired(r.cfg))
	cases.GET("/:id/todos", r.caseTodo.ListByCase)
	cases.POST("/:id/todos", r.caseTodo.Create)
	cases.GET("/:id/assignees", r.caseTodo.Assignees)

	todos := g.Group("/case-todos")
	todos.Use(middleware.AuthRequired(r.cfg))
	todos.POST("/:id/complete", r.caseTodo.Complete)
}
