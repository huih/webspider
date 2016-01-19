package main

import (
//	"io/ioutil"
//	"strings"
//	"net/http"
//	"regexp"
//	"strings"
//	"io/ioutil"
//	"fmt"
//	"strings"
//	"regexp"
//	"os"
//	"fmt"
	"output"
//	"github.com/gotools/logs"
//	"github.com/PuerkitoBio/goquery"
//	"rule"
)

//psql -h 127.0.0.1 -p 5432 -U test postgres
func main(){
//	papter_title := "testtest"
//	papter_content := "testlkjdfksadfkljsdakfjsdkafjkdsaljfkds"
//	source_title := "zhiliaoyuan"
//	source_address := "http://www.zhiliaoyuan.com"
//	output.PostDataToRemote(papter_title, papter_content, source_title, source_address)
//	imagePath := "http://images2015.cnblogs.com/blog/464220/201512/464220-20151218091553412-243296128.png";
//	tmpPath, err := output.PostImageToRemote(imagePath);
//	if err != nil {
//		fmt.Printf("%s", err.Error())
//		return;
//	}
//	fmt.Println(tmpPath)

//	f, err := ioutil.ReadFile("log.txt")
//	if err != nil {
//		return;
//	}
//	content := string(f)
//	fmt.Print(content)
	
//	reg := regexp.MustCompile("http://.*.(png|jpg|jpeg|gif)")
//	imgPath := reg.FindAllString(content, -1)
//	fmt.Println(imgPath)
//	for _, value := range imgPath {
//		newValue := value + "123.jpg"
//		content = strings.Replace(content, value, newValue, -1)
//	}
//	fmt.Println(content)

//	paper_title := "test"
//	paper_content := "test"
//	source_title := "test"
//	source_url := "test"
//	output.SaveDataToLocalDB(paper_title, paper_content, source_title, source_url)
//	rst, err := output.DataInDB("http://www.baidu.com")
//	if err != nil {
//		logs.Debug("err ")
//	} else {
//		logs.Debug("rst: %d", rst)
//	}
//{
//	content, err := ioutil.ReadFile("log3.txt")
//	if err != nil {
//		logs.Debug("read all file data")
//		return;
//	}
//	query, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
//	//logs.Debug(query.Html())
//	toolbar := query.Find("span.cnblogs_code_collapse")
//	toolbar.Remove()
//	img := query.Find("img.code_img_opened")
//	if img != nil {
//		img.Remove()
//	}
	
//	img = query.Find("img.code_img_closed")
//	if img != nil {
//		img.Remove()
//	}
//	logs.Debug(query.Html())
//}

//{
//	content, err := ioutil.ReadFile("log3.txt")
//	if err != nil {
//		logs.Debug("read all file data")
//		return;
//	}
//	query, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
//	dom := query.Find("#cnblogs_post_body")
//	//logs.Debug(query.Html())
//	removeDom := []string{"div.cnblogs_code_toolbar", "img.code_img_closed", "img.code_img_opened","span.cnblogs_code_collapse"}
//	dom = rule.FormatContent(removeDom, dom)
//	logs.Debug(dom.Html())
//}

//{
//	content, _ := ioutil.ReadFile("log.txt")
//	reg := regexp.MustCompile(`<span style=\"color: #008080;?\">\s*(&nbsp;)?(\d+)<\/span>`)
//	tmp := reg.ReplaceAllLiteralString(string(content), "")
//	logs.Debug(tmp)
//}
	
//{
//	content, _ := ioutil.ReadFile("log3.txt")
//	tmp := output.HandlePaperImagePath(string(content))
//	logs.Debug(tmp)
//}

//var deleteContent = []string {
//	`<span style=\"color: #008080;?\">\s*(&nbsp;)?(\d+)<\/span>`,
//	`<table border=\"0\".*>(\n\r)*<tbody><tr>`,
//	`<!--[\w\W\r\n]*?-->`,
//	`<p>[&nbsp;]*</p>`,
//	`(\n\r)*`,
//}
//{
//	bodyHtml, _ := ioutil.ReadFile("log3.txt")
//	bodyContent := string(bodyHtml)
//	for _, item :=  range deleteContent {
//		reg := regexp.MustCompile(item)
//		bodyContent = reg.ReplaceAllLiteralString(bodyContent,"")
//	}
//	logs.Debug(bodyContent)
//}
//{
//	query, _ := goquery.NewDocument("http://www.ibm.com/developerworks/cn/linux/l-processor-utilization-difference-aix-lop-trs/index.html")
	
//	listContent := query.Find("div.SpaceList")
//	entityContent := listContent.Find("div.BlogEntity")					
//	bodyContent := entityContent.Find("div.BlogContent")
					
//	docstr, _ := bodyContent.Html()
//	logs.Debug(docstr)
//}

//{
//	url := "http://www.ibm.com/developerworks/cn/linux/1511_cyq_tool/index.html"
//	client := &http.Client{}
	
//	req, err := http.NewRequest("GET", url, strings.NewReader("name=cjb"))
//	if err != nil {
//		logs.Debug("xxxxxxxxxxxx")
//		return;
//	}
	
//	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8");
//	//req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
//	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
//	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
////	req.Header.Set("Cookie", "PD_STATEFUL_4016dd7c-0b03-11e5-bbb7-a628a0b3ab03=%2Fgateway; mm_session=1; cmTPSet=Y; CMAVID=none; CookieChecker=set; ibmSurvey=1452992727406; UnicaNIODID=JLzqLSXW2nC-ZkAZZZU; PD_STATEFUL_dfec5590-0b01-11e5-9c53-62ee8e3f6b03=%2Fgateway; CoreM_State=1~-1~-1~-1~-1~3~3~5~3~3~7~7~|~~|~~|~~|~||||||~|~~|~~|~~|~~|~~|~~|~~|~; CoreM_State_Content=6~|~~|~|; PrefID=242-108914777; optimizelyEndUserId=oeu1452992801165r0.4865624117664993; optimizelySegments=%7B%223512541636%22%3A%22none%22%2C%223532581527%22%3A%22false%22%2C%223535552913%22%3A%22direct%22%2C%223539991679%22%3A%22gc%22%7D; optimizelyBuckets=%7B%7D; mmcore.tst=0.016; mmid=-31383215%7CBgAAAApE1atT5wwAAA%3D%3D; mmcore.pd=-1670244131%7CBgAAAAoBQkTVq1PnDLhf2LEBABkWWoLaHtNIDwAAAFWgLDXaHtNIAAAAAP//////////AAZEaXJlY3QB5wwBAAAAAAAAAAAAAP+MAACAigAA/4wAAAAAAAAAAUU%3D; mmcore.srv=lvsvwcgus02; mm_criteria_tt=%7B%22TrafficType%22%3A%22NotDriveToTraffic%22%7D; CoreID6=33638109469014529927172&ci=50200000|GLOBAL_50200000|IBMTESTWWW_50200000|MKTIBMCLOUD_50200000|DEVWRKS; pSite=http%3A%2F%2Fwww.ibm.com%2Fdeveloperworks%2Fcn%2Flinux%2Fl-processor-utilization-difference-aix-lop-trs%2Findex.html; Hm_lvt_7cb65dbe40a6c02096dd8f8dbf47a19f=1452847901,1452992716,1452992779,1452992895; Hm_lpvt_7cb65dbe40a6c02096dd8f8dbf47a19f=1452992911; __asc=bd34e0771524d1d600dbe8eb3cb; __auc=bd34e0771524d1d600dbe8eb3cb; JSESSIONID=0000DST0B5LeaF_FqewRRpMBqCr:17bqvrk0o; utag_main=v_id:01524d1d5b2b001fc5a90012059c2306c001906400876$_sn:1$_ss:0$_pn:11%3Bexp-session$_st:1452994709910$ses_id:1452992715563%3Bexp-session$dc_visit:1$dc_event:11%3Bexp-session$dc_region:eu-west-1%3Bexp-session; 50200000_clogin=v=1&l=1452992717&e=1452994714160")
	
//	resp, err := client.Do(req)
//	defer resp.Body.Close()
	
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logs.Debug("xxxxxxxxxxxxxxxx")
//	}
//	logs.Debug(string(body))
	
	
	
	
//}


	//output.Auto_add_paper_tags()
	output.Auto_submit_paper_url()
}
