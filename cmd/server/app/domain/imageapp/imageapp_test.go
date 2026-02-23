package imageapp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/internal/paths"
	sd "github.com/kawai-network/veridium/pkg/stablediffusion"
)

type mockImageEngine struct {
	ready bool
	err   error
	calls []sd.ImgGenParams
}

func (m *mockImageEngine) IsReady() bool {
	return m.ready
}

func (m *mockImageEngine) GenerateImage(imgGenParams *sd.ImgGenParams, newImagePath string) error {
	if m.err != nil {
		return m.err
	}
	if imgGenParams != nil {
		m.calls = append(m.calls, *imgGenParams)
	}
	return os.WriteFile(newImagePath, []byte("png-bytes"), 0644)
}

func newTestLogger() *logger.Logger {
	return logger.New(io.Discard, logger.LevelInfo, "TEST", nil)
}

func TestGenerations_ContractSuccess(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	paths.SetDataDir(t.TempDir())

	engine := &mockImageEngine{ready: true}
	a := newApp(Config{
		Log:    newTestLogger(),
		Engine: engine,
	})

	reqBody := `{"prompt":"a cat astronaut","n":1,"size":"640x480","response_format":"b64_json","quality":"hd"}`
	r := httptest.NewRequest("POST", "/v1/images/generations", strings.NewReader(reqBody))

	enc := a.generations(context.Background(), r)
	resp, ok := enc.(*ImageGenerationResponse)
	if !ok {
		t.Fatalf("expected *ImageGenerationResponse, got %T", enc)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 image, got %d", len(resp.Data))
	}
	if resp.Data[0].B64JSON == "" {
		t.Fatal("expected base64 image payload")
	}
	if resp.Data[0].RevisedPrompt != "a cat astronaut" {
		t.Fatalf("unexpected revised prompt: %q", resp.Data[0].RevisedPrompt)
	}

	if len(engine.calls) != 1 {
		t.Fatalf("expected 1 engine call, got %d", len(engine.calls))
	}
	got := engine.calls[0]
	if got.Width != 640 || got.Height != 480 {
		t.Fatalf("unexpected image size: %dx%d", got.Width, got.Height)
	}
	if got.SampleSteps != qwenImageHDSteps {
		t.Fatalf("unexpected steps: %d", got.SampleSteps)
	}
	if got.CfgScale != qwenImageCFGScale {
		t.Fatalf("unexpected cfg scale: %f", got.CfgScale)
	}
}

func TestGenerations_UnavailableWhenEngineNotReady(t *testing.T) {
	a := newApp(Config{
		Log:    newTestLogger(),
		Engine: &mockImageEngine{ready: false},
	})

	r := httptest.NewRequest("POST", "/v1/images/generations", strings.NewReader(`{"prompt":"test"}`))
	enc := a.generations(context.Background(), r)
	errResp, ok := enc.(*errs.Error)
	if !ok {
		t.Fatalf("expected *errs.Error, got %T", enc)
	}
	if errResp.Code != errs.Unimplemented {
		t.Fatalf("expected %s, got %s", errs.Unimplemented, errResp.Code)
	}
}

func TestEdits_ContractSuccess(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	paths.SetDataDir(t.TempDir())

	editEngine := &mockImageEngine{ready: true}
	a := newApp(Config{
		Log:        newTestLogger(),
		EditEngine: editEngine,
	})

	imgData := base64.StdEncoding.EncodeToString([]byte("init-image"))
	maskData := base64.StdEncoding.EncodeToString([]byte("mask-image"))
	payload := map[string]any{
		"prompt":          "add sunglasses",
		"image":           "data:image/png;base64," + imgData,
		"mask":            "data:image/png;base64," + maskData,
		"n":               1,
		"size":            "512x512",
		"response_format": "b64_json",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	r := httptest.NewRequest("POST", "/v1/images/edits", strings.NewReader(string(body)))
	enc := a.edits(context.Background(), r)
	resp, ok := enc.(*ImageGenerationResponse)
	if !ok {
		t.Fatalf("expected *ImageGenerationResponse, got %T", enc)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 image, got %d", len(resp.Data))
	}

	if len(editEngine.calls) != 1 {
		t.Fatalf("expected 1 edit engine call, got %d", len(editEngine.calls))
	}
	got := editEngine.calls[0]
	if got.Prompt != "add sunglasses" {
		t.Fatalf("unexpected prompt: %q", got.Prompt)
	}
	if got.InitImagePath == "" {
		t.Fatal("expected init image path to be set")
	}
	if got.MaskImagePath == "" {
		t.Fatal("expected mask image path to be set")
	}
}

func TestEdits_InvalidRequest(t *testing.T) {
	a := newApp(Config{
		Log:        newTestLogger(),
		EditEngine: &mockImageEngine{ready: true},
	})

	r := httptest.NewRequest("POST", "/v1/images/edits", strings.NewReader(`{"prompt":"x"}`))
	enc := a.edits(context.Background(), r)
	errResp, ok := enc.(*errs.Error)
	if !ok {
		t.Fatalf("expected *errs.Error, got %T", enc)
	}
	if errResp.Code != errs.InvalidArgument {
		t.Fatalf("expected %s, got %s", errs.InvalidArgument, errResp.Code)
	}
}

func TestGenerations_EngineErrorReturnsInternal(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	paths.SetDataDir(t.TempDir())

	a := newApp(Config{
		Log:    newTestLogger(),
		Engine: &mockImageEngine{ready: true, err: errors.New("boom")},
	})

	r := httptest.NewRequest("POST", "/v1/images/generations", strings.NewReader(`{"prompt":"test"}`))
	enc := a.generations(context.Background(), r)
	errResp, ok := enc.(*errs.Error)
	if !ok {
		t.Fatalf("expected *errs.Error, got %T", enc)
	}
	if errResp.Code != errs.Internal {
		t.Fatalf("expected %s, got %s", errs.Internal, errResp.Code)
	}
}
