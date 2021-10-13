package main

import (
	"time"
)

type todoitem struct {
	Value     string
	CreatedAt time.Time
}

func NewTodoItem(value string) todoitem {
	return todoitem{Value: value, CreatedAt: time.Now()}
}
