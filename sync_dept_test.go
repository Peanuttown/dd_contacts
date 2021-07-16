package dd_contacts

import(
	"context"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"path"
	"os"
	"strings"
		"testing"
	"github.com/Peanuttown/dd_contacts/ent"
  "github.com/Peanuttown/tzzGoUtil/encoding"
  "github.com/Peanuttown/tzzGoUtil/log"
	"github.com/Peanuttown/dd_api"
)

func testDBClient()*ent.ClientWrapper{
	wd,err :=os.Getwd()
	if err != nil{
		panic(err)
	}
	urlBytes,err := ioutil.ReadFile(path.Join(wd,"debug","test_db_url.txt"))
	if err != nil{
		panic(err)
	}
	url := strings.TrimSpace(string(urlBytes))
	fmt.Println(url)
	client,err := ent.DialMysql(url)
	if err != nil{
		panic(err)
	}
	err = client.Schema.Create(context.Background())
	if err != nil{
		panic(err)
	}
	return client
}

func testDDApiClient() *dd_api.Client{
	wd,err :=os.Getwd()
	if err != nil{
		panic(err)
	}
	cfg :=&dd_api.Cfg{}
err =	encoding.UnMarshalByFile(
			path.Join(wd,"debug","test_dingding_cfg.json"),
			cfg,
			json.Unmarshal,
	)
	if err != nil{
		panic(err)
	}
	cli:= dd_api.NewClient(cfg)
	return cli

}

func TestSyncDept(t *testing.T){
	ctx := context.Background()
	logger := log.NewEmptyLogger()
	err := SyncDept(ctx,testDBClient(),testDDApiClient(),logger)
	if err !=nil{
		t.Fatal(err)
	}

}
