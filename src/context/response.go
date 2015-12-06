package context

import (
	"net/http"
)
type Response struct {	
	//the request is crawled by spider that contains url and relevent information 
	request *Request
	
	//crawl result
	response *http.Response
	
	//the text is body of response
	text string
}

func (self *Response) GetRequest() *Request{
	return self.request
}

func (self *Response) SetRequest(req *Request) {
	self.request = req
}

func (self *Response) GetResponse() *http.Response {
	return self.response
}

func (self *Response) SetResponse(resp *http.Response) {
	self.response = resp
}

func (self *Response) GetText() string {
	return self.text
}

func (self *Response) SetText(text string) {
	self.text = text
}