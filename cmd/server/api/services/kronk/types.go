package kronk

// SetupResult contains the result of setup
type SetupResult struct {
	WalletCreated   bool
	WalletAddress   string
	LibraryReady    bool
	WhisperReady    bool
	StableDiffReady bool
	TTSReady        bool
	LLMReady        bool
	Errors          []error
	Warnings        []error
}
