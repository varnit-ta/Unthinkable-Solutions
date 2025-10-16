// Package vision provides AI-powered image analysis for ingredient detection.
// It defines interfaces for vision services and common result structures.
package vision

import (
	"context"
	"fmt"
)

// VisionService defines the interface for AI-powered ingredient detection from images.
// Implementations should use computer vision or image captioning models
// to identify food ingredients in uploaded images.
type VisionService interface {
	// DetectIngredients analyzes an image and extracts ingredient names.
	// Returns a DetectionResult with ingredients and confidence metrics,
	// or an error if detection fails.
	DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error)
}

// DetectionResult contains the ingredients detected from an image
// along with confidence scores and metadata about the detection process.
type DetectionResult struct {
	Ingredients []string               `json:"ingredients"`           // List of detected ingredient names
	RawResponse string                 `json:"rawResponse,omitempty"` // Raw AI model output (e.g., caption)
	Confidence  float64                `json:"confidence"`            // Overall confidence score (0.0 to 1.0)
	Provider    string                 `json:"provider"`              // Name of the AI service provider
	Metadata    map[string]interface{} `json:"metadata,omitempty"`    // Additional metadata (model name, etc.)
}

// DetectionError is a custom error type for vision service failures.
// It includes the provider name to help with debugging and monitoring.
type DetectionError struct {
	Provider string // Name of the vision service provider (e.g., "huggingface")
	Err      error  // Underlying error that occurred
}

// Error returns a formatted error message including the provider name.
func (e *DetectionError) Error() string {
	return fmt.Sprintf("vision detection error (%s): %v", e.Provider, e.Err)
}

// Unwrap returns the underlying error for error chain inspection.
func (e *DetectionError) Unwrap() error {
	return e.Err
}
