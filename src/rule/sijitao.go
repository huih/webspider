package rule

// 基础包
import (
//	"regexp"
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
//	. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/gotools/logs"                   //信息输出

	// net包
//	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	 "regexp"
//	"strconv"
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
	"output"
)

var sijitaoDeleteContent = []string {
	`<div[&nbsp;]*style=[\"'a-zA-Z0-9;#]*[&nbsp]*>[\n\r]*src=\"http://crs.baidu.com/t.js?.*</div>`,
	`<!--[\w\W\r\n]*?-->`,
	`<p>[&nbsp;]*</p>`,
	`(\n\r)*`,
}

func init() {
	sijitao.Register()
}

func formatsijitaoContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	
	preDom := bodyContent.Find("pre")
	if preDom != nil {
		for index := 0; index < preDom.Length(); index++{
			dom := preDom.Eq(index)
			dom.RemoveClass("brush:cpp;toolbar: true; auto-links: false;")
			dom.AddClass("prettyprint linenums")
		}
	}
	
	//remove the spicial dom
	if removeDom != nil {
		for _, dom := range removeDom {
			tmpDom := bodyContent.Find(dom)
			if tmpDom != nil {
				tmpDom.Remove()
			}
		}
	}
	return bodyContent
}

var sijitao = &Spider{
	Name:        "sijitao",
	Description: "sijitao博文采集 [http://www.sijitao.net/]",
	Host: "http://www.sijitao.net/",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:      USE,
	EnableCookie: false,
	MaxItemNum: RULE_SPIDER_PAPER_NUM,
	StopSpider: false,
	RemoveDom: []string{"h1.title","div.info","div#singlead","#comments", "div.contentnav", "div.copyinfo", "div.clear", "div[style]", "div.bdlikebutton", "div#bdshare"},
//	FilterDom: []string{},
//	FilterUrl: []string{"depsi", "liwei225"},
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&context.Request{
							Url:   "http://www.sijitao.net/",
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					itemBody := query.Find("div.template")
	
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}
					
					//get detail page url
					for index := 0; index < itemBody.Length() && ctx.Spider.ContinueGetItem(); index++ {
						div := itemBody.Eq(index)
						hrefa := div.Find("h2 a[href]")
						href,_ := hrefa.Attr("href")
						
						var tc = false	
						for _, u := range ctx.Spider.GetFilterUrl() {
							tc = strings.Contains(href, u)
							if tc {
								break
							}
						}
						if tc {
							continue
						}
						
						ctx.AddQueue(&context.Request{
							Url:    href,
							Rule:   "详细页面",
						})
					}
					
					//get next page url
					if (ctx.Spider.GetStopSpider()) {
						logs.Debug("stop spider: %s", ctx.Spider.GetHostUrl())
						return
					}
					
					nextPage := query.Find("div.pagenav #next")
					if nextPage.Length() <= 0 {
						return;
					}
					
					text,_ := nextPage.Attr("href")
					
					if text != "" {
						isLast := strings.Contains(ctx.GetUrl(), text)
						if (isLast == false) && ctx.Spider.ContinueGetItem() {
							nextPageUrl := ctx.Spider.GetHostUrl() + "/blog" + text
							ctx.AddQueue(&context.Request{
							Url:    nextPageUrl,
							Rule:   "生成请求",
							})
						}
					}
				},
			},

			"详细页面": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"paper_title",
					"paper_content",
					"web_title",
					"paper_address",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					
					//find title
					titleDom := query.Find("div.template h1.title")
					
					
					paperTitle := titleDom.Text()
					paperTitle = strings.TrimSpace(paperTitle)
					
					//find content					
					bodyContent := query.Find("div.template")
					
					//filter bad paper
					for _, dom := range ctx.Spider.GetFilterDom() {
						tmpDom := bodyContent.Find(dom)
						if tmpDom.Length() > 1 {
							return
						}
					}
					
					bodyContent = formatsijitaoContent(ctx.Spider.GetRemoveDom(), bodyContent)
					bodyHtml, err := bodyContent.Html()
					if err != nil || bodyHtml == "" {
						logs.Debug("get body content html err: %s", err.Error())
						return;
					}
					//remove bad content
					for _, item := range sijitaoDeleteContent {
						reg := regexp.MustCompile(item)
						bodyHtml = reg.ReplaceAllLiteralString(bodyHtml,"")
					}
					paperTitle = strings.TrimSpace(paperTitle);
					bodyHtml = strings.TrimSpace(bodyHtml);				
						
//					save result to local db for replication
					output.SaveDataToLocalDB(paperTitle, bodyHtml, ctx.Spider.GetName(), ctx.GetUrl())
					
					if IsAllCode(bodyContent) == true {
						logs.Debug("current page is all code")
						return;
					}
							
					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: paperTitle,
						1: bodyHtml,
						2: ctx.Spider.GetName(),
						3: ctx.GetUrl(),
					})
				},
			},
		},
	},
}
