package api

import (
	"context"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/config"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/error_codes"
	pb "gitlab.ozon.dev/N0fail/price-tracker-validator/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(clientMain pb.AdminClient) pb.AdminServer {
	return &implementation{
		clientMain: clientMain,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	clientMain pb.AdminClient
}

func (i *implementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	if in.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrEmptyCode.Error())
	}

	if len(in.GetName()) < config.MinNameLength {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrNameTooShortError.Error())
	}

	return i.clientMain.ProductCreate(ctx, in)
}

func (i *implementation) ProductList(ctx context.Context, in *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	if in.GetOrderBy() != pb.ProductListRequest_code && in.GetOrderBy() != pb.ProductListRequest_name {
		return &pb.ProductListResponse{}, status.Error(codes.InvalidArgument, "Not implemented order_by")
	}

	if in.GetResultsPerPage() == 0 {
		in.ResultsPerPage = config.DefaultResultsPerPage
	}

	return i.clientMain.ProductList(ctx, in)
}

func (i *implementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	return i.clientMain.ProductDelete(ctx, in)
}

func (i *implementation) PriceTimeStampAdd(ctx context.Context, in *pb.PriceTimeStampAddRequest) (*pb.PriceTimeStampAddResponse, error) {
	if in.GetPrice() < 0 {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrNegativePrice.Error())
	}

	if len(in.GetCode()) == 0 {
		return nil, status.Error(codes.InvalidArgument, error_codes.ErrEmptyCode.Error())
	}

	return i.clientMain.PriceTimeStampAdd(ctx, in)
}
func (i *implementation) PriceHistory(ctx context.Context, in *pb.PriceHistoryRequest) (*pb.PriceHistoryResponse, error) {
	return i.clientMain.PriceHistory(ctx, in)
}
