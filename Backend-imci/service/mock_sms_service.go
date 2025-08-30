// service/mock_sms_service.go
package service

import (
	"context"
	"fmt"
	"log"
)

type MockSMSService struct {
}

func NewMockSMSService() *MockSMSService {
	return &MockSMSService{}
}

func (m *MockSMSService) SendOTP(ctx context.Context, toPhone, otpCode string) error {
	// In development, just log the OTP instead of sending real SMS
	log.Printf("MOCK SMS: OTP for %s is %s", toPhone, otpCode)
	fmt.Printf("📱 MOCK SMS to %s: Your OTP is %s\n", toPhone, otpCode)
	return nil
}