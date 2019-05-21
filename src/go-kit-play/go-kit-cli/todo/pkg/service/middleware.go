package service

import (
	"context"
	io "go-kit-play/go-kit-cli/todo/pkg/io"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(TodoService) TodoService

type loggingMiddleware struct {
	logger log.Logger
	next   TodoService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a TodoService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next TodoService) TodoService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Get(ctx context.Context) (t []io.Todo, error error) {
	defer func() {
		l.logger.Log("method", "Get", "t", t, "error", error)
	}()
	return l.next.Get(ctx)
}
func (l loggingMiddleware) Add(ctx context.Context, todo io.Todo) (t io.Todo, error error) {
	defer func() {
		l.logger.Log("method", "Add", "todo", todo, "t", t, "error", error)
	}()
	return l.next.Add(ctx, todo)
}
func (l loggingMiddleware) SetComplete(ctx context.Context, id string) (error error) {
	defer func() {
		l.logger.Log("method", "SetComplete", "id", id, "error", error)
	}()
	return l.next.SetComplete(ctx, id)
}
func (l loggingMiddleware) RemoveComplete(ctx context.Context, id string) (error error) {
	defer func() {
		l.logger.Log("method", "RemoveComplete", "id", id, "error", error)
	}()
	return l.next.RemoveComplete(ctx, id)
}
func (l loggingMiddleware) Delete(ctx context.Context, id string) (error error) {
	defer func() {
		l.logger.Log("method", "Delete", "id", id, "error", error)
	}()
	return l.next.Delete(ctx, id)
}

func (l loggingMiddleware) GetById(ctx context.Context, id string) (t io.Todo, error error) {
	defer func() {
		l.logger.Log("method", "GetById", "id", id, "t", t, "error", error)
	}()
	return l.next.GetById(ctx, id)
}

type authMiddleware struct {
	next TodoService
}

// AuthMiddleware returns a TodoService Middleware.
func AuthMiddleware() Middleware {
	return func(next TodoService) TodoService {
		return &authMiddleware{next}
	}

}
func (a authMiddleware) Get(ctx context.Context) (t []io.Todo, error error) {
	// Implement your middleware logic here

	return a.next.Get(ctx)
}
func (a authMiddleware) Add(ctx context.Context, todo io.Todo) (t io.Todo, error error) {
	// Implement your middleware logic here

	return a.next.Add(ctx, todo)
}
func (a authMiddleware) SetComplete(ctx context.Context, id string) (error error) {
	// Implement your middleware logic here

	return a.next.SetComplete(ctx, id)
}
func (a authMiddleware) RemoveComplete(ctx context.Context, id string) (error error) {
	// Implement your middleware logic here

	return a.next.RemoveComplete(ctx, id)
}
func (a authMiddleware) Delete(ctx context.Context, id string) (error error) {
	// Implement your middleware logic here

	return a.next.Delete(ctx, id)
}
func (a authMiddleware) GetById(ctx context.Context, id string) (t io.Todo, error error) {
	// Implement your middleware logic here

	return a.next.GetById(ctx, id)
}
