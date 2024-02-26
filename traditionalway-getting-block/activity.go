package traditionalway_getting_block

import (
	"context"
	"fmt"
	"time"

	"traditionalway-getting-block/model"

	"go.temporal.io/sdk/activity"
)

func Withdraw(ctx context.Context) (bool, error) {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")

	url := fmt.Sprintf("http://localhost:8080/withdraw/")

	var requestType string
	var responseType string
	err := PostApi(url, &requestType, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return false, err
	}
	fmt.Println("💡Withdraw activity completed")
	return true, nil
}

func ResendOtp(ctx context.Context) error {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")

	url := fmt.Sprintf("http://localhost:8080/otp/resend/")

	var responseType string
	err := GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡Resend OTP activity completed")
	return nil
}

func UserInputOtp(ctx context.Context, inputt bool) (bool, error) {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")
	now := time.Now().Unix()
	// set 1 minute otp expiration
	for time.Now().Unix()-now <= 60 {
		var input string
		fmt.Println("Enter your OTP: ")
		fmt.Scanln(&input)

		var otpRequest = model.RequestOtp{
			Otp: input,
		}

		url := fmt.Sprintf("http://localhost:8080/otp/verify/")
		var responseType model.ResponseOtp
		err := PostApi(url, &otpRequest, &responseType)
		if err != nil && err.Error() == "" {
			log.Error("Error when calling API", err)
			return false, err
		}

		log.Info("Response", responseType)
		if responseType.Check == true {
			fmt.Println("💡The OTP is valid")
			return true, nil
		} else {
			fmt.Println("🔥The OTP is invalid")
		}
	}
	fmt.Println("🔥The OTP is expired")
	return false, nil
}
func Notification(ctx context.Context, inputt bool) error {
	log := activity.GetLogger(ctx)
	log.Info("Notification activity started")

	url := fmt.Sprintf("http://localhost:8080/notification/")

	var responseType string
	err := GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡Notification: " + responseType)
	return nil
}

func LongAcitivity(ctx context.Context) error {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")
	url := fmt.Sprintf("http://localhost:8080/OCB/info/")
	var responseType string
	err := GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡OCB info: " + responseType)
	return nil
}
