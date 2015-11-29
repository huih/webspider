package downloader

import (
	"net/http"
	"logs"
	"io/ioutil"
	"errors"
)

func DownLoad(req *Request) (response *Response, err error){
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
	
	res := &Response{}
	res.SetRequest(req)
	res.SetResponse(resp)
	res.SetText(string(body))
	
	return res, nil
}