package dd_contacts

import(
	dt "github.com/Peanuttown/tzzGoUtil/datastruct"
	daoFac  "github.com/Peanuttown/dd_contacts/dao/factory"
	"context"
)

type DeptTreeByDB struct{
	tree *dt.Tree
}

type DeptNodeValue struct{
	DeptId uint
	Name string
}

type deptTreeBuilderByDB struct{
	deptId uint
	daoFac daoFac.DaoFactoryI
}

func newDeptTreeBuilderByDB(deptId uint,daoF daoFac.DaoFactoryI) *deptTreeBuilderByDB{
	return &deptTreeBuilderByDB{
		deptId:deptId,
		daoFac:daoF,
	}
}

func (this *deptTreeBuilderByDB) GetValue(ctx context.Context)(interface{},error){
	dept,err := this.daoFac.NewDaoDept().FindDept(ctx,this.deptId)
	if err != nil{
		return nil,err
	}
	return &DeptNodeValue{
		DeptId:dept.ID,
		Name:dept.Name,
	},nil
}

func(this *deptTreeBuilderByDB) GetChildren(ctx context.Context)([]dt.TreeBuilderI,error){
	subDeptIds,err := this.daoFac.NewDaoDept().FindSubDeptIds(ctx,this.deptId)
	if err != nil{
		return nil,err
	}
	treeBuilders := make([]dt.TreeBuilderI,0,len(subDeptIds,))
	for _,id := range subDeptIds{
		treeBuilders = append(treeBuilders,newDeptTreeBuilderByDB(id,this.daoFac))
	}
	return treeBuilders,nil
}


func BuildDeptTreeByDB(ctx context.Context,daoF daoFac.DaoFactoryI,deptId uint)(*DeptTreeByDB,error){
	tree,err := dt.BuildTree(ctx,newDeptTreeBuilderByDB(deptId,daoF))
	if err != nil{
		return nil,err
	}
	return &DeptTreeByDB{
		tree:tree,
	},nil

}
