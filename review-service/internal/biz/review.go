package biz

import (
	"context"
	"review-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义数据库的接口，哪个数据库实现了这个SaveReview方法，该数据库就是这个接口的实现者，因此可以使用不同的数据库进行存储或者测试，例如mysql，MongoDB，甚至是内存数据库等。
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
}

// 定义一个结构体，成员是数据库接口和日志工具
type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

// 构造函数，接收数据库接口和日志工具作为参数
func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// 创建评价，实现业务逻辑的地方
// service 调用该方法
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview req:%#v", review)
	//1.数据校验
	//2.生成评价ID
	//3.查询订单和商品信息
	//4.拼装数据保存到数据库
	return uc.repo.SaveReview(ctx, review)

}
