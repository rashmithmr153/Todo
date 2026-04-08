package store

import (
	"encoding/json"
	"fmt"
	"os"
	"todo/todo"
)

type Store struct {
	FilePath string
	Todos    []todo.Todo
	LastID   int
}

func NewStore(filepath string) *Store {
	return &Store{
		FilePath: filepath,
		LastID:   0,
	}
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("Unable to open stored todos from file(%s) : %s", s.FilePath, err)
	}

	var existTodo []todo.Todo
	if err := json.Unmarshal(data, &existTodo); err != nil {
		return fmt.Errorf("Unable to parse json from (%s) : %s", s.FilePath, err)
	}
	s.Todos = existTodo
	maxID := 0
	for _, item := range existTodo {
		if item.Id > maxID {
			maxID = item.Id
		}
	}
	s.LastID = maxID
	return nil
}

func (s *Store) Save() error {
	data, err := json.Marshal(s.Todos)
	if err != nil {
		return fmt.Errorf("Unable to parse todos to json:%s", err)
	}
	if err := os.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("Unable to write to %s file : %s", s.FilePath, err)
	}
	return nil
}

func (s *Store) Add(title string) error {
	newTodo := todo.New(title)
	s.LastID++
	newTodo.Id = s.LastID
	s.Todos = append(s.Todos, *newTodo)
	return s.Save()
}

func (s *Store) Delete(Id int) error {
	for k, value := range s.Todos {
		if value.Id == Id {
			s.Todos = append(s.Todos[:k], s.Todos[k+1:]...)
		}
	}
	return s.Save()
}

func (s *Store) MarkDone(Id int) error {
	for k, v := range s.Todos {
		if v.Id == Id {
			s.Todos[k].Done = true
		}
	}
	return s.Save()
}
