package vision

import (
	"context"
	"fmt"
)

// VisionService defines the interface for ingredient detection from images
type VisionService interface {
	DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error)
}

// DetectionResult contains the detected ingredients and metadata
type DetectionResult struct {
	Ingredients []string               `json:"ingredients"`
	RawResponse string                 `json:"rawResponse,omitempty"`
	Confidence  float64                `json:"confidence"`
	Provider    string                 `json:"provider"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DetectionError represents an error during detection
type DetectionError struct {
	Provider string
	Err      error
}

func (e *DetectionError) Error() string {
	return fmt.Sprintf("vision detection error (%s): %v", e.Provider, e.Err)
}

func (e *DetectionError) Unwrap() error {
	return e.Err
}
