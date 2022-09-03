package api

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/config"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/kafka"
	kafkaConfig "gitlab.ozon.dev/N0fail/price-tracker-validator/internal/kafka/config"
	pb "gitlab.ozon.dev/N0fail/price-tracker-validator/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func New(clientMain pb.AdminClient) pb.AdminServer {
	return &implementation{
		clientMain:     clientMain,
		kafkaRequester: kafka.New(),
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	clientMain     pb.AdminClient
	kafkaRequester kafka.RequestProducerI
}

func (i *implementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	if in.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrEmptyCode.Error())
	}

	if len(in.GetName()) < config.MinNameLength {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrNameTooShortError.Error())
	}

	if strings.ContainsAny(in.GetCode(), config.InvalidCodeSymbols) {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrCodeWithInvalidSymbol.Error())
	}

	//return i.clientMain.ProductCreate(ctx, in)
	err := i.kafkaRequester.ProductCreate(kafkaConfig.ProductCreateRequest{
		Code: in.GetCode(),
		Name: in.GetName(),
	})

	if err != nil {
		logrus.Errorf("ProductCreate error during sending kafka request: %v", err.Error())
		return nil, err
	}

	return &pb.ProductCreateResponse{}, nil
}

func (i *implementation) ProductList(ctx context.Context, in *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	if in.GetOrderBy() != pb.ProductListRequest_code && in.GetOrderBy() != pb.ProductListRequest_name {
		return &pb.ProductListResponse{}, status.Error(codes.InvalidArgument, "Not implemented order_by")
	}

	if in.GetResultsPerPage() == 0 {
		in.ResultsPerPage = config.DefaultResultsPerPage
	}

	err := i.kafkaRequester.ProductList(kafkaConfig.ProductListRequest{
		PageNumber:     in.GetPageNumber(),
		ResultsPerPage: in.GetResultsPerPage(),
		OrderBy:        int32(in.GetOrderBy()),
	})

	if err != nil {
		logrus.Errorf("ProductList error during sending kafka request: %v", err.Error())
		return nil, err
	}

	//return i.clientMain.ProductList(ctx, in)
	return &pb.ProductListResponse{}, nil
}

func (i *implementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	err := i.kafkaRequester.ProductDelete(kafkaConfig.ProductDeleteRequest{
		Code: in.GetCode(),
	})

	if err != nil {
		logrus.Errorf("ProductDelete error during sending kafka request: %v", err.Error())
		return nil, err
	}

	return &pb.ProductDeleteResponse{}, nil
	//return i.clientMain.ProductDelete(ctx, in)
}

func (i *implementation) PriceTimeStampAdd(ctx context.Context, in *pb.PriceTimeStampAddRequest) (*pb.PriceTimeStampAddResponse, error) {
	if in.GetPrice() < 0 {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrNegativePrice.Error())
	}

	if len(in.GetCode()) == 0 {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrEmptyCode.Error())
	}

	err := i.kafkaRequester.PriceTimeStampAdd(kafkaConfig.PriceTimeStampAddRequest{
		Code:  in.GetCode(),
		Price: in.GetPrice(),
		Ts:    in.GetTs(),
	})

	if err != nil {
		logrus.Errorf("PriceTimeStampAdd error during sending kafka request: %v", err.Error())
		return nil, err
	}

	// TODO read from redis
	return &pb.PriceTimeStampAddResponse{}, nil
	//return i.clientMain.PriceTimeStampAdd(ctx, in)
}
func (i *implementation) PriceHistory(ctx context.Context, in *pb.PriceHistoryRequest) (*pb.PriceHistoryResponse, error) {
	err := i.kafkaRequester.PriceHistory(kafkaConfig.PriceHistoryRequest{
		Code: in.GetCode(),
	})

	if err != nil {
		logrus.Errorf("PriceHistory error during sending kafka request: %v", err.Error())
		return nil, err
	}

	// TODO read from redis
	return &pb.PriceHistoryResponse{}, nil
	// return i.clientMain.PriceHistory(ctx, in)
}
