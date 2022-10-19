package models

type Todo struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	Completed bool   `json:"completed" gorm:"default:false"`
}
