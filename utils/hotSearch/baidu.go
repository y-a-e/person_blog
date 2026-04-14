package hotSearch

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"server/model/other"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type Baidu struct {
}

func (*Baidu) GetHotSearchData(maxNum int) (HotSearchData other.HotSearchData, err error) {
	// 可尝试使用postman获取网站数据，里面包括热搜
	resp, err := http.Get("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		return other.HotSearchData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.HotSearchData{}, err
	}

	var jsonStr string
	// <!--s-data:{"data":{"cards":[{"updateTime":1234567890,"content":[{"word":"热搜1","desc":"描述1"}]}]}}-->
	// reg = {"data":{"cards":[{"updateTime":1234567890,"content":[{"word":"热搜1","desc":"描述1"}]}]}}
	reg := regexp.MustCompile(`<!--s-data:({.*?})-->`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	// result[0][0] = "<!--s-data:{"name":"test"}-->", result[0][1] = "{"name":"test"}"
	if len(result) > 0 && len(result[0]) > 1 {
		jsonStr = result[0][1]
	} else {
		return other.HotSearchData{}, errors.New("failed to get data")
	}

	updateTime := time.Unix(gjson.Get(jsonStr, "data.cards.0.updateTime").Int(), 0).Format("2006-01-02 15:04:05")

	var hotList []other.HotItem
	for i := 0; i < maxNum; i++ {
		if index := gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".index"); !index.Exists() {
			break
		}
		hotList = append(hotList, other.HotItem{
			Index:       i + 1,
			Title:       gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".word").Str,
			Description: gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".desc").Str,
			Image:       gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".img").Str,
			Popularity:  gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".hotScore").Str,
			URL:         gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".rawUrl").Str,
		})
	}

	return other.HotSearchData{Source: "百度热搜", UpdateTime: updateTime, HotList: hotList}, nil
}
