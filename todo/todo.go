package todo

import "time"

type Todo struct {
	Id        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

func New(title string) *Todo {
	return &Todo{
		Id:        0, //later call Store.NewID() func to create new id
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
}
