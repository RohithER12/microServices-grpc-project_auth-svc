package utils

import (
	"fmt"

	"github.com/RohithER12/auth-svc/pkg/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	TWILIO_ACCOUNT_SID string
	TWILIO_AUTH_TOKEN  string
	TWILIO_SERVICE_ID  string
	client             *twilio.RestClient
)

func init() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		return
	}

	TWILIO_ACCOUNT_SID = config.TWILIO_ACCOUNT_SID
	TWILIO_AUTH_TOKEN = config.TWILIO_AUTH_TOCKEN
	TWILIO_SERVICE_ID = config.TWILIO_SERVICE_ID
	fmt.Println(
		"\nTWILIO_ACCOUNT_SID\n", TWILIO_ACCOUNT_SID,
		"\nTWILIO_AUTH_TOKEN\n", TWILIO_AUTH_TOKEN,
		"\nTWILIO_SERVICE_ID\n", TWILIO_SERVICE_ID,
	)
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})
}

func SendOtp(phone string) (string, error) {

	to := "+91" + phone
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(TWILIO_SERVICE_ID, params)

	fmt.Println(
		"mob\n", to,
		"\nparams\n\n", params, "\n\n...", TWILIO_SERVICE_ID)

	if err != nil {
		fmt.Println("Invalid PhoneNumber\n", err.Error())
		return "", err
	}

	key := *resp.Sid
	fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	return key, nil
}

func CheckOtp(phoneNumber string, code string) error {
	to := "+91" + phoneNumber
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(TWILIO_SERVICE_ID, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		fmt.Println("Correct!")
	} else {
		fmt.Println("Incorrect!")
	}
	return nil
}
