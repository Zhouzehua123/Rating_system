package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string

}
// Business is a Business model.

// BusinessRepo is a Business repo.
type BusinessRepo interface {
	Reply(context.Context, *ReplyParam)(int64,error)
}

// BusinessUsecase is a Business usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewBusinessUsecase new a Business usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}
//创建一个回复
//service层调用的方法
func (uc *BusinessUsecase) CreateReply(ctx context.Context, param *ReplyParam)(int64,error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply param:%v", param)
	return uc.repo.Reply(ctx, param)
}