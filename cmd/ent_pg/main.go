package main

import (
	"context"
	"log"
	"github.com/Peanuttown/dd_contacts/ent"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client,err := ent.DialMysql("root:tzzjkl@tcp(127.0.0.1:3306)/ent_pg")
	if err != nil{
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Schema.Create(ctx)
	if err != nil{
		log.Fatal(err)
	}
}
