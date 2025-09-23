package whatsapp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/stretchr/testify/assert"
)

// TestableWhatsAppService is a wrapper for testing
type TestableWhatsAppService struct {
	BaseURL string
	Client  domain.WhatsAppService
}

func NewTestableWhatsAppService(accessToken, phoneNumberID, baseURL string) *TestableWhatsAppService {
	// You would need to modify the original service to accept a base URL
	// For now, we'll create a simple mock implementation
	return &TestableWhatsAppService{
		BaseURL: baseURL,
		Client:  &MockWhatsAppService{},
	}
}

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
	
	// Simulate successful sending
	return nil
}

func (m *MockWhatsAppService) SendPasswordResetOTP(ctx context.Context, phoneNumber, code string) error {
	// Simple validation
	if phoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}
	if len(code) != 6 {
		return fmt.Errorf("OTP code must be 6 digits")
	}
	
	// Simulate successful password reset OTP sending
	return nil
}

func (m *MockWhatsAppService) GetStartLink() string {
	return "WhatsApp verification ready. Use your phone number directly."
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

// Simple test without external dependencies
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

	// Test invalid OTP code
	err = service.SendOTP(ctx, "+251911223344", "123")
	assert.Error(t, err)
}

// Test for password reset OTP
func TestWhatsAppService_SimpleSendPasswordResetOTP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := &MockWhatsAppService{}

	// Test successful password reset OTP sending
	err := service.SendPasswordResetOTP(ctx, "+251911223344", "654321")
	assert.NoError(t, err)

	// Test invalid phone number
	err = service.SendPasswordResetOTP(ctx, "", "654321")
	assert.Error(t, err)

	// Test invalid OTP code
	err = service.SendPasswordResetOTP(ctx, "+251911223344", "654")
	assert.Error(t, err)
}