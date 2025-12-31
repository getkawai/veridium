package catalog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/cli/cli/command"

	"github.com/kawai-network/veridium/cmd/docker-mcp/hints"
)

func Ls(ctx context.Context, dockerCli command.Cli, format Format) error {
	cfg, err := ReadConfigWithDefaultCatalog(ctx)
	if err != nil {
		return err
	}

	if format == JSON {
		data, err := json.Marshal(cfg)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		humanPrintCatalog(dockerCli, *cfg)
	}

	return nil
}

func humanPrintCatalog(dockerCli command.Cli, cfg Config) {
	if len(cfg.Catalogs) == 0 {
		fmt.Println("No catalogs configured.")
		return
	}

	for name, catalog := range cfg.Catalogs {
		fmt.Printf("%s: %s\n", name, catalog.DisplayName)
	}
	if hints.Enabled(dockerCli) {
		hints.TipCyan.Print("Tip: To browse a catalog's servers, use ")
		hints.TipCyanBoldItalic.Print("docker mcp catalog show <catalog-name>")
		fmt.Println()
	}
}
