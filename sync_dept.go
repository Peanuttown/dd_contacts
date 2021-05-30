package dd_contacts

import(
	"time"
	"fmt"
	"context"
	"github.com/Peanuttown/dd_contacts/ent"
	daoFac "github.com/Peanuttown/dd_contacts/dao/factory"
	"github.com/Peanuttown/dd_contacts/dao/models"
	"github.com/Peanuttown/dd_api"
  dt "github.com/Peanuttown/tzzGoUtil/datastruct"
)

func SyncDept(ctx context.Context,dbClient *ent.ClientWrapper,ddApiClient *dd_api.Client)(error){
	// < sync_dept
	// steps: 
	// << 
	//   1. get all depts from db, compared with deptTree build by ddApi,
	//   2.  add, and update dept to db,  
	//   3. then sync user
	//   4. deleted dept at last
	// >>
	var generation = uint(time.Now().Unix())
	return daoFac.NewDaoFactoryTx(dbClient).TxDo(
		ctx,
		func (ctx context.Context,daoF daoFac.DaoFactoryI)error{
			deptTreeBuildByApi,err := dd_api.BuildDeptTreeByApi(ctx,dd_api.ROOT_DEPT_ID,ddApiClient)
			if err != nil{
				return err
			}
			fmt.Println(deptTreeBuildByApi.ToString(ctx))
			// TODO GOON
			err = deptTreeBuildByApi.DepthFirstDo(
				ctx,
				func(
					ctx context.Context,
					node *dt.Node,
				)error{
					dept := node.GetValue().(*dd_api.DeptNodeValue)
					parentDeptNode := node.GetParent()
					var parentId  uint = 0
					if parentDeptNode != nil{
						parentId = (parentDeptNode.GetValue().(*dd_api.DeptNodeValue)).DeptId
					}
					// upsert to db
					err = daoF.NewDaoDept().Upsert(
						ctx,
						models.NewDeptRequiredFields(
							dept.DeptId,
							dept.Name,
						),
						models.DeptOptionalParentDeptId(parentId),
						models.DeptOptionalGeneration(generation),
					)
					if err != nil{
						return err
					}
					// update uesr in the dept
					userIdsInTheDepts,err := dd_api.NewApiDeptGetUserIds(dept.DeptId).ExecBy(ctx,ddApiClient)
					fmt.Println("users: ",userIdsInTheDepts)
					if err != nil{
						return err
					}
					// query user info
					for _,userId:= range userIdsInTheDepts{
						userInfo,err := dd_api.NewApiUserGetDetail(userId).ExecBy(ctx,ddApiClient)
						if err != nil{
							return err
						}
						pts := make([]models.UserPropertiesInDepts,0,len(userInfo.IsLeaderInDepts))
						for _,p := range userInfo.IsLeaderInDepts{
							pts = append(pts,models.UserPropertiesInDepts{
								DeptId:p.DeptId,
								IsDeptLeader:p.Leader,
							})
						}
						fmt.Println("userId is ",userId)
						err = daoF.NewDaoUser().Upsert(
							ctx,
							models.NewUserRequiredFields(
								userId,
								userInfo.Name,
								pts,
								userInfo.Mobile,
							),
						models.UserOptionlGeneration(generation),
						)
						if err != nil{
							return err
						}
					}
					return nil
				},
			)
			if err != nil{
				return err
			}
			fmt.Println("start clear ")
			// clean uesr and dept
			deptsToRemove,err := daoF.NewDaoDept().FindByNotGeneration(ctx,generation)
			if err != nil{
				return fmt.Errorf("Find Dept ByNotGeneration failed: %w",err)
			}
			for _,deptId := range deptsToRemove{
				err = daoF.NewDaoDept().Delete(ctx,deptId)
				if err != nil{
					return err
				}
			}
			usersToRemove,err := daoF.NewDaoUser().FindByNotGeneration(ctx,generation)
			if err != nil{
				return fmt.Errorf("Find User ByNotGeneration failed: %w",err)
			}
			for _,user := range usersToRemove{
				err := daoF.NewDaoUser().DeleteUser(ctx,user)
				if err != nil{
					return fmt.Errorf("delete user by id %v failed : %w",user,err)
				}
			}
			return nil
		},
	)
	// >
}

