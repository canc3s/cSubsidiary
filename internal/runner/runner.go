package runner

import (
	"fmt"
	"github.com/canc3s/cSubsidiary/internal/fileutil"
	"github.com/canc3s/cSubsidiary/internal/gologger"
	"os"
	"regexp"
)

type Targets struct {
	ID			[]string
	Name		[]string
}

func RunEnumeration(options *Options) {
	if options.InputFile != "" {
		fin, error := os.OpenFile(options.InputFile, os.O_RDONLY, 0)
		if error != nil {
			gologger.Fatalf("文件读取失败：%s",error)
		}
		defer fin.Close()
		imf := fileutil.ReadImf(fin)
		targets := TransImf(imf)
		for _,id := range targets.ID {
			options.CompanyID = id
			SubsidiariesById(options)
		}
		for _,name := range targets.Name {
			options.CompanyName = name
			GetCompanyID(options)
			if options.CompanyID != "" {
				SubsidiariesById(options)
			}
		}
	}else{
		if options.CompanyName != "" && options.CompanyID == "" {
			GetCompanyID(options)
		}
		if options.CompanyID != "" {
			SubsidiariesById(options)
		}
	}
}

func SubsidiariesById(options *Options) {
	var subsidiaries []Subsidiary
	gologger.Infof("正在查询 https://www.tianyancha.com/company/%s 的子公司\n", options.CompanyID)
	subsidiaries = ListInvest(options)

	for _,subsidiary := range subsidiaries {
		fmt.Printf("%s\t%s\n", subsidiary.Url, subsidiary.Name)
		//fmt.Printf("%s\t%s\t%s\n", subsidiary.Url, subsidiary.Name, subsidiary.Percent)
	}
	if options.Output != "" {
		file, err := os.OpenFile(options.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			gologger.Fatalf("结果无法写入文件：\n%s\n", err)
		}
		defer file.Close()

		for _,subsidiary := range subsidiaries {
			file.WriteString(subsidiary.Url+"\t"+subsidiary.Name+"\n")
		}
	}

}

func TransImf(imf []string) Targets {
	var targets Targets
	for _,i := range imf {
		re := regexp.MustCompile(`(\d{6,11})`)
		buf := re.FindStringSubmatch(i)
		if buf == nil {
			targets.Name = append(targets.Name, i)
		}else{
			targets.ID = append(targets.ID, buf[0])
		}
	}
	return targets
}