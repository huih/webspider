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
//	"output"
)

func init() {
//	linuxnews.Register()
}

func linuxnewsFormatContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	csblogCode := bodyContent.Find("div.linuxnews_code")
	if csblogCode != nil {
		for index := 0; index < csblogCode.Length(); index++ {
			dom := csblogCode.Eq(index)
			dom.RemoveClass("linuxnews_code")
		}
	}
	
	preDom := bodyContent.Find("pre")
	if preDom != nil {
		for index := 0; index < preDom.Length(); index++{
			dom := preDom.Eq(index)
			dom.AddClass("prettyprint linenums")
		}
	}
	
	tableDom := bodyContent.Find("table")
	if tableDom.Length() > 0 {
		for index := 0; index < tableDom.Length(); index++ {
			dom := tableDom.Eq(index)
			dom = dom.SetAttr("style", "width:100%;border:solid thin black;")
			dom = dom.SetAttr("cellpadding", "2")
			dom = dom.SetAttr("cellspacing", "2")
		}
	}
	
	//remove the spicial dom
	for _, dom := range removeDom {
		tmpDom := bodyContent.Find(dom)
		if tmpDom != nil {
			tmpDom.Remove()
		}
	}
	return bodyContent
}

var linuxnews = &Spider{
	Name:        "linux news",
	Description: "linux news博文采集 [https://linux.cn/news/]",
	Host: "https://linux.cn/news/",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	StopSpider: false,
	Keyword:      USE,
	EnableCookie: false,
	MaxItemNum: RULE_SPIDER_PAPER_NUM,
	RemoveDom:[]string{},
	FilterDom: []string{},
	FilterUrl: []string{},
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&context.Request{
							Url:   "http://www.linuxnews.com",
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					
					tmpquery, _ := query.Html()
					logs.Debug(tmpquery)
					
					itemBody := query.Find("ul.article-list")
	
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}
					
					//get detail page url
					for index := 0; index < itemBody.Length() && ctx.Spider.ContinueGetItem(); index++ {
						div := itemBody.Eq(index)
						title := div.Find("span.title a.new")
						href,_ := title.Attr("href")
						
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
						
						logs.Debug(href)
//						ctx.AddQueue(&context.Request{
//							Url:    href,
//							Rule:   "详细页面",
//						})
					}
					
					//get next page url
					if (ctx.Spider.GetStopSpider()) {
						logs.Debug("stop spider: %s", ctx.Spider.GetHostUrl())
						return
					}
					
					//get next page url
					allPage := query.Find("#pager_bottom #paging_block .pager")
					pageList := allPage.Children()
					if pageList.Length() <= 0 {
						return
					}
					nextPage := pageList.Eq(pageList.Length() - 1)
					text,_ := nextPage.Attr("href")
					
					if text != "" {
						isLast := strings.Contains(ctx.GetUrl(), text)
						if (isLast == false) && ctx.Spider.ContinueGetItem() {
							nextPageUrl := ctx.Spider.GetHostUrl() + text
							//output.SaveDataToLocalDB("category", "category", ctx.Spider.GetName(), ctx.GetUrl())
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
//					titleDom := query.Find("#cb_post_title_url");
//					paperTitle := titleDom.Text()
					
					//find content
					bodyContent := query.Find("#linuxnews_post_body")
					
					//filter bad paper
					for _, dom := range ctx.Spider.GetFilterDom() {
						tmpDom := bodyContent.Find(dom)
						if tmpDom.Length() > 1 {
							return
						}
					}
					
					bodyContent = linuxnewsFormatContent(ctx.Spider.GetRemoveDom(), bodyContent)
					
					bodyHtml, err := bodyContent.Html()
					if err != nil || bodyHtml == "" {
						logs.Debug("get body content html err: %s", err.Error())
						return;
					}
					
						//remove all code linenum
					reg := regexp.MustCompile(`<span style=\"color: #008080;?\">\s*(&nbsp;)?(\d+)<\/span>`)
					bodyHtml = reg.ReplaceAllLiteralString(bodyHtml, "")
					
					//save result to local db for replication
//					output.SaveDataToLocalDB(paperTitle, bodyHtml, ctx.Spider.GetName(), ctx.GetUrl())
					
					// 结果存入Response中转
//					ctx.Output(map[int]interface{}{
//						0: paperTitle,
//						1: bodyHtml,
//						2: ctx.Spider.GetName(),
//						3: ctx.GetUrl(),
//					})
				},
			},
		},
	},
}
