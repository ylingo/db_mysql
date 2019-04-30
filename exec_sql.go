package db_mysql

import (
	_ "strings"
)

func ExecNoQuery(strStmt string, args ...interface{}) error {
	sqlStmt, err := database.Prepare(strStmt)
	defer checkerr()
	defer sqlStmt.Close()
	if err != nil {
		errorinfo("ExecNoQuery,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return err
	}
	_, err = sqlStmt.Exec(args...)
	if err != nil {
		errorinfo("ExecNoQuery,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return err
	}
	loginfo("ExecNoQuery，sql:%s,param:%s", strStmt, argsToString(args...))
	return nil
	//	if _, err := database.Exec(strStmt, args...); err != nil {
	//		return err
	//	}
	//	return nil
}

func ExecInsertGetLastId(strStmt string, args ...interface{}) (int64, error) {
	sqlStmt, err := database.Prepare(strStmt)
	defer checkerr()
	defer sqlStmt.Close()
	if err != nil {
		errorinfo("ExecInsertGetLastId,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return 0, err
	}
	sqlResult, err := sqlStmt.Exec(args...)
	if err != nil {
		errorinfo("ExecInsertGetLastId,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return 0, err
	}
	id, err := sqlResult.LastInsertId()
	if err != nil {
		errorinfo("ExecInsertGetLastId,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return 0, err
	}
	loginfo("ExecInsertGetLastId,sql:%s,param:%s", strStmt, argsToString(args...))
	return id, nil
}

func ExecQuery(dest interface{}, strStmt string, args ...interface{}) error {
	defer checkerr()
	//方法一
	if err := database.Select(dest, strStmt, args...); err != nil {
		errorinfo("ExecQuery,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
		return err
	}
	loginfo("ExecQuery,sql:%s,param:%s", strStmt, argsToString(args...))
	return nil

	//方法二
	//	rows, err := database.Query(strStmt, args...)
	//	defer rows.Close()
	//	if err != nil {
	//		errorinfo("ExecQuery,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
	//		return err
	//	}
	//	if err := sqlx.StructScan(rows, dest); err != nil {
	//		errorinfo("ExecQuery,sql:%s,param:%s,err:%s", strStmt, argsToString(args...), err)
	//		return err
	//	}
	//  loginfo("ExecQuery,sql:%s,param:%s", strStmt, argsToString(args...))
	//	return nil
}

func ExecCheckExists(strStrmt string, args ...interface{}) (bool, error) {
	defer checkerr()
	strSql := "select  coalesce(exists(" + strStrmt + "),0) as tt"
	var tmp []string
	if err := ExecQuery(&tmp, strSql, args...); err != nil {
		return false, err
	}
	if len(tmp) > 0 && tmp[0] == "1" {
		return true, nil
	}
	return false, nil
}
