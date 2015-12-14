package output

import (

	"io/ioutil"
	"net/url"
	"github.com/gotools/logs"
	"net/http"
	"privatehub" //save private data, for example : user name, user password and so on
)

func PostDataToRemote(paperTitle string, paperContent string, 
	source_title string, source_address string){
	urlValue := url.Values{
		privatehub.WebAddress_Title1:{privatehub.WebAddress_security1},
		privatehub.WebAddress_Title2:{privatehub.WebAddress_security2},
		privatehub.WebAddress_Title3:{paperContent},
		privatehub.WebAddress_Title4:{paperTitle},
		privatehub.WebAddress_Title5:{source_title},
		privatehub.WebAddress_Title6:{source_address},
	}
	resp, err := http.PostForm(privatehub.WebAddress, urlValue)
	if err != nil {
		logs.Debug("post paper err: %s", err.Error())
		return;
	}
	
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug(" post paper err: %s", err.Error())
		return
	}
}