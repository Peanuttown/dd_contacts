package schema

import(
	"errors"
	"github.com/Peanuttown/dd_contacts/ent"
	up "github.com/Peanuttown/dd_contacts/ent/userpropertyindept"
	_ "github.com/go-sql-driver/mysql"
	"github.com/Peanuttown/dd_contacts/ent/dept"
	ent_user "github.com/Peanuttown/dd_contacts/ent/user"
	_"github.com/Peanuttown/dd_contacts/ent/runtime"
	"fmt"
	"context"
	"testing"
)

func testClient()*ent.Client{
	client,err := ent.Open("mysql","root:tzzjkl@tcp(127.0.0.1:3306)/ent_pg")
	if err != nil{
		panic(err)
	}
	err =client.Schema.Create(context.Background())
	if err != nil{
		panic(err)
	}
	return client

}

func TestDept(t *testing.T){
	cli := testClient()
	var testDeptId uint = 123
	var testDeptId2 uint =1234
	var testUserId = "testUserId"
	var testUserId2 = "testUserId2"
	var testUserId3 = "testUserId3"
	ctx := context.Background()
	// < ready data
	// << clean data first
	// ignore err
	var err error
	tx,err := cli.Tx(ctx)
	if err != nil{
		t.Fatal(err)
	}
	 err =tx.Dept.DeleteOneID(testDeptId).Exec(ctx)
	 errNotFound := &ent.NotFoundError{}
	 if err != nil && !errors.As(err,&errNotFound){
		 t.Fatal(err)
	 }
	 err = tx.Dept.DeleteOneID(testDeptId2).Exec(ctx)
	 if err != nil && !errors.As(err,&errNotFound){
		 t.Fatal(err)
	 }
	tx.User.DeleteOneID(testUserId).Exec(ctx)
	 if err != nil && !errors.As(err,&errNotFound){
		 t.Fatal(err)
	 }
	tx.User.DeleteOneID(testUserId2).Exec(ctx)
	 if err != nil && !errors.As(err,&errNotFound){
		 t.Fatal(err)
	 }
	tx.Commit()
	tx.User.DeleteOneID(testUserId3).Exec(ctx)
	 if err != nil && !errors.As(err,&errNotFound){
		 t.Fatal(err)
	 }
	// >>
	// >

	// < insert test data
	testDept ,err := cli.Dept.Create().SetID(testDeptId).SetName("testName").Save(ctx)
	if err != nil{
		t.Fatal(err)
	}
	testDept2 ,err := cli.Dept.Create().SetID(testDeptId2).SetName("testName").Save(ctx)
	if err != nil{
		t.Fatal(err)
	}
	// << add user to the dept
	user,err := cli.User.Create().SetID(testUserId).SetName("testName").SetPhone("testPhone").AddDepts(testDept).Save(ctx)
	if err != nil{
		t.Fatal(err)
	}
	user2,err := cli.User.Create().SetID(testUserId2).SetName("testName").SetPhone("testPhone").AddDepts(testDept,testDept2).Save(ctx)
	if err != nil{
		t.Fatal(err)
	}
	user3,err := cli.User.Create().SetID(testUserId3).SetName("testName").SetPhone("testPhone").AddDepts(testDept,testDept2).Save(ctx)
	if err != nil{
		t.Fatal(err)
	}
	user2PropertiesInDept,err := cli.UserPropertyInDept.Create().SetDept(testDept).SetUser(user2).SetIsLeader(true).Save(ctx)
	if err !=nil{
		t.Fatal(err)
	}
	user3PropertiesInDept,err := cli.UserPropertyInDept.Create().SetDept(testDept).SetUser(user3).SetIsLeader(true).Save(ctx)
	if err !=nil{
		t.Fatal(err)
	}
	// >>
	// >

	// < delete dept, expect the user also deleted
	tx,err = cli.Tx(ctx)
	if err != nil{
		t.Fatal(err)
	}
	err = tx.Dept.DeleteOne(testDept).Exec(ctx)
	if err != nil{
		t.Fatal(err)
	}
	err = tx.Commit()
	if err != nil{
		t.Fatal(err)
	}
	// << query dept weather deleted
	count,err := cli.Dept.Query().Where(dept.IDEQ(testDept.ID)).Count(ctx)
	if err != nil{
		t.Fatal(err)
	}
	if count > 0 {
		t.Fatal(fmt.Errorf("dept deleted failed"))
	}
	// >>
	count,err =cli.UserPropertyInDept.Query().Where(up.IDEQ(user3PropertiesInDept.ID)).Count(ctx)
	if err != nil{
		t.Fatal(err)
	}
	if count != 0{
		t.Fatal(fmt.Errorf("when dept deleted , the userPropertyInDepts still exists"))
	}
	// << query user weather deleted
	count,err =cli.User.Query().Where(ent_user.IDEQ(user.ID)).Count(ctx)
	if err != nil{
		t.Fatal(err)
	}
	if count > 0{
		t.Fatal("when the dept delted , the user in the dept still exists")
	}
	count,err =cli.User.Query().Where(ent_user.IDEQ(user2.ID)).Count(ctx)
	if err != nil{
		t.Fatal(err)
	}
	if count != 1{
		t.Fatal("the user in two depts also deleted")
	}
	tx,err = cli.Tx(ctx)
	if err != nil{
		t.Fatal(err)
	}
	err = tx.User.DeleteOne(user2).Exec(ctx)
	if err != nil{
		t.Fatal(err)
	}
	err =tx.Commit()
	if err != nil{
		t.Fatal(err)
	}
	count,err = cli.UserPropertyInDept.Query().Where(up.IDEQ(user2PropertiesInDept.ID)).Count(ctx)
	if err != nil{
		t.Fatal(err)
	}
	if count != 0{
		t.Fatal("when user deleted , userPropertiesInDept still exists")
	}
	// >>

	// >

}
