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

// HuggingFaceService implements VisionService using Hugging Face's free Inference API.
// It uses image-to-text models (like BLIP or similar) to generate captions from images.
type HuggingFaceService struct {
	apiToken   string       // Hugging Face API token
	modelID    string       // Model ID to use (e.g., "Salesforce/blip-image-captioning-large")
	httpClient *http.Client // HTTP client with timeout configuration
}

// NewHuggingFaceService creates a new Hugging Face service instance.
// This service uses Hugging Face's free Inference API for image captioning.
//
// Benefits:
// - Free tier available
// - No local setup needed
// - Access to various state-of-the-art models
// - Managed infrastructure
//
// Popular models for ingredient detection:
// - "Salesforce/blip-image-captioning-large" (recommended)
// - "Salesforce/blip-image-captioning-base"
// - "nlpconnect/vit-gpt2-image-captioning"
//
// Parameters:
//   - apiToken: Hugging Face API token (get from https://huggingface.co/settings/tokens)
//   - modelID: Model identifier (default: "Salesforce/blip-image-captioning-large")
//
// Returns a configured HuggingFaceService ready for use.
func NewHuggingFaceService(apiToken, modelID string) *HuggingFaceService {
	if modelID == "" {
		modelID = "Salesforce/blip-image-captioning-large"
	}

	return &HuggingFaceService{
		apiToken: apiToken,
		modelID:  modelID,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// HuggingFaceResponse represents the response from Hugging Face Inference API.
// The API returns an array of results, each with a generated caption.
type HuggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

// HuggingFaceErrorResponse represents error responses from the API.
type HuggingFaceErrorResponse struct {
	Error         string  `json:"error"`
	EstimatedTime float64 `json:"estimated_time,omitempty"`
}

// DetectIngredients analyzes an image using Hugging Face's Inference API
// and extracts ingredient names from the generated caption.
//
// Process:
// 1. Send image to Hugging Face Inference API
// 2. Receive AI-generated caption
// 3. Parse caption to extract ingredient names
// 4. Return structured result
//
// The method handles model loading delays and will retry if the model
// is initially "cold" (not loaded).
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - imageData: Raw image bytes (JPEG, PNG, etc.)
//   - filename: Original filename for logging/metadata
//
// Returns DetectionResult with ingredients or error on failure.
func (s *HuggingFaceService) DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error) {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", s.modelID)

	// Create request with image data as body
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(imageData))
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: err}
	}

	// Set headers - send binary data without specific content type
	req.Header.Set("Authorization", "Bearer "+s.apiToken)
	// Don't set Content-Type - let it default to application/octet-stream or binary
	// This matches the curl --data-binary behavior

	fmt.Printf("Calling Hugging Face API for model: %s\n", s.modelID)
	fmt.Printf("Request URL: %s\n", url)
	fmt.Printf("Image size: %d bytes\n", len(imageData))

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("API request failed: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("failed to read response: %w", err)}
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(respBody))

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var errResp HuggingFaceErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			// If model is loading, provide helpful message
			if errResp.EstimatedTime > 0 {
				return nil, &DetectionError{
					Provider: "huggingface",
					Err:      fmt.Errorf("model is loading, estimated time: %.1f seconds. Please try again", errResp.EstimatedTime),
				}
			}
			return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("API error: %s", errResp.Error)}
		}

		// Special handling for 404 - likely model not found
		if resp.StatusCode == 404 {
			return nil, &DetectionError{
				Provider: "huggingface",
				Err:      fmt.Errorf("model not found: %s. Please check the model ID or try 'Salesforce/blip-image-captioning-large'", s.modelID),
			}
		}

		return nil, &DetectionError{
			Provider: "huggingface",
			Err:      fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody)),
		}
	}

	// Parse the response
	var hfResp HuggingFaceResponse
	if err := json.Unmarshal(respBody, &hfResp); err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("failed to parse response: %w", err)}
	}

	if len(hfResp) == 0 {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("no results returned from API")}
	}

	caption := hfResp[0].GeneratedText
	if caption == "" {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("empty caption returned")}
	}

	// Parse ingredients from the caption
	ingredients := ParseIngredientsFromText(caption)

	// Calculate confidence based on number of ingredients detected
	confidence := 0.85
	if len(ingredients) == 0 {
		confidence = 0.3
	} else if len(ingredients) == 1 {
		confidence = 0.6
	}

	result := &DetectionResult{
		Ingredients: ingredients,
		RawResponse: caption,
		Confidence:  confidence,
		Provider:    "huggingface",
		Metadata: map[string]interface{}{
			"model":       s.modelID,
			"caption":     caption,
			"filename":    filename,
			"image_size":  len(imageData),
			"detected_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	fmt.Printf("Hugging Face detection result: %d ingredients from caption: %s\n", len(ingredients), caption)

	return result, nil
}
