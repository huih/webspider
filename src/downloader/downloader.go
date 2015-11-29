package downloader

import (
	"net/http"
	"logs"
	"io/ioutil"
)

func DownLoad(req *Request) *Response{
	resp, err := http.Get(req.GetUrl())
	if err != nil {
		logs.Error("http get error")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read response body error")
		return nil
	}
	res := &Response{}
	res.SetPageContent(string(body))
	res.SetContentType(resp.Header.Get("content-type"))
	res.SetAcceptLanguage(resp.Header.Get("accept-language"))
	
	return res
}