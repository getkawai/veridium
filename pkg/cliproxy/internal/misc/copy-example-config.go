package misc

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func CopyConfigTemplate(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := in.Close(); errClose != nil {
			slog.Warn("failed to close source config file", "error", errClose)
		}
	}()

	if err = os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := out.Close(); errClose != nil {
			slog.Warn("failed to close destination config file", "error", errClose)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
