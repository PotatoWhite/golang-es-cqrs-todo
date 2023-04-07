package spec

type CreateTodo struct {
	Title string `binding:"required"`
}

type UpdateTitle struct {
	Title string `binding:"required"`
}

type UpdateStatus struct {
	Status string `binding:"required"`
}
