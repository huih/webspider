package spider

import (
	"logs"
)

type Spider struct {
	Id int             // spider id, auto generate by system
	Name string        // spider name
	Description string // spider description
	PauseTime int      // pause time after spider scrable page every time  
	Address string     // spider scrable first address  
}

func init(){
	logs.Debug("xxxxxxxxxxxxxx")
}

func (self *Spider) Register(){
	SpiderList.AddSpider(self)
}

