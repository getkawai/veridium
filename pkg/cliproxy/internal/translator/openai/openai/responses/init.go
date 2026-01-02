package responses

import (
	. "github.com/kawai-network/veridium/pkg/cliproxy/internal/constant"
	"github.com/kawai-network/veridium/pkg/cliproxy/internal/interfaces"
	"github.com/kawai-network/veridium/pkg/cliproxy/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenaiResponse,
		OpenAI,
		ConvertOpenAIResponsesRequestToOpenAIChatCompletions,
		interfaces.TranslateResponse{
			Stream:    ConvertOpenAIChatCompletionsResponseToOpenAIResponses,
			NonStream: ConvertOpenAIChatCompletionsResponseToOpenAIResponsesNonStream,
		},
	)
}
