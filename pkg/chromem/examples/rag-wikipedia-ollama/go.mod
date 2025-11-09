module github.com/kawai-network/veridium/pkg/chromem/examples/rag-wikipedia-ollama

go 1.21

require (
	github.com/kawai-network/veridium/pkg/chromem v0.0.0
	github.com/sashabaranov/go-openai v1.17.9
)

replace github.com/kawai-network/veridium/pkg/chromem => ./../..
