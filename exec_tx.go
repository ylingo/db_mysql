package db_mysql

//	"bytes"
//	"fmt"
//	"strings"

//func ExecNoQueryTx(strStmts []string) error {
//	var _sb bytes.Buffer
//	for i := range strStmts {
//		_sb.WriteString(strStmts[i])
//		if !strings.HasSuffix(strStmts[i], ";") {
//			_sb.WriteString(";")
//		}
//	}

//	fmt.Println("aaaa")
//	tx, err := database.Beginx()
//	fmt.Println("bbbb")
//	if err != nil {
//		return err
//	}

//	sqlStmt, err := tx.Prepare(_sb.String())
//	defer sqlStmt.Close()
//	defer tx.Rollback()
//	if err != nil {
//		return err
//	}

//	_, err = sqlStmt.Exec()
//	if err != nil {
//		return err
//	}

//	tx.Commit()
//	return nil
//}
