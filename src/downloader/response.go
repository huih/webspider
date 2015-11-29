package downloader

type Response struct {
	pageContent string
	acceptLanguage string
	cookie string
	contentType string
}

func (self *Response) SetPageContent(pageContent string) {
	self.pageContent = pageContent
}

func (self *Response) GetPageContent() string{
	return self.pageContent
}

func (self *Response) SetAcceptLanguage(acceptLanguage string) {
	self.acceptLanguage = acceptLanguage
}

func (self *Response) GetAcceptLanguage() string {
	return self.acceptLanguage
}

func (self *Response) SetContentType(contentType string) {
	self.contentType = contentType
}

func (self *Response) GetContentType() string{
	return self.contentType
}