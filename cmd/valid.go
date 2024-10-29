package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"proxyCrawler/model"
	"proxyCrawler/module/validator"
	"strconv"
	"strings"
	"time"
)

var validCmd = &cobra.Command{
	Use:   "valid",
	Short: "验证代理,目前支持 http和socks",
	Long:  os.Args[0] + " valid http://user:pass@localhost:8088",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		valid(args[0])
	},
}

func valid(u string) {
	parse, err := url.Parse(u)
	if err != nil {
		fmt.Println("URL不合法: ", err)
		os.Exit(0)
	}
	host := strings.Split(parse.Host, ":")[0]
	portStr := parse.Port()
	port, _ := strconv.Atoi(portStr)
	password, _ := parse.User.Password()
	proxy := model.Proxy{
		CType:             "",
		Protocol:          parse.Scheme,
		Host:              host,
		Port:              port,
		Validated:         false,
		Latency:           0,
		ValidateDate:      time.Time{},
		ToValidateDate:    time.Time{},
		ValidateFailedCnt: 0,
		User:              parse.User.Username(),
		Pass:              password,
		Link:              u,
		Country:           "",
	}
	vp := validator.ValidProxy(proxy)
	fmt.Printf("%s -> %dms", vp.Link, vp.Latency)
}
