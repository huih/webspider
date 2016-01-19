package rule

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
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
