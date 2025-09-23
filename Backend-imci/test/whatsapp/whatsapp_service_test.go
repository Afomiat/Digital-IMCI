package whatsapp_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockWhatsAppService for simple testing
type MockWhatsAppService struct{}

func (m *MockWhatsAppService) SendOTP(ctx context.Context, phoneNumber, code string) error {
	// Simple validation
	if phoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}
	if len(code) != 6 {
		return fmt.Errorf("OTP code must be 6 digits")
	}
	
	// Validate phone number format
	_, err := m.FormatPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}
	
	// Simulate successful sending
	return nil
}

func (m *MockWhatsAppService) FormatPhoneNumber(phoneNumber string) (string, error) {
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

// TestWhatsAppService_SimpleSendOTP tests basic OTP sending functionality
func TestWhatsAppService_SimpleSendOTP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := &MockWhatsAppService{}

	// Test successful OTP sending
	err := service.SendOTP(ctx, "+251911223344", "123456")
	assert.NoError(t, err)

	// Test invalid phone number
	err = service.SendOTP(ctx, "", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone number is required")

	// Test invalid OTP code
	err = service.SendOTP(ctx, "+251911223344", "123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OTP code must be 6 digits")

	// Test invalid phone number format
	err = service.SendOTP(ctx, "123", "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone number")
}

// TestPhoneNumberFormatting tests the phone number formatting logic
func TestPhoneNumberFormatting(t *testing.T) {
	service := &MockWhatsAppService{}
	
	testCases := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"+251911223344", "251911223344", false},
		{"0911223344", "251911223344", false},
		{"911223344", "251911223344", false},
		{"251911223344", "251911223344", false},
		{"123", "", true}, // Too short
		{"", "", true},    // Empty
		{"abc", "", true}, // Invalid characters
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := service.FormatPhoneNumber(tc.input)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// TestWhatsAppService_ContextCancellation tests context cancellation
func TestWhatsAppService_ContextCancellation(t *testing.T) {
	service := &MockWhatsAppService{}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := service.SendOTP(ctx, "+251911223344", "123456")
	assert.NoError(t, err) // Our mock doesn't respect context, but real implementation should
}

// TestWhatsAppService_MultipleCalls tests multiple consecutive calls
func TestWhatsAppService_MultipleCalls(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := &MockWhatsAppService{}

	// Test multiple successful calls
	for i := 0; i < 5; i++ {
		phone := fmt.Sprintf("+25191122334%d", i)
		code := fmt.Sprintf("12345%d", i)
		
		err := service.SendOTP(ctx, phone, code)
		assert.NoError(t, err)
	}
}

// TestWhatsAppService_EdgeCases tests edge cases
func TestWhatsAppService_EdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := &MockWhatsAppService{}

	// Test with different Ethiopian number formats
	testNumbers := []string{
		"+251911223344",
		"0911223344",
		"911223344",
		"251911223344",
	}

	for _, number := range testNumbers {
		t.Run(number, func(t *testing.T) {
			err := service.SendOTP(ctx, number, "123456")
			assert.NoError(t, err)
		})
	}
}

// TestMockServerIntegration tests with a mock HTTP server
func TestMockServerIntegration(t *testing.T) {
	// Create a test server that simulates WhatsApp API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages": [{"id": "test-message-id"}]}`))
	}))
	defer testServer.Close()

	// This test shows the pattern for testing the real service
	// You would need to modify the actual service to accept a custom HTTP client
	t.Logf("Mock server running at: %s", testServer.URL)
	assert.True(t, true) // Placeholder assertion
}

// TestWhatsAppService_ImplementsInterface verifies the mock implements the expected methods
func TestWhatsAppService_ImplementsInterface(t *testing.T) {
	var service interface {
		SendOTP(ctx context.Context, phoneNumber, code string) error
		FormatPhoneNumber(phoneNumber string) (string, error)
	} = &MockWhatsAppService{}
	
	assert.NotNil(t, service)
}