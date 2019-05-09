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

func Test_ExecInsertGetLastId(t *testing.T) {
	Init()

	lastId, err := ExecInsertGetLastId("insert into tbl1(c1,c2,c3) values(?,?,?)", "a", "\"b", "\"c")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("执行成功,ID:", lastId)
	}
}

type Product struct {
	SeqId     int64  `db:"seqid" json:"seqid"`
	UserId    string `db:"userid" json:"userid"`
	UserName  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password"`
	MobileNum string `db:"mobilenum" json:"mobilenum"`
	OnOff     bool   `db:"onoff" json:"onoff"`
}

func Test_ExecQuery(t *testing.T) {
	Init()
	seqId := 10
	p := []Product{}
	strSql := `select 
		t.seqid,
		coalesce(t.userId,'') as userid,
		coalesce(t.userName,'') as username,
		'' as password,
		coalesce(t.mobileNum,'') as mobilenum,
		coalesce(t.onoff,false) as onoff
		from t_admin_user as t
		where t.seqid=?`
	err := ExecQuery(&p, strSql, seqId)
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
	tblName := "t_admin_user as t"
	fields := `t.seqid,
	    coalesce(t.userId,'') as userid,
		coalesce(t.userName,'') as username,
		'' as password,
		coalesce(t.mobileNum,'') as mobilenum,
		coalesce(t.onOff,false) as onoff`

	join := ""
	where := ""
	orderby := ""
	totalcount, totalpage, outpageindex, err := QueryByPage(&dest, tblName, fields, where,
		join,
		orderby, 100, 1)
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
