// service/sms_service.go
package service

import (
	"context"
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// Define interface
type SMSService interface {
	SendOTP(ctx context.Context, toPhone, otpCode string) error
}

// Real implementation
type TwilioSMSService struct {
	client      *twilio.RestClient
	fromNumber  string
}

func NewSMSService(accountSID, authToken, fromNumber string) *TwilioSMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	
	return &TwilioSMSService{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (s *TwilioSMSService) SendOTP(ctx context.Context, toPhone, otpCode string) error {
	message := fmt.Sprintf("Your Digital IMCI verification code is: %s. This code expires in 5 minutes.", otpCode)
	
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(s.fromNumber)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	
	log.Printf("OTP sent via SMS to %s", toPhone)
	return nil
}