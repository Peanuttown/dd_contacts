package ent

import(
	"fmt"
	"context"
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

func (this *ClientWrapper) IsTx() bool {
	_,ok := this.driver.(*txDriver)
	return ok
}

type TxDoFunc = func(ctx context.Context,tx *ClientWrapper)error

func (this *ClientWrapper) TxDoIfClientNotTx(ctx context.Context,f TxDoFunc)error{
	if this.IsTx(){
		return f(ctx,this)
	}
	return this.TxDo(ctx,f)
}

func (this *ClientWrapper) TxDo(ctx context.Context,f TxDoFunc) (err error){
	tx,err := this.Tx(ctx)
	if err != nil{
		return fmt.Errorf("Start tx failed: %w",err)
	}
	defer func(){
		e := recover()
		if e != nil{
			err = fmt.Errorf("%v",e)
			tx.Rollback()
			return
		}
		if err != nil{
			tx.Rollback()
			return 
		}
		err = tx.Commit()
	}()
	err = f(ctx,NewClientWrapper(tx.Client()))
	return 
}
