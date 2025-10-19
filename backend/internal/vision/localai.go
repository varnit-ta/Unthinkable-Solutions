// Package vision provides AI-powered image analysis for ingredient detection.
package vision

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

// LocalAIService implements VisionService using a local Python AI service.
// It communicates with a FastAPI server running the Salesforce BLIP model.
type LocalAIService struct {
	serviceURL string
	httpClient *http.Client
}

// NewLocalAIService creates a new local AI service instance.
// This service calls a Python FastAPI server that runs the BLIP model locally.
//
// Benefits:
// - No API key needed
// - Faster inference (local)
// - No rate limits
// - Model is pre-loaded
//
// Parameters:
//   - serviceURL: URL of the Python AI service (e.g., http://localhost:8000)
//
// Returns a configured LocalAIService ready for use.
func NewLocalAIService(serviceURL string) *LocalAIService {
	return &LocalAIService{
		serviceURL: serviceURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// AIServiceResponse represents the structure of the response from the Python AI service.
type AIServiceResponse struct {
	Success     bool                   `json:"success"`
	Ingredients []string               `json:"ingredients"`
	Cuisine     string                 `json:"cuisine"`
	DishType    string                 `json:"dish_type"`
	Caption     string                 `json:"caption"`
	Confidence  float64                `json:"confidence"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Model       interface{}            `json:"model,omitempty"`
	Device      string                 `json:"device"`
}

// DetectIngredients analyzes an image using the local Python AI service
// and extracts ingredient names from the generated caption.
//
// Process:
// 1. Send image to local AI service
// 2. Receive AI-generated caption
// 3. Parse caption to extract ingredient names
// 4. Return structured result
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - imageData: Raw image bytes (JPEG, PNG, etc.)
//   - filename: Original filename for logging/metadata
//
// Returns DetectionResult with ingredients or error on failure.
func (s *LocalAIService) DetectIngredients(ctx context.Context, imageData []byte, filename string) (*DetectionResult, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	contentType := getContentTypeFromFilename(filename)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("failed to create form file: %w", err)}
	}

	if _, err := part.Write(imageData); err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("failed to write image data: %w", err)}
	}

	if err := writer.Close(); err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("failed to close writer: %w", err)}
	}

	url := s.serviceURL + "/detect"
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: err}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Printf("Calling local AI service at: %s (content-type: %s)\n", url, contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("AI service request failed: %w (is the service running?)", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("failed to read response: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("AI service error. Status: %d, Body: %s\n", resp.StatusCode, string(respBody))
		return nil, &DetectionError{
			Provider: "local-ai",
			Err:      fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(respBody)),
		}
	}

	var aiResp AIServiceResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("failed to parse response: %w", err)}
	}

	if !aiResp.Success {
		return nil, &DetectionError{Provider: "local-ai", Err: fmt.Errorf("AI service returned success=false")}
	}

	ingredients := aiResp.Ingredients

	if len(ingredients) == 0 && aiResp.Caption != "" {
		ingredients = ParseIngredientsFromText(aiResp.Caption)
		fmt.Printf("No ingredients in response, parsed %d from caption\n", len(ingredients))
	}

	fmt.Printf("Local AI service detected %d ingredients: %v\n", len(ingredients), ingredients)

	modelInfo := "local-ai"
	if aiResp.Model != nil {
		switch v := aiResp.Model.(type) {
		case map[string]interface{}:
			if clipModel, ok := v["clip"].(string); ok {
				modelInfo = clipModel
			} else if blipModel, ok := v["blip"].(string); ok {
				modelInfo = blipModel
			}
		case map[string]string:
			if clipModel, ok := v["clip"]; ok {
				modelInfo = clipModel
			} else if blipModel, ok := v["blip"]; ok {
				modelInfo = blipModel
			}
		case string:
			modelInfo = v
		}
	}

	result := &DetectionResult{
		Ingredients: ingredients,
		RawResponse: aiResp.Caption,
		Confidence:  aiResp.Confidence,
		Provider:    "local-ai",
		Metadata: map[string]interface{}{
			"model":       modelInfo,
			"device":      aiResp.Device,
			"caption":     aiResp.Caption,
			"cuisine":     aiResp.Cuisine,
			"dish_type":   aiResp.DishType,
			"details":     aiResp.Details,
			"filename":    filename,
			"image_size":  len(imageData),
			"detected_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	return result, nil
}

// getContentTypeFromFilename determines the MIME type based on file extension
func getContentTypeFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	default:
		return "image/jpeg"
	}
}
