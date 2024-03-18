package server

import (
	"context"

	"github.com/lengocson131002/go-clean-core/pipeline"
	"github.com/ocb/mcs-money-transfer/domain/account"
	pb "github.com/ocb/mcs-money-transfer/presentation/grpc/pb"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
}

func NewAccountServer() pb.AccountServiceServer {
	return &AccountServer{}
}

func (s *AccountServer) CheckBalance(ctx context.Context, req *pb.CheckAccountBalanceRequest) (*pb.CheckAccountBalanceResponse, error) {
	pipelineReq := &account.CheckBalanceRequest{
		Account: req.AccountNumber,
	}

	res, err := pipeline.Send[*account.CheckBalanceRequest, *account.CheckBalanceResponse](ctx, pipelineReq)
	if err != nil {
		return nil, err
	}

	return &pb.CheckAccountBalanceResponse{
		Balance:  res.Balance,
		Currency: res.Currency,
	}, nil
}
