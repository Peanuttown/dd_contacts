package ent

import(
	"fmt"
)


type ClientWrapper struct{
	*Client
}

func NewClientWrapper(client *Client)*ClientWrapper{
	return &ClientWrapper{
		Client:client,
	}

}

func DialMysql(mysqlUrl string)(*ClientWrapper,error){
	client,err := Open("mysql",mysqlUrl)
	if err != nil{
		return nil,fmt.Errorf("Connect mysql <%v> error : %w",mysqlUrl,err)
	}
	return NewClientWrapper(client),nil
}
