package runner

import (
	"flag"
	"github.com/canc3s/cSubsidiary/internal/fileutil"
	"github.com/canc3s/cSubsidiary/internal/gologger"
	"os"
)

const banner = `
 ██████╗███████╗██╗   ██╗██████╗ ███████╗██╗██████╗ ██╗ █████╗ ██████╗ ██╗   ██╗
██╔════╝██╔════╝██║   ██║██╔══██╗██╔════╝██║██╔══██╗██║██╔══██╗██╔══██╗╚██╗ ██╔╝
██║     ███████╗██║   ██║██████╔╝███████╗██║██║  ██║██║███████║██████╔╝ ╚████╔╝ 
██║     ╚════██║██║   ██║██╔══██╗╚════██║██║██║  ██║██║██╔══██║██╔══██╗  ╚██╔╝  
╚██████╗███████║╚██████╔╝██████╔╝███████║██║██████╔╝██║██║  ██║██║  ██║   ██║   
 ╚═════╝╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝╚═════╝ ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   
											v`

// Version is the current version of C
const Version = `0.0.4`


type Options struct {
	CompanyName			string
	CompanyID			string                 // Target is a single URL/Domain to scan usng a template
	InputFile			string                 // Targets specifies the targets to scan using templates.
	Cookie				string
	Percent				int
	Funds				int
	Timeout 			int
	Output              string                 // Output is the file to write found subdomains to.
	Silent				bool
	NoColor				bool
	Verbose				bool
	Version             bool                   // Version specifies if we should just show version and exit
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.CompanyName, "n", "", "公司名称")
	flag.StringVar(&options.CompanyID, "i", "", "公司ID号码")
	flag.StringVar(&options.InputFile, "f", "", "包含公司ID号码的文件")
	flag.StringVar(&options.Cookie, "c", "", "天眼查的Cookie")
	flag.IntVar(&options.Percent, "p", 100, "显示投资比例为x以上的公司")
	flag.IntVar(&options.Funds, "w", 100, "显示投资资金为x万以上的公司")
	flag.IntVar(&options.Timeout, "timeout", 15, "连接超时时间")
	flag.StringVar(&options.Output, "o", "", "结果输出的文件(可选)")
	flag.BoolVar(&options.Silent, "silent", false, "Silent mode")
	flag.BoolVar(&options.NoColor, "no-color", false, "No Color")
	flag.BoolVar(&options.Verbose, "verbose", false, "详细模式")
	flag.BoolVar(&options.Version, "version", false, "显示软件版本号")

	flag.Parse()

	options.configureOutput()

	showBanner()

	if options.Version {
		gologger.Infof("Current Version: %s\n", Version)
		os.Exit(0)
	}

	options.validateOptions()

	return options
}

func (options *Options) validateOptions() {
	if options.CompanyID != "" && len(options.CompanyID) < 5 {
		gologger.Fatalf("公司ID %s 不正确!\n", options.CompanyID)
	}
	if options.InputFile != "" && !fileutil.FileExists(options.InputFile) {
		gologger.Fatalf("文件 %s 不存在!\n", options.InputFile)
	}
	if options.CompanyName == "" && options.CompanyID == "" && options.InputFile == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}
}


// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Printf("%s%s\n", banner,Version)
	gologger.Printf("\t\thttps://github.com/canc3s/cSubsidiary\n\n")

	//gologger.Labelf("请谨慎使用,您应对自己的行为负责\n")
	//gologger.Labelf("开发人员不承担任何责任，也不对任何滥用或损坏负责.\n")
}

func (options *Options) configureOutput() {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.MaxLevel = gologger.Verbose
	}
	if options.NoColor {
		gologger.UseColors = false
	}
	if options.Silent {
		gologger.MaxLevel = gologger.Silent
	}
}