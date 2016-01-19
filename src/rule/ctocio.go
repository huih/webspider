package rule

// 基础包
import (
	"strconv"
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/gotools/logs"                   //信息输出

	// 字符串处理包
	 "regexp"
	"strings"
	"output"
)
//http://[\S]*.(png|jpg|jpeg|gif)
var deleteContent = []string {
	`<p><strong>第一时间获取面向IT决策者的独家深度资讯，敬请关注IT经理网微信号：<span style=\"color: #4bacc6;\">ctociocom</span></strong></p>`,
	`<p>[\r\n]*<a href=\"http://www.ctocio.com.*qrcode_for_gh_.*.jpg\">[\r\n]*<img class=\"alignnone size-full wp-image-19931\" src=\"http://www.ctocio.com.*qrcode_for_gh_.*.jpg\" alt=\"qrcode_for_gh_.*>[\r\n]*</a>[\r\n]*</p>`,
	`<span style=\"color: #008080;?\">\s*(&nbsp;)?(\d+)<\/span>`,
	`<!--[\w\W\r\n]*?-->`,
	`<p>[&nbsp;]*</p>`,
	`(\n\r)*`,
}

func init() {
	ctocio.Register()
}


func formatContent(removeDom []string, bodyContent *goquery.Selection) *goquery.Selection{
	csblogCode := bodyContent.Find("div.cnblogs_code")
	if csblogCode != nil {
		for index := 0; index < csblogCode.Length(); index++ {
			dom := csblogCode.Eq(index)
			dom.RemoveClass("cnblogs_code")
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
		for index :=0; index < tmpDom.Length(); index++ {
			tmp := tmpDom.Eq(index)
			tmp.Remove()
		}
	}

	return bodyContent
}

var ctocio = &Spider{
	Name:        "ctocio",
	Description: "ctocio博文采集 [http://www.ctocio.com/category/ccnews]",
	Host: "http://www.ctocio.com",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:      USE,
	StopSpider:   false,
	EnableCookie: false,
	MaxItemNum: RULE_SPIDER_PAPER_NUM,
	RemoveDom:[]string{"div.wp-pagenavi", "span.pages"},
	//FilterDom: []string{"#codeSnippetWrapper","meta","pre.best-text mb-10"},
	//FilterUrl: []string{"depsi", "liwei225"},
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&context.Request{
							Url:   "http://www.ctocio.com/category/ccnews",
							Rule:   aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					itemBody := query.Find("div.post")
	
					if itemBody.Length() <= 0 {
						logs.Debug("finished %s", ctx.GetUrl())
						return;
					}
					
					//get detail page url
					for index := 0; index < itemBody.Length() && ctx.Spider.ContinueGetItem(); index++ {
						div := itemBody.Eq(index)
						titleDom := div.Find("h2")
						hrefa := titleDom.Find("a[href]")
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
					allPage := query.Find("div.wp-pagenavi span.pages")
					if allPage.Length() <= 0 {
						return
					}
					
					pageText := allPage.Text()
					reg := regexp.MustCompile(`[\d]+`)
					numstrs := reg.FindAllString(pageText, -1)
					if len(numstrs) < 2 {
						return
					}
					curPage,_ := strconv.Atoi(numstrs[0])
					totalPage,_ := strconv.Atoi(numstrs[1])
					if curPage > 0 && curPage < totalPage && ctx.Spider.ContinueGetItem() {
						nextPageUrl := "http://www.ctocio.com/category/ccnews/page/" + strconv.Itoa((curPage + 1)) 
						ctx.AddQueue(&context.Request{
						Url:    nextPageUrl,
						Rule:   "生成请求",
						})
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
					
					paperBody := query.Find("div.post")
					
					//find title
					titleDom := paperBody.Find("h1")
					paperTitle := titleDom.Text()
					
					//find content
					bodyContent := paperBody.Find("div.entrys")
					
					//filter bad paper
					for _, dom := range ctx.Spider.GetFilterDom() {
						tmpDom := bodyContent.Find(dom)
						if tmpDom.Length() > 1 {
							return
						}
					}
					
					bodyContent = formatContent(ctx.Spider.GetRemoveDom(), bodyContent)
					
					bodyHtml, err := bodyContent.Html()
					if err != nil || bodyHtml == "" {
						logs.Debug("get body content html err: %s", err.Error())
						return;
					}
					
					//remove bad content
					for _, item := range deleteContent {
						reg := regexp.MustCompile(item)
						bodyHtml = reg.ReplaceAllLiteralString(bodyHtml,"")
					}
					
					//save result to local db for replication
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
