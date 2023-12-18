package activity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	utils "demo-temporal"
	model "demo-temporal/model"

	"go.temporal.io/sdk/activity"
)

func RegisterAccount(ctx context.Context) (model.Account, error) {
	log := activity.GetLogger(ctx)
	log.Info("RegisterAccount activity started")

	url := "http://localhost:8080/account/register"

	content, err := utils.ReceiveFromApi(url, "POST")
	if err != nil {
		return model.Account{}, err
	}

	var output model.Account

	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		return model.Account{}, err
	}

	return output, nil
}

func RegisterSms(ctx context.Context, account *model.Account) (model.Account, error) {
	log := activity.GetLogger(ctx)
	log.Info("RegisterAccount activity started")

	id := account.Cif
	url := fmt.Sprintf("http://localhost:8080/account/register/sms/%s/", id)

	content, err := utils.ReceiveFromApi(url, "POST")
	if err != nil {
		return model.Account{}, err
	}

	var output model.Account
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		return model.Account{}, err
	}

	return output, nil
}

func RegisterEmail(ctx context.Context, account *model.Account) (model.Account, error) {
	log := activity.GetLogger(ctx)
	log.Info("RegisterAccount activity started")

	id := account.Cif
	url := fmt.Sprintf("http://localhost:8080/account/register/email/%s/", id)

	content, err := utils.ReceiveFromApi(url, "POST")
	if err != nil {
		return model.Account{}, err
	}

	var output model.Account
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		return model.Account{}, err
	}

	return output, nil
}

func NotificationSms(ctx context.Context, account *model.Account) error {
	if account.IsSms == true {
		fmt.Println("🎇You have register SMS notification successfully!")
		return nil
	} else {
		return errors.New("Register SMS failed")
	}
}
