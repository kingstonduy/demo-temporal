package workflow

import (
	"context"
	"errors"
	"fmt"
	"kingstonduy/demo-temporal/async/model"
	"time"

	activity "demo-temporal/activity"

	"go.temporal.io/sdk/workflow"
)

func AsyncWorkFlow(ctx workflow.Context) error {

	// retryPolicy := &temporal.RetryPolicy{
	// 	InitialInterval: time.Second,
	// 	MaximumAttempts: 1, // unlimited retries
	// }

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 1,
		HeartbeatTimeout:    time.Minute * 1,
		// RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	future := workflow.ExecuteActivity(ctx, activity.GetOcbInfo)

	var flag bool
	err := workflow.ExecuteActivity(ctx, activity.Withdraw, nil).Get(ctx, &flag)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activity.UserInputOtp, flag).Get(ctx, &flag)

	for flag == false {
		err = workflow.ExecuteActivity(ctx, activity.ResendOtp, nil).Get(ctx, nil)
		if err != nil {
			return err
		}

		err = workflow.ExecuteActivity(ctx, activity.UserInputOtp, flag).Get(ctx, &flag)
		if err != nil {
			return err
		}
	}

	if flag == true {
		err = workflow.ExecuteActivity(ctx, activity.Notification, flag).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	err = future.Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func Withdraw(ctx context.Context) (bool, error) {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")

	url := fmt.Sprintf("http://localhost:8080/withdraw/")

	var requestType string
	var responseType string
	err := utils.PostApi(url, &requestType, &responseType)
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
	err := utils.GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡Resend OTP activity completed")
	return nil
}

func GetOcbInfo(ctx context.Context) error {
	log := activity.GetLogger(ctx)
	log.Info("get OCB info activity started")

	url := fmt.Sprintf("http://localhost:8080/OCB/info")

	var responseType string
	err := utils.GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡" + responseType)
	return nil
}

func UserInputOtp(ctx context.Context, inputt bool) (bool, error) {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")
	now := time.Now().Unix()
	// set 1 minute otp expiration
	for time.Now().Unix()-now <= 20 {
		var input string
		fmt.Println("Enter your OTP: ")
		fmt.Scanln(&input)

		var otpRequest = model.RequestOtp{
			Otp: input,
		}

		url := fmt.Sprintf("http://localhost:8080/otp/verify/")
		var responseType model.ResponseOtp
		err := utils.PostApi(url, &otpRequest, &responseType)
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
	return false, errors.New("The OTP is expired")
}

func Notification(ctx context.Context, inputt bool) error {
	log := activity.GetLogger(ctx)
	log.Info("Withdraw activity started")

	url := fmt.Sprintf("http://localhost:8080/notification/")

	var responseType string
	err := utils.GetApi(url, &responseType)
	if err != nil {
		log.Error("🔥Error when calling API", err)
		return err
	}
	fmt.Println("💡Notification: " + responseType)
	return nil
}
