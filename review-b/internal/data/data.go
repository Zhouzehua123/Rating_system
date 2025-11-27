package data

import (
	"context"
	v1 "review-b/api/review/v1"
	"review-b/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewReviewServiceClient, NewData, NewBusinessRepo)

// Data .
type Data struct {
	//之前是嵌入一个query对象，现在是要嵌入一个grpc的client端，通过这个client去调用review-service服务
	rc  v1.ReviewClient
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, rc v1.ReviewClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		rc:  rc,
		log: log.NewHelper(logger),
	}, cleanup, nil
}

// 创建一个连接 review-service服务的grpc client端
func NewReviewServiceClient() v1.ReviewClient {
	//这里要实现一个grpc的client去连接review-service服务
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9092"),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewReviewClient(conn)
}
