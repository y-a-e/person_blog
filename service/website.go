package service

import (
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type WebsiteService struct {
}

func (websiteService *WebsiteService) WebsiteAddCarousel(req request.WebsiteCarouselOperation) error {
	return utils.ChangeImagesCategory(global.DB, []string{req.Url}, appTypes.Carousel)
}

func (websiteService *WebsiteService) WebsiteCancelCarousel(req request.WebsiteCarouselOperation) error {
	return utils.InitImagesCategory(global.DB, []string{req.Url})
}

func (websiteService *WebsiteService) WebsiteCarousel() ([]string, error) {
	var urls []string
	// Pluck 用于从查询结果中提取指定列的值到一个切片中
	if err := global.DB.Model(&database.Image{}).Where("category = ?", appTypes.Carousel).Pluck("url", &urls).Error; err != nil {
		return []string{}, err
	}
	return urls, nil
}

func (websiteService *WebsiteService) WebsiteNews(sourceStr string) (other.HotSearchData, error) {
	hotSearchData, err := ServiceGroupApp.HotSearchService.GetHotSearchDataBySource(sourceStr)
	if err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}

func (websiteService *WebsiteService) WebsiteCalendar(sourceStr string) (other.Calendar, error) {
	calendarData, err := ServiceGroupApp.CalendarService.GetCalendarByDate(sourceStr)
	if err != nil {
		return other.Calendar{}, err
	}
	return calendarData, nil
}

func (websiteService *WebsiteService) WebsiteFooterLink() []database.FooterLink {
	var footerLinkList []database.FooterLink
	global.DB.Model(&database.FooterLink{}).Find(&footerLinkList)
	return footerLinkList
}

func (websiteService *WebsiteService) WebsiteCreateFooterLink(req database.FooterLink) error {
	return global.DB.Save(&req).Error
}

func (websiteService *WebsiteService) WebsiteDeleteFooterLink(req database.FooterLink) error {
	return global.DB.Delete(&req).Error
}
