// service/whatsapp_service.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
)

type MetaWhatsAppClient struct {
	accessToken   string
	phoneNumberID string
	httpClient    *http.Client
}

func NewMetaWhatsAppService(accessToken, phoneNumberID string) domain.WhatsAppService {
	if accessToken == "" {
		log.Fatal("WABA_ACCESS_TOKEN is required")
	}
	if phoneNumberID == "" {
		log.Fatal("WABA_PHONE_NUMBER_ID is required")
	}
	
	log.Printf("WhatsApp Service initialized with Phone Number ID: %s", phoneNumberID)
	log.Printf("Access token length: %d characters", len(accessToken))
	
	return &MetaWhatsAppClient{
		accessToken:   accessToken,
		phoneNumberID: phoneNumberID,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (m *MetaWhatsAppClient) SendOTP(ctx context.Context, phoneNumber, code string) error {
	log.Printf("Attempting to send WhatsApp OTP to: %s", phoneNumber)

	cleanNumber, err := m.formatPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	log.Printf("Formatted phone number: %s", cleanNumber)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                cleanNumber,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "digital_imci_otp", 
			"language": map[string]interface{}{
				"code": "en",
			},
			"components": []map[string]interface{}{
				{
					"type": "body",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": code,
						},
					},
				},
				{
					"type": "button",
					"sub_type": "url",
					"index": "0",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": code,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", m.phoneNumberID)
	log.Printf("Sending request to: %s", url)
	log.Printf("Request payload: %s", string(jsonData))
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.accessToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Meta API Response Status: %s", resp.Status)
	log.Printf("Meta API Response Body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return fmt.Errorf("WhatsApp API error: %s - %s", resp.Status, string(body))
		}
		
		if errorData, ok := errorResponse["error"].(map[string]interface{}); ok {
			errorCode := "unknown"
			if code, ok := errorData["code"].(float64); ok {
				errorCode = fmt.Sprintf("%.0f", code)
			}
			
			errorMessage := "unknown error"
			if msg, ok := errorData["message"].(string); ok {
				errorMessage = msg
			}
			
			return fmt.Errorf("WhatsApp API error [%s]: %s", errorCode, errorMessage)
		}
		
		return fmt.Errorf("WhatsApp API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("WhatsApp OTP %s sent successfully to %s", code, cleanNumber)
	return nil
}

func (m *MetaWhatsAppClient) GetStartLink() string {
	return "WhatsApp verification ready. Use your phone number directly."
}

func (m *MetaWhatsAppClient) formatPhoneNumber(phoneNumber string) (string, error) {
	re := regexp.MustCompile(`\D`)
	cleanNumber := re.ReplaceAllString(phoneNumber, "")

	// Handle Ethiopian numbers (+251)
	if strings.HasPrefix(cleanNumber, "251") {
		return cleanNumber, nil
	}
	
	// If number starts with 0, remove it and add 251
	if len(cleanNumber) > 0 && cleanNumber[0] == '0' {
		cleanNumber = cleanNumber[1:] // Remove leading 0
		if len(cleanNumber) == 9 { // Ethiopian numbers are 9 digits after 0
			cleanNumber = "251" + cleanNumber
		}
	}

	// If still doesn't start with country code, assume it's Ethiopian
	if !strings.HasPrefix(cleanNumber, "251") && len(cleanNumber) == 9 {
		cleanNumber = "251" + cleanNumber
	}

	if len(cleanNumber) < 10 {
		return "", fmt.Errorf("invalid phone number length: %s", cleanNumber)
	}

	return cleanNumber, nil
}

func (m *MetaWhatsAppClient) SendPasswordResetOTP(ctx context.Context, phoneNumber, code string) error {
	log.Printf("Attempting to send WhatsApp password reset OTP to: %s", phoneNumber)

	cleanNumber, err := m.formatPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	log.Printf("Formatted phone number: %s", cleanNumber)

	// Use a different template for password reset
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                cleanNumber,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "digital_imci_password_reset", // Different template for password reset
			"language": map[string]interface{}{
				"code": "en",
			},
			"components": []map[string]interface{}{
				{
					"type": "body",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": code,
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", m.phoneNumberID)
	log.Printf("Sending password reset request to: %s", url)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.accessToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Meta API Response Status: %s", resp.Status)
	log.Printf("Meta API Response Body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return fmt.Errorf("WhatsApp API error: %s - %s", resp.Status, string(body))
		}
		
		if errorData, ok := errorResponse["error"].(map[string]interface{}); ok {
			errorCode := "unknown"
			if code, ok := errorData["code"].(float64); ok {
				errorCode = fmt.Sprintf("%.0f", code)
			}
			
			errorMessage := "unknown error"
			if msg, ok := errorData["message"].(string); ok {
				errorMessage = msg
			}
			
			return fmt.Errorf("WhatsApp API error [%s]: %s", errorCode, errorMessage)
		}
		
		return fmt.Errorf("WhatsApp API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("WhatsApp password reset OTP %s sent successfully to %s", code, cleanNumber)
	return nil
}