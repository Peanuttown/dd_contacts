package main

import (
	"os"
	"context"
	"fmt"
	"flag"
	"github.com/Peanuttown/dd_api"
	"github.com/Peanuttown/dd_contacts/ent"
	"github.com/Peanuttown/dd_contacts"
	"github.com/Peanuttown/tzzGoUtil/log"
)



// oapi.dingtalk.com
func main() {
	logger := log.NewLogger()
	var ddAppId string
	var ddAppSecret string
	var ddApiHost string
	var sqlSource string
	flag.StringVar(&ddAppId,"ddAppId","","DingDing app id")
	flag.StringVar(&ddAppSecret,"ddAppSecret","","DingDing app secret")
	flag.StringVar(&ddApiHost,"ddApiHost","","DingDing api host")
	flag.StringVar(&sqlSource,"sqlSource","","Mysql source")
	flag.Parse()
	// < check params
	if len(ddAppId) == 0{
		panic(fmt.Errorf("input error, ddAppId is nil"))
	}
	if len(ddAppSecret) == 0 {
		panic(fmt.Errorf("input error, ddAppSecret is nil"))
	}
	if len(sqlSource) == 0{
		panic(fmt.Errorf("sqlSource is nil"))
	}
	if len(ddApiHost) == 0{
		panic(fmt.Errorf("DingDing api host is nil"))
	}
	// >

	logger.Debug("Connecting to db")
	dbCli,err := ent.DialMysql(sqlSource)
	if err != nil{
		panic(fmt.Errorf("Connect to db failed: %w",err))
	}
	ctx := context.Background()
	// < migrate table
	err = dbCli.Schema.Create(ctx)
	if err != nil{
		panic(fmt.Errorf("Migrate table failed: %w",err))
	}
	// >
	logger.Debug("Connect db success")
	ddApiCli:=dd_api.NewClient(&dd_api.Cfg{
		AppKey:ddAppId,
		AppSecret:ddAppSecret,
		ApiHost:ddApiHost,
	})
	logger.Debug("Syncing dingding contacts")
	err = dd_contacts.SyncDept(ctx,dbCli,ddApiCli,logger)
	if err != nil{
		logger.Error("❌ Sync dingding contacts failed: %w",err)
		os.Exit(1)
	}
	logger.Info("✅ Sync dingding contacts success")
}
