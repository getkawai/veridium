package types

import "fmt"

// ProcessorError represents an error that occurred in a processor
type ProcessorError struct {
	ProcessorName string
	Message       string
	OriginalError error
}

func (e *ProcessorError) Error() string {
	if e.OriginalError != nil {
		return fmt.Sprintf("[%s] %s: %v", e.ProcessorName, e.Message, e.OriginalError)
	}
	return fmt.Sprintf("[%s] %s", e.ProcessorName, e.Message)
}

func (e *ProcessorError) Unwrap() error {
	return e.OriginalError
}

// PipelineError represents an error that occurred in the pipeline
type PipelineError struct {
	Message       string
	ProcessorName string
	OriginalError error
}

func (e *PipelineError) Error() string {
	if e.ProcessorName != "" {
		return fmt.Sprintf("Pipeline error in [%s]: %s", e.ProcessorName, e.Message)
	}
	return fmt.Sprintf("Pipeline error: %s", e.Message)
}

func (e *PipelineError) Unwrap() error {
	return e.OriginalError
}

