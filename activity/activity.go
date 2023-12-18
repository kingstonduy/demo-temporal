package activity

import (
	"context"
	"encoding/json"
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
	log.Info("Cif: ", account)

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
