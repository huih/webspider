package spider

import(
	"logs"
)

type SpiderInf interface {
	Add(*Spider)
}

type spiderList struct{
	list []*Spider
	sorted bool
}

func newSpiderList() *spiderList{
	return &spiderList {
		list: []*Spider{},
	}
}

var SpiderList = newSpiderList()

func (self *spiderList) AddSpider(spider *Spider) {
	self.list = append(self.list, spider)
}

func ListAllSpiders(){
	l := len(SpiderList.list)
	logs.Debug("xxxxxxxxxxxxxxxxxxxx: %d", l)
	for i := 0; i < l; i++ {
		logs.Debug("%s:%s", SpiderList.list[i].Name, SpiderList.list[i].Address)
	}
}