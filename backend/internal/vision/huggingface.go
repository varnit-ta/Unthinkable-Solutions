package vision

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HuggingFaceService implements VisionService using Hugging Face Inference API
type HuggingFaceService struct {
	apiKey     string
	httpClient *http.Client
	modelURL   string
}

// NewHuggingFaceService creates a new Hugging Face vision service
// Uses the BLIP-2 model for image captioning which is excellent for food detection
func NewHuggingFaceService(apiKey string) *HuggingFaceService {
	return &HuggingFaceService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Using Salesforce BLIP image captioning - great for food detection
		modelURL: "https://api-inference.huggingface.co/models/Salesforce/blip-image-captioning-large",
	}
}

// HuggingFaceResponse represents the API response structure
type HuggingFaceResponse struct {
	GeneratedText string `json:"generated_text"`
}

// DetectIngredients sends image to Hugging Face and extracts ingredients
func (s *HuggingFaceService) DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", s.modelURL, bytes.NewReader(imageData))
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: err}
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("API request failed: %w", err)}
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &DetectionError{Provider: "huggingface", Err: fmt.Errorf("failed to read response: %w", err)}
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, &DetectionError{
			Provider: "huggingface",
			Err:      fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body)),
		}
	}

	// Parse response - can be array or object depending on model state
	var responses []HuggingFaceResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		// Try parsing as single object (when model is loading)
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

	// Extract ingredients from generated text
	generatedText := responses[0].GeneratedText
	ingredients := ParseIngredientsFromText(generatedText)

	// Build result
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

// calculateConfidence estimates confidence based on caption quality and ingredient count
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

	// Check if caption contains food-related words
	foodWords := []string{"food", "dish", "plate", "bowl", "ingredients", "cooking", "meal"}
	for _, word := range foodWords {
		if contains(caption, word) {
			confidence += 0.05
			break
		}
	}

	// Cap at 0.95 (never 100% certain)
	if confidence > 0.95 {
		confidence = 0.95
	}

	return confidence
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return bytes.Contains([]byte(s), []byte(substr))
}

// toLower converts string to lowercase
func toLower(s string) string {
	return string(bytes.ToLower([]byte(s)))
}

// Helper method to encode image to base64 (for alternative API endpoints if needed)
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
