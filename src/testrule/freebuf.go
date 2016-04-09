package rule

// 基础包
import (
	"regexp"
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
//	"strconv"
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
//	"output"
"common"
)

var freebufDeleteContent = []string {
	`<p>.*转载请注明来自FreeBuf黑客与极客.*</p>`,
}


func formatFreebufContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	
	preDom := bodyContent.Find("pre")
	if preDom != nil {
		for index := 0; index < preDom.Length(); index++{
			dom := preDom.Eq(index)
			dom.RemoveClass("brush:cpp;toolbar: true; auto-links: false;")
			dom.AddClass("prettyprint linenums")
		}
	}
	
	tableDom := bodyContent.Find("table")
	if tableDom != nil {
		for index := 0; index < tableDom.Length(); index++{
			dom := tableDom.Eq(index)
			dom.RemoveAttr("border")
			dom.SetAttr("border", "1")
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
	
	
	imgDom := bodyContent.Find("img")
	if preDom != nil {
		for index := 0; index < imgDom.Length(); index++{
			dom := imgDom.Eq(index)
			imgorgsrc, ok := dom.Attr("data-original")
			if ok == true {
				dom = dom.RemoveAttr("data-original")
				dom = dom.SetAttr("src", imgorgsrc)	
			}
			
			imgsrc, ok := dom.Attr("src")
			if ok {
				if strings.HasSuffix(imgsrc, "!small") {
					imgsrc = strings.Replace(imgsrc, "!small","",-1)
					dom = dom.SetAttr("src", imgsrc)
				}
			}
		}
	}
	
	return bodyContent
}

func init() {
	freebuf.Register()
}

//http://www.freebuf.com/
var freebuf = &Spider{
	Name:        "freebuf.com",
	Description: "freebuf博文采集 [http://www.freebuf.com/]",
	Host: "http://www.freebuf.com/",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:      USE,
	EnableCookie: false,
	MaxItemNum: RULE_SPIDER_PAPER_NUM,
	StopSpider: false,
	RemoveDom: []string{"noscript"},
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
							Url:   "http://www.freebuf.com/",
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					itemBody := query.Find("div.row div.main-mid #timeline div.news_inner")
	
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}
					
					//get detail page url
					for index := 0; index < itemBody.Length() && ctx.Spider.ContinueGetItem(); index++ {
						div := itemBody.Eq(index)
						hrefa := div.Find("div.news-info dl dt a[href]")
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
//					if (ctx.Spider.GetStopSpider()) {
//						logs.Debug("stop spider: %s", ctx.Spider.GetHostUrl())
//						return
//					}
					
					nextPage := query.Find("div.row div.main-mid div.news-more a[href]")
					if nextPage.Length() <= 0 {
						return;
					}
					
					var text string
					for index := 0; index < nextPage.Length(); index++ {
						pageItem := nextPage.Eq(index)
						text, _ = pageItem.Attr("href")
						break
					}
					
					if text != "" {
						isLast := strings.Contains(ctx.GetUrl(), text)
						if (isLast == false) && ctx.Spider.ContinueGetItem() {
							nextPageUrl := text
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
					path := "div.container div.row div.article-wrap div.articlecontent div.title h2"
					paperTitle, err := common.GetHtmlContent(query, path, ctx, []string{}, false)
					if err != nil {
						logs.Debug("get paper title error")
						return
					}
					
					//find content					
					path = "div.container div.row div.article-wrap div.articlecontent #contenttxt"
					contentObj := query.Find(path)
		
					for _, dom := range ctx.Spider.GetFilterDom() {
						tmpDom := contentObj.Find(dom)
						if tmpDom.Length() > 1 {
							return
						}
					}
		
					contentObj = formatFreebufContent(ctx.Spider.GetRemoveDom(), contentObj)
					
					contentHtml, err := contentObj.Html()
					if err != nil {
						logs.Debug("get html error, err:%s", err.Error())
						return
					}
					
					//remove bad content
					for _, item := range freebufDeleteContent {
						reg := regexp.MustCompile(item)
						contentHtml = reg.ReplaceAllLiteralString(contentHtml,"")
					}
			
					bodyHtml := strings.TrimSpace(contentHtml)
					
//					save result to local db for replication
//					output.SaveDataToLocalDB(paperTitle, bodyHtml, ctx.Spider.GetName(), ctx.GetUrl())
							
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
