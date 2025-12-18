package main

import (
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/htmltomarkdown/converter"
	"github.com/kawai-network/veridium/pkg/htmltomarkdown/plugin/base"
	"github.com/kawai-network/veridium/pkg/htmltomarkdown/plugin/commonmark"
)

func main() {
	input := `<strong>Bold Text</strong>`

	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("__"),
				// ...additional configurations for the plugin
			),
		),
	)

	markdown, err := conv.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: __Bold Text__
}
