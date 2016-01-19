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

var oschinaDeleteContent = []string {
	`<!--[\w\W\r\n]*?-->`,
	`<p>[&nbsp;]*</p>`,
	`<span id=\"OSC_h[0-9]+_[0-9]+\"></span>`,
	`(\n\r)*`,
}

func init() {
	oschina.Register()
}

func formatOschinaContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	
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

var oschina = &Spider{
	Name:        "oschina",
	Description: "oschina博文采集 [http://www.oschina.net/blog]",
	Host: "http://www.oschina.net",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:      USE,
	EnableCookie: false,
	StopSpider: false,
	MaxItemNum: RULE_SPIDER_PAPER_NUM,
	RemoveDom:[]string{"div.toolbar", "td.gutter"},
//	FilterDom: []string{"#codeSnippetWrapper","meta","pre.best-text mb-10"},
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
							Url:   "http://www.oschina.net/blog",
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					itemBody := query.Find("#OSC_Screen #OSC_Content #BlogChannel #BlogsList #RecentBlogs ul.BlogList li")
	
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}

					//get detail page url
					for index := 0; index < itemBody.Length() && ctx.Spider.ContinueGetItem(); index++ {
						div := itemBody.Eq(index)
						hrefa := div.Find("h3 a[href]")
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
					
					//get next page url
					allPage := query.Find("ul.pager")
					pageList := allPage.Children()
					if pageList.Length() <= 0 {
						return
					}
					nextPage := pageList.Eq(pageList.Length() - 1)
					nextPageItem := nextPage.Find("a[href]")
					text,_ := nextPageItem.Attr("href")
					
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
					titleDom := query.Find("div.BlogTitle")
					titleDomText := titleDom.Find("h1")
					
					paperTitle := titleDomText.Text()
					if (strings.Contains(paperTitle, "转")){
						return;
					}
					
					paperTitle = strings.Replace(paperTitle, "原", "",-1)
					paperTitle = strings.Replace(paperTitle, "荐", "",-1)
					paperTitle = strings.Replace(paperTitle, "顶", "",-1)
					paperTitle = strings.TrimSpace(paperTitle)
					
					//find content					
					listContent := query.Find("div.SpaceList")
					entityContent := listContent.Find("div.BlogEntity")					
					bodyContent := entityContent.Find("div.BlogContent")
					//filter bad paper
					for _, dom := range ctx.Spider.GetFilterDom() {
						tmpDom := bodyContent.Find(dom)
						if tmpDom.Length() > 1 {
							return
						}
					}
					
					bodyContent = formatOschinaContent(ctx.Spider.GetRemoveDom(), bodyContent)
					bodyHtml, err := bodyContent.Html()
					if err != nil || bodyHtml == "" {
						logs.Debug("get body content html err: %s", err.Error())
						return;
					}
					//remove bad content
					for _, item := range oschinaDeleteContent {
						reg := regexp.MustCompile(item)
						bodyHtml = reg.ReplaceAllLiteralString(bodyHtml,"")
					}
					paperTitle = strings.TrimSpace(paperTitle);
					bodyHtml = strings.TrimSpace(bodyHtml);				
							
//					//save result to local db for replication
					output.SaveDataToLocalDB(paperTitle, bodyHtml, ctx.Spider.GetName(), ctx.GetUrl())
					
					if IsAllCode(bodyContent) == true {
						logs.Debug("current page is all code")
						return;
					}		
					
//					// 结果存入Response中转
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
