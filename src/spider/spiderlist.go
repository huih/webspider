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

func (self *spiderList) Size() int {
	return len(self.list)
}

func (self *spiderList) Get(i int) *Spider{
	if i < self.Size() {
		return self.list[i]
	} else {
		return nil
	}
}

func ListAllSpiders(){
	l := len(SpiderList.list)
	logs.Debug("xxxxxxxxxxxxxxxxxxxx: %d", l)
	for i := 0; i < l; i++ {
		logs.Debug("%s:%s", SpiderList.list[i].Name, SpiderList.list[i].Address)
	}
}