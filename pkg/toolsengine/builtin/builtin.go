package builtin

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/Knetic/govaluate"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/toolsengine"
)

// RegisterAllBuiltinTools registers all builtin tools
func RegisterAllBuiltinTools(engine *toolsengine.ToolsEngine) error {
	log.Printf("Registering builtin tools...")

	// Register web search
	if err := RegisterWebSearchTool(engine); err != nil {
		log.Printf("Warning: Failed to register web search tool: %v", err)
	} else {
		log.Printf("✅ Registered: web-search")
	}

	// Register calculator
	if err := RegisterCalculatorTool(engine); err != nil {
		log.Printf("Warning: Failed to register calculator tool: %v", err)
	} else {
		log.Printf("✅ Registered: calculator")
	}

	log.Printf("Builtin tools registration complete")
	return nil
}

// RegisterWebSearchTool registers a web search tool
func RegisterWebSearchTool(engine *toolsengine.ToolsEngine) error {
	tool, err := toolsengine.NewToolBuilder("web-search", "web_search").
		WithDescription("Search the web for information").
		WithParameter("query", schema.String, "The search query", true).
		WithParameter("max_results", schema.Integer, "Maximum number of results (default: 10)", false).
		WithCategory("search").
		WithVersion("1.0.0").
		WithExecutor(func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			query, ok := args["query"].(string)
			if !ok {
				return nil, fmt.Errorf("query parameter is required")
			}

			maxResults := 10
			if mr, ok := args["max_results"].(float64); ok {
				maxResults = int(mr)
			}

			// TODO: Integrate with actual search API
			results := map[string]interface{}{
				"query": query,
				"results": []map[string]interface{}{
					{
						"title":   "Example Result 1",
						"url":     "https://example.com/1",
						"snippet": fmt.Sprintf("Search result for: %s", query),
					},
					{
						"title":   "Example Result 2",
						"url":     "https://example.com/2",
						"snippet": fmt.Sprintf("Another result for: %s", query),
					},
				},
				"count":       2,
				"max_results": maxResults,
			}

			return results, nil
		}).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build web search tool: %w", err)
	}

	return engine.RegisterTool(tool)
}

// RegisterCalculatorTool registers a calculator tool
func RegisterCalculatorTool(engine *toolsengine.ToolsEngine) error {
	tool, err := toolsengine.NewToolBuilder("calculator", "calculator").
		WithDescription("Perform mathematical calculations. Supports: +, -, *, /, sqrt(), sin(), cos(), tan(), pow(), pi, e").
		WithParameter("expression", schema.String, "Mathematical expression (e.g., '2 + 2', 'sqrt(16)', 'sin(pi/2)')", true).
		WithCategory("utility").
		WithVersion("1.0.0").
		WithExecutor(func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			expression, ok := args["expression"].(string)
			if !ok {
				return nil, fmt.Errorf("expression parameter is required")
			}

			// Create evaluator with math functions
			functions := map[string]govaluate.ExpressionFunction{
				"sqrt": func(args ...interface{}) (interface{}, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("sqrt requires 1 argument")
					}
					val, ok := args[0].(float64)
					if !ok {
						return nil, fmt.Errorf("sqrt argument must be a number")
					}
					return math.Sqrt(val), nil
				},
				"sin": func(args ...interface{}) (interface{}, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("sin requires 1 argument")
					}
					val, ok := args[0].(float64)
					if !ok {
						return nil, fmt.Errorf("sin argument must be a number")
					}
					return math.Sin(val), nil
				},
				"cos": func(args ...interface{}) (interface{}, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("cos requires 1 argument")
					}
					val, ok := args[0].(float64)
					if !ok {
						return nil, fmt.Errorf("cos argument must be a number")
					}
					return math.Cos(val), nil
				},
				"tan": func(args ...interface{}) (interface{}, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("tan requires 1 argument")
					}
					val, ok := args[0].(float64)
					if !ok {
						return nil, fmt.Errorf("tan argument must be a number")
					}
					return math.Tan(val), nil
				},
				"pow": func(args ...interface{}) (interface{}, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("pow requires 2 arguments")
					}
					base, ok1 := args[0].(float64)
					exp, ok2 := args[1].(float64)
					if !ok1 || !ok2 {
						return nil, fmt.Errorf("pow arguments must be numbers")
					}
					return math.Pow(base, exp), nil
				},
			}

			// Create expression
			expr, err := govaluate.NewEvaluableExpressionWithFunctions(expression, functions)
			if err != nil {
				return nil, fmt.Errorf("invalid expression: %w", err)
			}

			// Add constants
			parameters := map[string]interface{}{
				"pi": math.Pi,
				"e":  math.E,
			}

			// Evaluate
			result, err := expr.Evaluate(parameters)
			if err != nil {
				return nil, fmt.Errorf("evaluation error: %w", err)
			}

			return map[string]interface{}{
				"expression": expression,
				"result":     result,
			}, nil
		}).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build calculator tool: %w", err)
	}

	return engine.RegisterTool(tool)
}

