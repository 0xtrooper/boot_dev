package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Validator interface {
	// Valid checks the object and returns any problems.
	// If len(problems) == 0 then the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}

func decode[T any](c *gin.Context, v T) (error) {
	if err := c.BindJSON(&v); err != nil{
		return fmt.Errorf("can not decode json: %w", err)
 	}
	return nil
}

func decodeValid[T Validator](ctx context.Context, c *gin.Context, v T) (map[string]string, error) {
	if err := c.BindJSON(&v); err != nil{
		return nil, fmt.Errorf("can not decode json: %w", err)
 	}

	if problems := v.Valid(ctx); len(problems) > 0 {
		return problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return nil, nil
}