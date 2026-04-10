package service

import (
	"server/global"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gofrs/uuid"
)

type FeedbackService struct {
}

func (feedbackService *FeedbackService) FeedbackCreate(req request.FeedbackCreate) error {
	return global.DB.Create(&database.Feedback{
		UserUUID: req.UUID,
		Content:  req.Content,
	}).Error
}

func (feedbackService *FeedbackService) FeedbackInfo(uuid uuid.UUID) (feedbacks []database.Feedback, err error) {
	if err := global.DB.Model(&database.Feedback{}).Where("user_uuid  = ?", uuid).Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (feedbackService *FeedbackService) FeedbackNew() (feedbacks []database.Feedback, err error) {
	// 按id倒序排序，查询5条并赋给feedbacks
	if err := global.DB.Order("id desc").Limit(5).Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (feedbackService *FeedbackService) FeedbackDelete(req request.FeedbackDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Delete(&database.Feedback{}, req.IDs).Error
}

func (feedbackService *FeedbackService) FeedbackReply(req request.FeedbackReply) error {
	return global.DB.Take(&database.Feedback{}, req.ID).Update("reply", req.Reply).Error
}

func (feedbackService *FeedbackService) FeedbackList(req request.PageInfo) (list interface{}, total int64, err error) {
	option := other.MySQLOption{
		PageInfo: req,
	}

	return utils.MySQLPagination(&database.Feedback{}, option)
}
