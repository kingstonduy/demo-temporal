package versionning_build_id

import (
	"context"
	"fmt"

	"kingstonduy/demo-temporal/versionning-build-id/services"
)

func GetInformation(ctx context.Context) (string, error) {
	var res string
	res = services.GetInformation()

	fmt.Println(res)
	return res, nil
}

func GetInformation1(ctx context.Context) (string, error) {
	var res string
	res = services.GetInformation1()

	fmt.Println(res)
	return res, nil
}

func GetInformation2(ctx context.Context) (string, error) {
	var res string
	res = services.GetInformation2()

	fmt.Println(res)
	return res, nil
}
