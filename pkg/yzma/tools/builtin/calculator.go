package builtin

import (
	"context"
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// RegisterCalculator registers the calculator tool
func RegisterCalculator(registry *tools.ToolRegistry) error {
	tool := tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "calculator",
		Description: "Perform mathematical calculations. Supports: +, -, *, /, sqrt(), sin(), cos(), tan(), pow(), pi, e",
		Parameters: map[string]any{
			"expression": map[string]any{
				"type":        "string",
				"description": "The mathematical expression to evaluate (e.g., '2 + 2', 'sqrt(16)', 'sin(pi/2)')",
			},
		},
		Required: []string{"expression"},
		Parallel: true, // Safe to run in parallel - pure computation, no side effects
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			expression, ok := args["expression"]
			if !ok || expression == "" {
				return "", fmt.Errorf("expression parameter is required")
			}

			// Create evaluator with math functions
			expr, err := govaluate.NewEvaluableExpression(expression)
			if err != nil {
				return "", fmt.Errorf("invalid expression: %w", err)
			}

			// Evaluate expression
			result, err := expr.Evaluate(nil)
			if err != nil {
				return "", fmt.Errorf("evaluation failed: %w", err)
			}

			// Format result
			return fmt.Sprintf("%v", result), nil
		},
	})

	return registry.Register(tool)
}
