// Package vision provides AI-powered image analysis for ingredient detection.
package vision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HuggingFaceService implements VisionService using the Hugging Face Inference API.
// It uses the BLIP-2 image captioning model which excels at food recognition.
type HuggingFaceService struct {
	apiKey     string       // Hugging Face API authentication key
	httpClient *http.Client // HTTP client with timeout configuration
	modelURL   string       // API endpoint URL for the vision model
}

// NewHuggingFaceService creates a new Hugging Face vision service instance.
// Uses the Salesforce BLIP image captioning model which is optimized for food detection.
//
// The BLIP model is excellent for:
// - Food and ingredient recognition
// - Detailed scene descriptions
// - High accuracy on kitchen/cooking images
//
// Parameters:
//   - apiKey: Hugging Face API authentication key
//
// Returns a configured HuggingFaceService ready for use.
func NewHuggingFaceService(apiKey string) *HuggingFaceService {
	return &HuggingFaceService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		modelURL: "https://api-inference.huggingface.co/models/Salesforce/blip-image-captioning-large",
	}
}

// HuggingFaceResponse represents the structure of the API response from Hugging Face.
// The model returns generated text describing the image content.
type HuggingFaceResponse struct {
	GeneratedText string `json:"generated_text"` // AI-generated caption describing the image
}

// DetectIngredients analyzes an image using Hugging Face's BLIP model
// and extracts ingredient names from the generated caption.
//
// Process:
// 1. Send image to Hugging Face API
// 2. Receive AI-generated caption
// 3. Parse caption to extract ingredient names
// 4. Calculate confidence score
// 5. Return structured result
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - imageData: Raw image bytes (JPEG, PNG, etc.)
//   - filename: Original filename for logging/metadata
//
// Returns DetectionResult with ingredients or error on failure.
func (s *HuggingFaceService) DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.modelURL, bytes.NewReader(imageData))
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("API request failed: %w", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("failed to read response: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &DetectionError{
			Provider: "huggingface",
			Err:      fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var responses []HuggingFaceResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		var errorResp map[string]interface{}
		if jsonErr := json.Unmarshal(body, &errorResp); jsonErr == nil {
			if errMsg, ok := errorResp["error"].(string); ok {
				return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("API error: %s", errMsg)}
			}
		}
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("failed to parse response: %w", err)}
	}

	if len(responses) == 0 {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("empty response from API")}
	}

	generatedText := responses[0].GeneratedText
	ingredients := ParseIngredientsFromText(generatedText)

	result := &DetectionResult{
		Ingredients: ingredients,
		RawResponse: generatedText,
		Confidence:  calculateConfidence(generatedText, ingredients),
		Provider:    "huggingface",
		Metadata: map[string]interface{}{
			"model":       "Salesforce/blip-image-captioning-large",
			"caption":     generatedText,
			"filename":    filename,
			"image_size":  len(imageData),
			"detected_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	return result, nil
}

// calculateConfidence estimates the detection confidence based on multiple factors.
// Higher confidence is assigned when:
// - More ingredients are detected
// - Caption contains food-related keywords
// - Caption is descriptive and detailed
//
// Returns a confidence score between 0.0 and 0.95 (never 100% certain).
func calculateConfidence(caption string, ingredients []string) float64 {
	if len(ingredients) == 0 {
		return 0.0
	}

	// Base confidence
	confidence := 0.7

	// Increase confidence with more ingredients detected
	if len(ingredients) >= 3 {
		confidence += 0.1
	}
	if len(ingredients) >= 5 {
		confidence += 0.1
	}

	foodWords := []string{"food", "dish", "plate", "bowl", "ingredients", "cooking", "meal"}
	for _, word := range foodWords {
		if contains(caption, word) {
			confidence += 0.05
			break
		}
	}

	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// contains performs a case-insensitive substring search.
func contains(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return bytes.Contains([]byte(s), []byte(substr))
}

// toLower converts a string to lowercase using byte operations.
func toLower(s string) string {
	return string(bytes.ToLower([]byte(s)))
}
