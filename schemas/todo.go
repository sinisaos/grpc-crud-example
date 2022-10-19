package schemas

type CreateTodoIn struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type UpdateTodoIn struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
