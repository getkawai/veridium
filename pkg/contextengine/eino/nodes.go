package eino

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// NodeBuilder helps build Eino nodes
type NodeBuilder struct{}

// NewNodeBuilder creates a new node builder
func NewNodeBuilder() *NodeBuilder {
	return &NodeBuilder{}
}

// CreateLambdaNode creates a lambda node from a function
// The function signature should match compose.InvokeWOOpt[I, O]
func CreateLambdaNode[T, U any](fn func(context.Context, T) (U, error)) *compose.Lambda {
	return compose.InvokableLambda(fn)
}

// Helper function to create message schema from content
func createMessageSchema(role schema.RoleType, content string) *schema.Message {
	return &schema.Message{
		Role:    role,
		Content: content,
	}
}

