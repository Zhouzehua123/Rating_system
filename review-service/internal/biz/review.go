package biz

import (
	"context"
	"fmt"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义数据库的接口，哪个数据库实现了这个SaveReview方法，该数据库就是这个接口的实现者，因此可以使用不同的数据库进行存储或者测试，例如mysql，MongoDB，甚至是内存数据库等。
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
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
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview req:%#v", review)
	//1.数据校验
	// 1.1 参数基础校验：正常来说不应该放在这一层，你在上一层或者框架层都应该能拦住
	// 1.2 参数业务校验：带业务逻辑的参数校验，比如已评价过的订单不能再创建评价
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) > 0 {
		// 已经评价过
		fmt.Printf("订单已评价len(reviews):%d\n", len(reviews))
		return nil, v1.ErrorOrderReviewed("订单:%d已评价", review.OrderID)
	}
	//2.生成评价ID
	// 这里可以使用雪花算法自己生成
	// 也可以直接接入公司内部的分布式ID生成服务 (前提是公司内部有这个服务)
	review.ReviewID = snowflake.GenID()
	//3.查询订单和商品快照信息
	//实际业务场景下需要查询订单服务和商家服务，比如通过RPC调用订单服务和商家服务
	//4.拼装数据保存到数据库
	return uc.repo.SaveReview(ctx, review)

}

// 创建评价回复
func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	//调用data层创建一个评价的回复
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply param:%v", param)
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}
