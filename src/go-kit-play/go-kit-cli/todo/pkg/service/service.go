package service

import (
	"context"
	"go-kit-play/go-kit-cli/todo/pkg/db"
	"go-kit-play/go-kit-cli/todo/pkg/io"

	"gopkg.in/mgo.v2/bson"
)

// TodoService describes the service.
type TodoService interface {
	Get(ctx context.Context) (t []io.Todo, error error)
	Add(ctx context.Context, todo io.Todo) (t io.Todo, error error)
	SetComplete(ctx context.Context, id string) (error error)
	RemoveComplete(ctx context.Context, id string) (error error)
	Delete(ctx context.Context, id string) (error error)
	// Example we want to add a get by id method.
	GetById(ctx context.Context, id string) (t io.Todo, error error)
}
type basicTodoService struct{}

func (b *basicTodoService) Get(ctx context.Context) (t []io.Todo, error error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return t, err
	}
	defer session.Close()
	c := session.DB("todo_app").C("todos")
	error = c.Find(nil).All(&t)
	return t, error
}

func (b *basicTodoService) Add(ctx context.Context, todo io.Todo) (t io.Todo, error error) {
	todo.Id = bson.NewObjectId()
	session, err := db.GetMongoSession()
	if err != nil {
		return t, err
	}
	defer session.Close()
	c := session.DB("todo_app").C("todos")
	error = c.Insert(&todo)
	return todo, error
}
func (b *basicTodoService) SetComplete(ctx context.Context, id string) (error error) {
	// TODO implement the business logic of SetComplete
	return error
}
func (b *basicTodoService) RemoveComplete(ctx context.Context, id string) (error error) {
	// TODO implement the business logic of RemoveComplete
	return error
}
func (b *basicTodoService) Delete(ctx context.Context, id string) (error error) {
	// TODO implement the business logic of Delete
	return error
}

// NewBasicTodoService returns a naive, stateless implementation of TodoService.
func NewBasicTodoService() TodoService {
	return &basicTodoService{}
}

// New returns a TodoService with all of the expected middleware wired in.
func New(middleware []Middleware) TodoService {
	var svc TodoService = NewBasicTodoService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func (b *basicTodoService) GetById(ctx context.Context, id string) (t io.Todo, error error) {
	// TODO implement the business logic of GetById
	return t, error
}
