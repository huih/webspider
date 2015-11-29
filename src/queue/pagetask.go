package tasklist

import (
	"container/list"
)

var pageList *list.List;
type PageObject struct{
	PageUrl string;
	HostUrl string;
}

func AddPageTask(pageUrl string, hostUrl string) {
	var object PageObject
	object.PageUrl = pageUrl
	object.HostUrl = hostUrl
	
	if pageList == nil {
		pageList = list.New()
	}
	pageList.PushBack(object)
}

func PopPageTask() PageObject {
	if pageList.Len() <= 0 {
		return PageObject{}
	}
	e := pageList.Front()
	return pageList.Remove(e).(PageObject)
}
