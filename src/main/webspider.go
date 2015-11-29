package main

import (
	"logs"
	"downloader"
)

func main(){
	logs.Start()
	logs.Debug("xxxxxxxxxxxxxxx");
	var req downloader.Request
	var resp *downloader.Response
	req.SetUrl("http://blog.csdn.net/experts.html?&page=2") 
	resp, err := downloader.DownLoad(&req)
	if err == nil {
		logs.Debug(*resp.GetText())
	}
}
