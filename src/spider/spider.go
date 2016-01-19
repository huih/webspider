package spider

import (
)

//parse rule
type Rule struct{
	//output field
	OutField []string
	ParseFunc func(*Context)
	
}

type RuleTree struct{
	Root func(*Context)
	//parse rule
	ParseFunc map[string]Rule 
}


type Spider struct {
	Id int             // spider id, auto generate by system
	Name string        // spider name
	Description string // spider description
	PauseTime int      // pause time after spider scrable page every time  
	Address string     // spider scrable first address  
	*RuleTree
}

func (self *Spider) Register(){
	SpiderList.AddSpider(self)
}

