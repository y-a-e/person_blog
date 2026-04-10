package service

import (
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"gorm.io/gorm"
)

type AdvertisementService struct {
}

func (advertisementService *AdvertisementService) AdvertisementCreate(req request.AdvertisementCreate) error {
	advertisementToCreate := database.Advertisement{
		AdImage: req.AdImage,
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.ChangeImagesCategory(tx, []string{advertisementToCreate.AdImage}, appTypes.AdImage); err != nil {
			return err
		}

		return tx.Create(&advertisementToCreate).Error
	})
}

func (advertisementService *AdvertisementService) AdvertisementDelete(req request.AdvertisementDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var advertisementDelete database.Advertisement
			if err := tx.Take(&advertisementDelete, id).Error; err != nil {
				return err
			}

			if err := utils.InitImagesCategory(tx, []string{advertisementDelete.AdImage}); err != nil {
				return err
			}

			if err := tx.Delete(&advertisementDelete).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (advertisementService *AdvertisementService) AdvertisementUpdate(req request.AdvertisementUpdate) error {
	advertisementUpdate := struct {
		Link    string `json:"link" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}{
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}

	return global.DB.Take(&database.Advertisement{}, req.ID).Updates(advertisementUpdate).Error
}

func (advertisementService *AdvertisementService) AdvertisementList(info request.AdvertisementList) (list interface{}, total int64, err error) {
	db := global.DB

	if info.Title != nil {
		db = db.Where("title LIKE ?", "%"+*info.Title+"%")
	}
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Advertisement{}, option)
}

func (advertisementService *AdvertisementService) AdvertisementInfo() (ads []database.Advertisement, total int64, err error) {
	if err := global.DB.Model(&database.Advertisement{}).Count(&total).Find(&ads).Error; err != nil {
		return nil, 0, err
	}
	return ads, total, nil
}
