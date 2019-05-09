package db_mysql

import (
	"testing"
)

func Init() {
	InitDb("mysql",
		"root:root@tcp(127.0.0.1:3306)/abc?charset=utf8&parseTime=true&timeout=20s&loc=Local&clientFoundRows=true",
		100,
		20, true, true)
}

func _Test_ExecNoQuery(t *testing.T) {
	Init()
	//	p := []Product{}
	//	ExecQuery(&p, "select Id,Product_id,Product_name from tbl2 limit 1,100")
	//	for i := range p {
	//		t.Log(p[i])
	//	}

	//err := ExecNoQuery("insert into tbl2(product_id,product_name) values('111','2222')")
	//strsql := "insert into tbl1(c2,c3,c4) values('1','11111','sssssss');"
	strsql := "insert into tbl2(c22,c33,c44) values('2','99999','hhhhhhhh');"
	// strsql += "insert into tbl1(c2,c3,c4) values('3','33333','sssssss');"
	// strsql += "insert into tbl1(c2,c3,c4) values('4','44444','sssssss');"
	// strsql += "insert into tbl1(c2,c3,c4) values('5','55555','sssssss');"
	for i := 0; i < 100; i++ {
		err := ExecNoQuery(strsql)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("执行成功")
		}
	}
}

// func Test_ExecInsertGetLastId(t *testing.T) {
// 	Init()

// 	lastId, err := ExecInsertGetLastId("insert into tbl2(product_id,product_name) values('111','2222')")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log("执行成功,ID:", lastId)
// 	}
// }

type Product struct {
	C1  int    `db:"c1"`
	C2  string `db:"c2"`
	C3  string `db:"c3"`
	C4  string `db:"c4"`
	C33 string `db:"c33"`
	C44 string `db:"c44"`
}

// func Test_ExecQuery(t *testing.T) {
// 	Init()

// 	p := []Product{}
// 	//p := make([]Product, 0)
// 	err := ExecQuery(&p, "select id,product_id,product_name from tbl2 limit 1,100")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		for i := range p {
// 			t.Log(p[i])
// 		}
// 	}
// }

func _TestQueryByPage(t *testing.T) {
	Init()
	dest := []Product{}
	totalcount, totalpage, outpageindex, err := QueryByPage(&dest, "tbl1 as t", "t.c1,t.c2,t.c3,t.c4,t2.c33,t2.c44", "",
		"left join tbl2 as t2 on t.c1=t2.c11",
		"", 100, 1)
	if err != nil {
		t.Error("totalcount", totalcount, "totalpage:", totalpage, "err:", err)
	} else {
		t.Log("totalcount", totalcount, ",totalpage:", totalpage, "outpageindex:", outpageindex)
		t.Log(len(dest))
		for i := range dest {
			t.Log(dest[i])
		}
	}
}

func TestBatchExecNoQuery(t *testing.T) {
	Init()
	var strSqls []string
	strSqls = append(strSqls, "insert into tbl1(c1,c2,c3) values('a','b','c')")
	strSqls = append(strSqls, "insert into tbl1(c1,c2,c3) values('a1','b1','c1')")
	if err := ExecBatchNoQuery(strSqls); err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
