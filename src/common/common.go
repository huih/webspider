package common

import (
	"errors"
	"regexp"
	"strings"
	"github.com/gotools/logs"
	"github.com/PuerkitoBio/goquery"
	. "github.com/henrylee2cn/pholcus/app/spider"
)

func IsAllCode(bodyContent *goquery.Selection) bool {
	preDom := bodyContent.Find("pre")
	if preDom != nil {
		for index := 0; index < preDom.Length(); index++{
			dom := preDom.Eq(index)
			dom.Remove()
		}
	}
	conStr := bodyContent.Text()
	
	conStr = strings.TrimSpace(conStr)
	if len(conStr) < 100 {
		return true
	}
	return false
}


func formatCommonContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	
	preDom := bodyContent.Find("pre")
	if preDom != nil {
		for index := 0; index < preDom.Length(); index++{
			dom := preDom.Eq(index)
			dom.RemoveClass("brush:cpp;toolbar: true; auto-links: false;")
			dom.AddClass("prettyprint linenums")
		}
	}
	
	tableDom := bodyContent.Find("table")
	if tableDom != nil {
		for index := 0; index < tableDom.Length(); index++{
			dom := tableDom.Eq(index)
			dom.RemoveAttr("border")
			dom.SetAttr("border", "1")
		}
	}
	
	//remove the spicial dom
	if removeDom != nil {
		for _, dom := range removeDom {
			tmpDom := bodyContent.Find(dom)
			if tmpDom != nil {
				tmpDom.Remove()
			}
		}
	}
	return bodyContent
}

//<img alt="图 11. 导入编译好的内核模块" src="img011.png" width="555"/>
func FindImagePath(content string) []string{
	
	rule := `src[\s|&nbsp]*=[\s|&nbsp]*['|\"]*[\S]*.(png|jpg|jpeg|gif)['|\"]`
	
	reg := regexp.MustCompile(rule)
	imageUrlArray := reg.FindAllString(content, -1)
	
	for index, imagePath := range imageUrlArray {
		splitStr := strings.Split(imagePath, "\"")
		if len(splitStr) < 3 {
			splitStr = strings.Split(imagePath, "'")
		}
		
		if len(splitStr) < 3 {
			return []string{}
		}
		
		imageUrlArray[index] = splitStr[1]
	}
	return imageUrlArray
}

func UpdateImagePath(content string, host string) string {
	imagePathArray := FindImagePath(content)
	for _, imagePath := range imagePathArray {
//		logs.Debug(imagePath)
		if strings.Contains(imagePath, "http://") {
			break
		}
		
		newImagePath := host + "/" + imagePath
//		logs.Debug(newImagePath)
		content = strings.Replace(content, imagePath, newImagePath, -1)
	}
	return content
}


func GetHtmlContent(query *goquery.Document, path string, ctx *Context, deleteArray []string, isHtml bool) (string, error) {
		contentHtml := ""
		err := errors.New("no error")
		
		contentObj := query.Find(path)
		
		for _, dom := range ctx.Spider.GetFilterDom() {
			tmpDom := contentObj.Find(dom)
			if tmpDom.Length() > 1 {
				return "", errors.New("filter the dom")
			}
		}
		
		contentObj = formatCommonContent(ctx.Spider.GetRemoveDom(), contentObj)
		
		if isHtml {
			contentHtml, err = contentObj.Html()
			if err != nil {
				logs.Debug("get html error, err:%s", err.Error())
				return "", err
			}
			
			//remove bad content
			for _, item := range deleteArray {
				reg := regexp.MustCompile(item)
				contentHtml = reg.ReplaceAllLiteralString(contentHtml,"")
			}
		} else {
			contentHtml = contentObj.Text()
		}
		
		contentHtml = strings.TrimSpace(contentHtml)
		
		return contentHtml, nil
}

func GetPageUrlDirectoryPath(url string) string {
	url = url[:strings.LastIndex(url, "/")]
	return url
}