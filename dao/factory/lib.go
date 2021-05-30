package fac

import(
	"context"
	"github.com/Peanuttown/dd_contacts/dao"
	"github.com/Peanuttown/dd_contacts/dao/impl"
	"github.com/Peanuttown/dd_contacts/ent"
)



type DaoFactoryI interface{
	NewDaoDept() dao.DaoDeptI
	NewDaoUser() dao.DaoUserI
}

type daoFactory struct{
	client *ent.ClientWrapper
}

func newDaoFactory(client *ent.ClientWrapper)DaoFactoryI{
	return &daoFactory{
		client:client,
	}
}

func (this *daoFactory) NewDaoDept() dao.DaoDeptI{
	return impl.NewDaoDept(this.client)
}

func (this *daoFactory) NewDaoUser() dao.DaoUserI{
	return impl.NewDaoUser(this.client)
}

type DaoFactoryTxI interface{
	TxDo(ctx context.Context,f func(ctx context.Context,daoF DaoFactoryI)error)error
}

type daoFactoryTx struct{
	client *ent.ClientWrapper
}

func NewDaoFactoryTx(client *ent.ClientWrapper)DaoFactoryTxI{
	return &daoFactoryTx{
		client:client,
	}
}

func (this *daoFactoryTx) TxDo(
	ctx context.Context,
	f func(ctx context.Context,daoF DaoFactoryI)error,
)error{
	return this.client.TxDo(
		ctx,
		func (ctx context.Context,tx *ent.ClientWrapper)error{
			daoF := newDaoFactory(tx)
			return f(ctx,daoF)
		},
	)
	
}
