package output

import (
	"strings"
	"errors"
	"regexp"
	"io/ioutil"
	"net/url"
	"github.com/gotools/logs"
	"net/http"
	"privatehub" //save private data, for example : user name, user password and so on
)

func PostDataToRemote(paperTitle string, paperContent string, 
	source_title string, source_address string){
	
	if len(paperContent) < 200 {//remove short paper
		return
	}
	
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

func PostImageToRemote(imagePath string) (string, error){
	urlValue := url.Values{
		privatehub.WebAddress_Title1:{privatehub.WebAddress_security1},
		privatehub.WebAddress_Title2:{privatehub.WebAddress_security2},
		privatehub.WebAddress_Title7:{imagePath},
	}
	resp, err := http.PostForm(privatehub.WebAddress_image, urlValue)
	if err != nil {
		logs.Debug("post paper err: %s", err.Error())
		return "", errors.New("post err");
	}
	
	defer resp.Body.Close()
	imgPath, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug(" post paper err: %s", err.Error())
		return "", errors.New("post err");
	}
	
	tmpImagePath := string(imgPath)
	reg := regexp.MustCompile("http://.*.(png|jpg|jpeg|gif)")
	tmpImagePath = reg.FindString(tmpImagePath)
	
	 
	return tmpImagePath, nil
}

func HandlePaperImagePath(paperContent string) string {
	reg := regexp.MustCompile(`http://[\S]*.(png|jpg|jpeg|gif)`)
	imageUrlArray := reg.FindAllString(paperContent, -1)
	
	for _, imagePath := range imageUrlArray {
		newImagePath, err := PostImageToRemote(imagePath)
		if err != nil {
			return "";
		}
		paperContent = strings.Replace(paperContent, imagePath, newImagePath, -1)
	}
	return paperContent
}

func HandleOutputData(paperTitle string, paperContent string, 
	source_title string, source_address string) {
	
	//find all url image
	paperContent = HandlePaperImagePath(paperContent)
	if paperContent == "" {
		return;
	}
	
	PostDataToRemote(paperTitle, paperContent, source_title, source_address)	
}

func Auto_add_paper_tags(){
	urlValue := url.Values{
		privatehub.WebAddress_Title1:{privatehub.WebAddress_security1},
		privatehub.WebAddress_Title2:{privatehub.WebAddress_security2},
		privatehub.WebAddress_Title8:{privatehub.WebAddress_Method_Value1},
	}
	resp, err := http.PostForm(privatehub.WebAddress_address_url3, urlValue)
	if err != nil {
		logs.Debug("post paper err: %s", err.Error())
		return;
	}
	
	defer resp.Body.Close()
	resContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug(" post paper err: %s", err.Error())
		return
	}
	logs.Debug(string(resContent))
}

func Auto_submit_paper_url(){
	urlValue := url.Values{
		privatehub.WebAddress_Title1:{privatehub.WebAddress_security1},
		privatehub.WebAddress_Title2:{privatehub.WebAddress_security2},
		privatehub.WebAddress_Title8:{privatehub.WebAddress_Method_Value2},
	}
	resp, err := http.PostForm(privatehub.WebAddress_address_url3, urlValue)
	if err != nil {
		logs.Debug("post paper err: %s", err.Error())
		return;
	}
	
	defer resp.Body.Close()
	resContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Debug(" post paper err: %s", err.Error())
		return
	} else {
		logs.Debug(string(resContent))
	}
}