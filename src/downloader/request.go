package downloader

type Request struct{
	downloaderId int
	spider string
	url string
}

func (self *Request) GetDownloaderId() int {
	return self.downloaderId
}

func (self *Request) SetDownloaderId(downloaderId int) {
	self.downloaderId = downloaderId
}

func (self *Request) GetSpider() string {
	return self.spider
}

func (self *Request) SetSpider(spiderName string) {
	self.spider = spiderName
}

func (self *Request) GetUrl() string {
	return self.url
}

func (self *Request) SetUrl(url string) {
	self.url = url
}