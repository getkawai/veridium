package chat_completions

import (
	. "github.com/kawai-network/veridium/pkg/cliproxy/internal/constant"
	"github.com/kawai-network/veridium/pkg/cliproxy/internal/interfaces"
	"github.com/kawai-network/veridium/pkg/cliproxy/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenAI,
		OpenAI,
		ConvertOpenAIRequestToOpenAI,
		interfaces.TranslateResponse{
			Stream:    ConvertOpenAIResponseToOpenAI,
			NonStream: ConvertOpenAIResponseToOpenAINonStream,
		},
	)
}
