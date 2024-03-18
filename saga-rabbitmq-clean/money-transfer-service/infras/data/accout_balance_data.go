package data

import (
	"context"

	"github.com/lengocson131002/go-clean-core/es"
	"github.com/ocb/mcs-money-transfer/domain/account"
)

const (
	IndexAccountBalance = "t24v2.fbnk_account.transf.1"
)

type accBalanceData struct {
	esClient es.ElasticSearchClient
}

type esAccountBalanceModel struct {
	OcbAccountNumber           string
	OcbBranchCode              string
	OcbCustomerNumber          string
	CustomerNumberJointProfile string
	Currency                   string
	AccountOpeningDate         string
	LastAccountStatusCode      string
	LastAccountStatusDate      string
	Category                   string
	AccountTitle               string
	ShortTitle                 string
	OpenActualBal              int64
	OnlineActualBal            int64
	WorkingBalance             int64
	AccountOfficer             string
	ConditionGroup             string
	CurrNo                     string
	Op_ts                      string
	Current_ts                 string
}

// GetBalance implements data.AccountData.
func (a *accBalanceData) GetBalance(ctx context.Context, accNumber string) (*account.AccountBalanceResponse, error) {
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"ocbAccountNumber": accNumber,
	// 		},
	// 	},
	// }

	// output, err := a.esClient.Search(
	// 	ctx,
	// 	fmt.Sprintf("%s*", IndexAccountBalance),
	// 	es.WithSearchQuery(query),
	// 	es.WithSearchSort([]string{"op_ts:desc"}),
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// if len(output) == 0 {
	// 	return nil, domain.ErrorAccountNotFound
	// }

	// var balRes esAccountBalanceModel
	// err = util.MapStruct(output[0], &balRes)
	// if err != nil {
	// 	return nil, err
	// }

	// return &account.AccountBalanceResponse{
	// 	Currency:        balRes.Currency,
	// 	OpenActualBal:   balRes.OpenActualBal,
	// 	WorkingBalance:  balRes.WorkingBalance,
	// 	OnlineActualBal: balRes.OnlineActualBal,
	// }, nil

	return &account.AccountBalanceResponse{
		Currency:        "VND",
		OpenActualBal:   1000000,
		WorkingBalance:  1000000,
		OnlineActualBal: 1000000,
	}, nil
}

func NewAccountBalanceData(esClient es.ElasticSearchClient) account.AccountBalanceData {
	return &accBalanceData{
		esClient: esClient,
	}
}
