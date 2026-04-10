package service

import (
	"server/global"
	"server/model/database"

	"gorm.io/gorm"
)

// 获取评论下的所有评论
func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment

	// 对评论进行预加载，加载用户信息
	userPreload := func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}
	// 查询评论的子评论
	if err := global.DB.Where("p_id = ?", comment.ID).Preload("User", userPreload).Find(&children).Error; err != nil {
		return err
	}

	// 递归加载子评论的子评论
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}

	comment.Children = children
	return nil
}

// 删除评论下的所有评论
func (commentService *CommentService) DeleteCommentAndChild(tx *gorm.DB, commentId uint) error {
	var children []database.Comment

	// 查询评论是否存在
	if err := tx.Where("p_id", commentId).Find(&children).Error; err != nil {
		return nil
	}

	// 查询子评论是否存在，存在则递归，否则删除该评论
	for _, child := range children {
		if err := commentService.DeleteCommentAndChild(tx, child.ID); err != nil {
			return err
		}
	}

	// 删除评论
	if err := tx.Delete(&database.Comment{}, commentId).Error; err != nil {
		return err
	}

	return nil
}

// 父子评论去重
func (commentService *CommentService) FindChildCommentsIDByRootCommentUserUUID(comments []database.Comment) map[uint]struct{} {
	result := make(map[uint]struct{})

	// 遍历所有根评论
	for _, rootComment := range comments {
		// 创建一个递归函数来查找与根评论相同 UserUUID 的子评论
		var findChildren func([]database.Comment)

		findChildren = func(children []database.Comment) {
			// 遍历当前子评论
			for _, child := range children {
				// 如果子评论的 UserUUID 与根评论相同，加入结果 map
				if child.UserUUID == rootComment.UserUUID {
					result[child.ID] = struct{}{}
				}
				// 如果有子评论，继续递归
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}
		}

		// 调用递归函数，查找根评论的所有子评论
		findChildren(rootComment.Children)
	}

	return result
}
