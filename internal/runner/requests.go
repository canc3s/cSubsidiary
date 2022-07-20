package runner

import (
	"crypto/tls"
	"encoding/json"
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

type InvestList struct {
	State	string `json:"state"`
	ErrorCode	int `json:"errorCode"`
	Data	Data `json:"data"`
}

type Data struct {
	Results	[]Result `json:"result"`
}

type Result struct {
	Name		string `json:"name"`
	Id			int64 `json:"id"`
	Amount		string `json:"amount"`
	RegStatus	string `json:"regStatus"`
	Percent		string `json:"percent"`
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

func ListInvest(options *Options) []Subsidiary {
	var transport = DefaultTransport()
	var client = &http.Client{
		Transport: transport,
		//Timeout:       time.Duration(options.Timeout),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	body := strings.NewReader("{\"gid\":\""+options.CompanyID+"\",\"pageSize\":200,\"pageNum\":1}")
	req, _ := http.NewRequest("POST", "https://capi.tianyancha.com/cloud-company-background/company/investListV2", body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("Content-Type", "application/json")
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
	resbody, _ := ioutil.ReadAll(resp.Body)
	var invests InvestList
	if err := json.Unmarshal(resbody, &invests); err != nil{
		gologger.Fatalf("查询失败：",err.Error(),"\n")
	}
	var subsidiaries []Subsidiary
	re := regexp.MustCompile(`(\d*)`)
	if invests.State == "ok" && invests.ErrorCode == 0 {
		for _,i := range invests.Data.Results {
			var subsidiary = Subsidiary{
				Name: i.Name,
				Url: "https://www.tianyancha.com/company/"+strconv.FormatInt(i.Id,10),
				Percent: i.Percent,
				Status:  i.RegStatus != "注销",
				Funds: re.FindStringSubmatch(i.Amount)[0],
			}
			percent1, _ := strconv.ParseFloat(subsidiary.Percent[:len(subsidiary.Percent)-1], 64)
			funds1, _ := strconv.Atoi(subsidiary.Funds)
			if percent1 >= float64(options.Percent) && funds1 >= options.Funds && subsidiary.Status {
				subsidiaries = append(subsidiaries, subsidiary)
			}
		}
	}
	return subsidiaries
}

func GetCompanyID(options *Options) {
	CompanyName := url.QueryEscape(options.CompanyName)
	resp := GetPage("https://sp0.tianyancha.com/search/suggestV2.json?key="+CompanyName, options)
	re := regexp.MustCompile(`,"id":(\d*),"comName":`)
	options.CompanyID = string(re.FindSubmatch(resp.Body)[1])
}
