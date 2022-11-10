package schemas

type TodoSchemaIn struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
