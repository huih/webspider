package taskqueue

import (
	"spider"
//	"logs"
	"context"
	"container/list"
	"sync"
)

type Task struct{
	*spider.Spider
	Request context.Request
	Response *context.Response
}

type TaskQueue struct{
	mu sync.Mutex
	l list.List
}
var DownloadTask = &TaskQueue{}
var PageQueue = &TaskQueue{}
 
func (self *TaskQueue) PushOneTask(task *Task) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.l.PushBack(task)
}

func (self *TaskQueue) Empty() bool {
	if self.l.Len() > 0 {
		return false
	} else {
		return true
	}
}

func (self *TaskQueue) PopOneTask() *Task {
	self.mu.Lock()
	defer self.mu.Unlock()
	task := self.l.Front()
	return self.l.Remove(task).(*Task)
}

func (self *TaskQueue) Size() int {
	return self.l.Len()
}

func init(){
	for i := 0; i < spider.SpiderList.Size(); i++ {
		task := &Task{}
		task.Spider = spider.SpiderList.Get(i)
		task.Request.SetDownloaderId(i)
		task.Request.SetSpider(task.Spider.Name)
		task.Request.SetUrl(task.Spider.Address)
		DownloadTask.PushOneTask(task)
	}
}