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
	resp = downloader.DownLoad(&req)
	logs.Debug(resp.GetPageContent())
}
