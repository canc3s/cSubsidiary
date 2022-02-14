package runner

import (
	"crypto/tls"
	"github.com/antchfx/htmlquery"
	"github.com/canc3s/cSubsidiary/internal/gologger"
	"golang.org/x/net/html"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Subsidiary struct {
	Name    string
	Url     string
	Percent string
	Funds   string
	Status  bool
}

type Request struct {
	Url    string
	Cookie string
}

type Response struct {
	Body []byte
	Page *html.Node
}

func DefaultTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}
	return transport
}

func GetPage(url string, options *Options) Response {

	var transport = DefaultTransport()
	var client = &http.Client{
		Transport: transport,
		//Timeout:       time.Duration(options.Timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:78.0) Gecko/20100101 Firefox/78.0")
	if options.Cookie != "" {
		req.Header.Set("Cookie", options.Cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		gologger.Fatalf("请求发生错误，请检查网络连接\n%s\n", err)
	}

	if resp.StatusCode == 403 {
		gologger.Fatalf("海外用户或者云服务器ip被禁止访问网站，请更换ip\n")
	} else if resp.StatusCode == 401 {
		gologger.Fatalf("天眼查Cookie有问题或过期，请重新获取\n")
	} else if resp.StatusCode == 302 {
		gologger.Fatalf("天眼查免费查询次数已用光，需要加Cookie\n")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	page, _ := htmlquery.Parse(strings.NewReader(string(body)))

	return Response{
		Body: body,
		Page: page,
	}
}

func GetInformation(resp Response, options *Options) []Subsidiary {
	list := htmlquery.Find(resp.Page, "//*[@id=\"_container_invest\"]/div/table/tbody/tr")

	subsidiaries := filter(list, options)
	return subsidiaries
}

func GetInformationWithCookie(page *html.Node, options *Options) []Subsidiary {
	list := htmlquery.Find(page, "/html/body/div/table/tbody/tr")

	subsidiaries := filter(list, options)
	return subsidiaries
}

func filter(list []*html.Node, options *Options) (subsidiaries []Subsidiary) {
	for _, n := range list {
		nodes := htmlquery.Find(n, "//td")
		nodeA := htmlquery.FindOne(n, "//td/div/a[1]")

		if len(nodes) < 11 {
			continue
		}
		re := regexp.MustCompile(`(\d*)`)
		state := htmlquery.InnerText(nodes[10])
		funds := htmlquery.InnerText(nodes[8])
		var subsidiary = Subsidiary{
			Name:    strings.Trim(htmlquery.InnerText(nodes[3]), "股权结构"),
			Url:     htmlquery.SelectAttr(nodeA, "href"),
			Funds:   re.FindStringSubmatch(funds)[0],
			Percent: htmlquery.InnerText(nodes[9]),
			Status:  state == "存续（在营、开业、在册）",
		}

		percent1, _ := strconv.ParseFloat(subsidiary.Percent[:len(subsidiary.Percent)-1], 64)
		funds1, _ := strconv.Atoi(subsidiary.Funds)
		if percent1 >= float64(options.Percent) && funds1 >= options.Funds && subsidiary.Status {
			subsidiaries = append(subsidiaries, subsidiary)
		}
	}
	return subsidiaries
}

func JudgePages(page *html.Node) int {
	list := htmlquery.Find(page, "/html/body/div[2]/div/div/div[5]/div[1]/div/div[3]/div[1]/div[7]/div[2]/div/div/ul/li/a")
	return len(list)
}

func JudgePagesI(page *html.Node) int {
	list := htmlquery.Find(page, "/html/body/div/div/ul/li/a")
	return len(list) - 1
}

func GetCompanyID(options *Options) {
	CompanyName := url.QueryEscape(options.CompanyName)
	resp := GetPage("https://sp0.tianyancha.com/search/suggestV2.json?key="+CompanyName, options)
	re := regexp.MustCompile(`,"id":(\d*),"comName":`)
	options.CompanyID = string(re.FindSubmatch(resp.Body)[1])
}
