package downloader

import (
	"net/http"
	"logs"
	"io/ioutil"
	"errors"
	"task/queue"
	"context"
)

func DownLoad(req *context.Request) (response *context.Response, err error){
	resp, err := http.Get(req.GetUrl())
	if err != nil {
		logs.Error("http get error")
		return nil, errors.New("http get error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read response body error")
		return nil, errors.New("read response body error")
	}
	
	res := &context.Response{}
	res.SetRequest(req)
	res.SetResponse(resp)
	res.SetText(string(body))
	
	return res, nil
}

func downloadPage(){
	//get task from taskqueue
	var task = taskqueue.DownloadTask.PopOneTask()
	logs.Debug("address: %s", task.Spider.Address)
	
	resp, ok := DownLoad(&task.Request)
	if ok != nil {
		logs.Debug("download error(%s)", ok.Error());
	}
	task.Response = resp
	taskqueue.PageQueue.PushOneTask(task)
	logs.Debug("text: %s", task.Response.GetText())
}


func Start(){
	downloadPage()
}