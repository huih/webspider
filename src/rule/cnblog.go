package rule

// 基础包
import (
//	"github.com/PuerkitoBio/goquery"                        //DOM解析
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
	// "regexp"
//	"strconv"
//	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	cnblogs.Register()
}

var cnblogs = &Spider{
	Name:        "cnblogs",
	Description: "cnblogs博文采集 [http://www.cnblogs.com/]",
	Host: "http://www.cnblogs.com",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:      USE,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&context.Request{
							Url:   ctx.Spider.GetHostUrl(),
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					itemBody := query.Find("div.post_item_body")
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}
					
					//get detail page url
					for index := 0; index < itemBody.Length(); index++ {
						div := itemBody.Eq(index)
						hrefa := div.Find("a[class]")
						href,_ := hrefa.Attr("href")
						ctx.AddQueue(&context.Request{
							Url:    href,
							Rule:   "详细页面",
						})
					}
					
//					//get next page url
//					allPage := query.Find("#pager_bottom #paging_block .pager")
//					pageList := allPage.Children()
//					if pageList.Length() <= 0 {
//						return
//					}
//					nextPage := pageList.Eq(pageList.Length() - 1)
//					text,_ := nextPage.Attr("href")
					
//					if text != "" {
//						isLast := strings.Contains(ctx.GetUrl(), text)
//						if isLast == false {
//							nextPageUrl := ctx.Spider.GetHostUrl() + text
//							logs.Debug("spiderName: %s", ctx.Spider.GetName())
//							logs.Debug("nextPageUrl: %s", nextPageUrl)
//							ctx.AddQueue(&context.Request{
//							Url:    nextPageUrl,
//							Rule:   "生成请求",
//							})
//						}
//					}
				},
			},

			"详细页面": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"网站名称",
					"网页内容",
					"链接",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					
					bodyContent := query.Find("#cnblogs_post_body")
					bodyHtml, err := bodyContent.Html()
					if err != nil || bodyHtml == "" {
						logs.Debug("get body content html err: %s", err.Error())
						return;
					}
					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: ctx.Spider.GetName(),
						1: bodyHtml,
						2: ctx.GetUrl(),
					})
				},
			},
		},
	},
}
