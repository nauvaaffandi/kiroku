package receipts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kiroku/backend/internal/models"
)

func CallOCRSpace(ctx context.Context, apiKey string, imageBase64 string, mockOCR bool) (string, error) {
	if mockOCR || apiKey == "mock" {
		// Mock receipt text for quick E2E testing without calling external API
		return `SUPERMARKET TOKYO
Date: 2026-10-15 14:30:00
Item 1: Fresh Vegetables - 850 JPY
Item 2: Milk - 320 JPY
Item 3: Bread - 280 JPY
Subtotal: 1,450 JPY
Tax (10%): 145 JPY
TOTAL: 1,595 JPY
Payment: Cash`, nil
	}

	if apiKey == "" {
		apiKey = "helloworld" // Default free OCR.space test key if none provided
	}

	apiURL := "https://api.ocr.space/parse/image"

	data := url.Values{}
	data.Set("base64Image", imageBase64)
	data.Set("isOverlayRequired", "false")
	data.Set("language", "eng")
	data.Set("detectOrientation", "true")
	data.Set("scale", "true")

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create ocr request: %w", err)
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ocr http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ocr api status: %d", resp.StatusCode)
	}

	var ocrResp models.OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return "", fmt.Errorf("decode ocr response: %w", err)
	}

	if ocrResp.OCRExitCode != 1 || ocrResp.IsErroredOnProcessing {
		errMsg := ocrResp.ErrorMessage
		if errMsg == "" {
			errMsg = ocrResp.ErrorDetails
		}
		if errMsg == "" && len(ocrResp.ParsedResults) > 0 {
			errMsg = ocrResp.ParsedResults[0].ErrorMessage
		}
		return "", fmt.Errorf("ocr processing error: %s", errMsg)
	}

	var fullText strings.Builder
	for _, res := range ocrResp.ParsedResults {
		if res.ParsedText != "" {
			fullText.WriteString(res.ParsedText)
			fullText.WriteString("\n")
		}
	}

	text := strings.TrimSpace(fullText.String())
	if text == "" {
		return "", fmt.Errorf("ocr extracted empty text")
	}

	return text, nil
}
