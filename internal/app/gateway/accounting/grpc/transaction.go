package grpc

import (
	"context"
	"erp/gen/proto"
	"erp/internal/grpcutil"
	"erp/pkg/discovery"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionGateWay struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *TransactionGateWay {
	return &TransactionGateWay{registry}
}

func (g *TransactionGateWay) SaveTransaction(ctx context.Context,d *proto.TransactionLedger)(*emptypb.Empty, error){
	fmt.Println("SENDING TRANSACTION",d)
	conn, err := grpcutil.ServiceConnection(ctx, "accounting", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := proto.NewTransactionServiceClient(conn)
	var resp *emptypb.Empty
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		resp, err = client.SaveTransaction(ctx,d)
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, err
}

func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}

