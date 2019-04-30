package db_mysql

import (
	"testing"
)

func Init() {
	InitDb("mysql",
		"root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true&timeout=20s&loc=Local&clientFoundRows=true",
		100,
		20, true, true)
}

func Test_ExecNoQuery(t *testing.T) {
	Init()
	//	p := []Product{}
	//	ExecQuery(&p, "select Id,Product_id,Product_name from tbl2 limit 1,100")
	//	for i := range p {
	//		t.Log(p[i])
	//	}

	//err := ExecNoQuery("insert into tbl2(product_id,product_name) values('111','2222')")
	err := ExecNoQuery("insert into tbl2(product_id,product_name) values(?,?)", "abc", "sdfdsfdsf")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("执行成功")
	}

}

func Test_ExecInsertGetLastId(t *testing.T) {
	Init()

	lastId, err := ExecInsertGetLastId("insert into tbl2(product_id,product_name) values('111','2222')")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("执行成功,ID:", lastId)
	}
}

type Product struct {
	Id           int    `db:"id"`
	Product_id   string `db:"product_id"`
	Product_name string `db:"product_name"`
}

func Test_ExecQuery(t *testing.T) {
	Init()

	p := []Product{}
	//p := make([]Product, 0)
	err := ExecQuery(&p, "select id,product_id,product_name from tbl2 limit 1,100")
	if err != nil {
		t.Error(err)
	} else {
		for i := range p {
			t.Log(p[i])
		}
	}
}

func TestQueryByPage(t *testing.T) {
	Init()
	dest := []Product{}
	totalcount, totalpage, err := QueryByPage(&dest, "tbl2", "id,product_id,product_name", "", "", "", 100, 1)
	if err != nil {
		t.Error("totalcount", totalcount, "totalpage:", totalpage, "err:", err)
	} else {
		t.Log("totalcount", totalcount, "totalpage:", totalpage)
		for i := range dest {
			t.Log(dest[i])
		}
	}
}

func Test_ProcedureExecQuery(t *testing.T) {
	Init()
	dest := []Product{}

	params := []PROCEDURE_PARAM{}
	var totalcount int = 0
	var pagecount int = 0
	params = append(params, PROCEDURE_PARAM{Name: "fields", Direct: "in", Value: "id,product_id,product_name"})
	params = append(params, PROCEDURE_PARAM{Name: "table", Direct: "in", Value: "tbl2"})
	params = append(params, PROCEDURE_PARAM{Name: "where", Direct: "in", Value: ""})
	params = append(params, PROCEDURE_PARAM{Name: "orderby", Direct: "in", Value: ""})
	params = append(params, PROCEDURE_PARAM{Name: "pageIndex", Direct: "in", Value: 1})
	params = append(params, PROCEDURE_PARAM{Name: "pageSize", Direct: "in", Value: 100})
	params = append(params, PROCEDURE_PARAM{Name: "totalcount", Direct: "out", Value: "", ParamPoint: &totalcount})
	params = append(params, PROCEDURE_PARAM{Name: "pagecount", Direct: "out", Value: "", ParamPoint: &pagecount})

	err := ExecProcedureQuery(&dest, "sp_page", params...)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("totalcount", totalcount, "pagecount:", pagecount)
		for i := range dest {
			t.Log(dest[i])
		}
	}
}

func Test_ProcedureExecNoQuery(t *testing.T) {
	Init()
	var param3 int
	var param4 int
	var param5 string
	params := []PROCEDURE_PARAM{}
	params = append(params, PROCEDURE_PARAM{Name: "param1", Direct: "in", Value: 1})
	params = append(params, PROCEDURE_PARAM{Name: "param2", Direct: "in", Value: "2"})
	params = append(params, PROCEDURE_PARAM{Name: "param3", Direct: "out", Value: 0, ParamPoint: &param3})
	params = append(params, PROCEDURE_PARAM{Name: "param4", Direct: "inout", Value: 4, ParamPoint: &param4})
	params = append(params, PROCEDURE_PARAM{Name: "param5", Direct: "inout", Value: "5", ParamPoint: &param5})

	err := ExecProcedureNoQuery("p1", params...)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("p3:", param3, "p4:", param4, "p5:", param5)
	}

}

func Test_ProcedureExecNoQuery1(t *testing.T) {
	Init()
	var param3 int
	var param4 int
	params := []PROCEDURE_PARAM{}
	params = append(params, PROCEDURE_PARAM{Name: "param1", Direct: "in", Value: 1})
	params = append(params, PROCEDURE_PARAM{Name: "param2", Direct: "in", Value: "2"})
	params = append(params, PROCEDURE_PARAM{Name: "param3", Direct: "out", Value: 0, ParamPoint: &param3})
	params = append(params, PROCEDURE_PARAM{Name: "param4", Direct: "out", Value: 4, ParamPoint: &param4})

	err := ExecProcedureNoQuery("p2", params...)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("p3:", param3, "p4:", param4)
	}

}

//func Test_Tx(t *testing.T) {
//	var strSqls []string
//	strSqls = append(strSqls, "insert into tbl2(product_id,product_name) values('1','11111')")
//	strSqls = append(strSqls, "insert into tbl2(product_id,product_name) values('2','22222')")
//	strSqls = append(strSqls, "insert into tbl2(product_id,product_name) values('3','33333')")
//	strSqls = append(strSqls, "insert into tbl2(product_id,product_name) values('4','44444')")
//	strSqls = append(strSqls, "insert into tbl2(product_id,product_name) values('5','55555')")
//	err := ExecNoQueryTx(strSqls)
//	if err != nil {
//		t.Error(err)
//	} else {
//		t.Log("搞定")
//	}
//}
