package db_mysql

import (
	"bytes"
	"fmt"
	"strings"
)

func QueryByPage(dest interface{},
	table, fields, where, join, orderby string,
	pageSize, pageIndex int) (totalcount, pagecount int, err error) {

	//先暂时使用这个拙劣的处理方式，后面需要研究一个优雅的方式来处理。
	return QueryWithFoundRows(dest, table, fields, where, join, orderby, pageSize, pageIndex)

	defer checkerr()

	params := []PROCEDURE_PARAM{}
	var totalcount_p int = 0
	var pagecount_p int = 0
	params = append(params, PROCEDURE_PARAM{Name: "_fields", Direct: "in", Value: fields})
	params = append(params, PROCEDURE_PARAM{Name: "_table", Direct: "in", Value: table})
	params = append(params, PROCEDURE_PARAM{Name: "_where", Direct: "in", Value: where})
	params = append(params, PROCEDURE_PARAM{Name: "_join", Direct: "in", Value: join})
	params = append(params, PROCEDURE_PARAM{Name: "_orderby", Direct: "in", Value: orderby})
	params = append(params, PROCEDURE_PARAM{Name: "_pageIndex", Direct: "in", Value: pageIndex})
	params = append(params, PROCEDURE_PARAM{Name: "_pageSize", Direct: "in", Value: pageSize})
	params = append(params, PROCEDURE_PARAM{Name: "_totalcount", Direct: "out", Value: 0, ParamPoint: &totalcount_p})
	params = append(params, PROCEDURE_PARAM{Name: "_pagecount", Direct: "out", Value: 0, ParamPoint: &pagecount_p})

	err = ExecProcedureQuery(dest, "sp_page", params...)
	if err != nil {
		return 0, 0, err
	}
	totalcount = totalcount_p
	pagecount = pagecount_p
	return totalcount, pagecount, nil
}

type PROCEDURE_PARAM struct {
	Name       string
	Direct     string
	Value      interface{}
	ParamPoint interface{}
}

func ExecProcedureQuery(dest interface{}, sqlProcedure string, params ...PROCEDURE_PARAM) error {
	defer checkerr()
	//var inParams []interface{}
	var outParamsPoint []interface{}
	var _sb bytes.Buffer
	outParams := []PROCEDURE_PARAM{}

	for i := range params {
		if (strings.EqualFold(params[i].Direct, "inout") ||
			strings.EqualFold(params[i].Direct, "in")) &&
			params[i].Value != nil {
			_sb.WriteString("set @")
			_sb.WriteString(params[i].Name)
			_sb.WriteString("=")
			switch params[i].Value.(type) {
			case int, int8, int16, int32, int64:
				_sb.WriteString(fmt.Sprintf("%d", params[i].Value))
			case float32, float64:
				_sb.WriteString(fmt.Sprintf("%f", params[i].Value))
			case string:
				_sb.WriteString("\"")

				//由于在执行存储过程的方式是先set变量值，再执行的CALL
				//设置变量时需要将转义符替换成双数，这样在mysql中才不至于出错
				_sb.WriteString(fmt.Sprintf("%s", strings.Replace(params[i].Value.(string), "\\", "\\\\", -1)))
				_sb.WriteString("\"")
			}
			_sb.WriteString(";")
		}
	}
	_sb.WriteString(" call ")
	_sb.WriteString(sqlProcedure)
	_sb.WriteString("(")
	for i := range params {
		if i != 0 {
			_sb.WriteString(",")
		}
		if strings.EqualFold(params[i].Direct, "in") {
			//_sb.WriteString("?")
			//inParams = append(inParams, params[i].Value)
			_sb.WriteString("@")
			_sb.WriteString(params[i].Name)
		} else if strings.EqualFold(params[i].Direct, "out") ||
			strings.EqualFold(params[i].Direct, "inout") {
			_sb.WriteString("@")
			_sb.WriteString(params[i].Name)
			outParamsPoint = append(outParamsPoint, params[i].ParamPoint)
			outParams = append(outParams, params[i])
		}

	}
	_sb.WriteString(");")
	//fmt.Println(_sb.String())
	//err := database.Select(dest, _sb.String(), inParams...)
	err := database.Select(dest, _sb.String())
	if err != nil {
		errorinfo("ExecProcedureQuery, %s ,err:%s", _sb.String(), err)
		return err
	}

	var _sb1 bytes.Buffer
	if len(outParams) > 0 {
		_sb1.WriteString("select ")
		for i := range outParams {
			if i != 0 {
				_sb1.WriteString(",")
			}
			_sb1.WriteString("@")
			_sb1.WriteString(outParams[i].Name)
			_sb1.WriteString(" as ")
			_sb1.WriteString(outParams[i].Name)
		}
		_sb1.WriteString(";") //try append semicolon
		rows1 := database.QueryRow(_sb1.String())
		if err = rows1.Scan(outParamsPoint...); err != nil {
			errorinfo("ExecProcedureQuery,%s %s,err:%s", sqlProcedure, _sb1.String(), err)
			return err
		}
	}
	loginfo("ExecProcedureQuery, %s %s", _sb.String(), _sb1.String())

	return nil
}

func ExecProcedureNoQuery(sqlProcedure string, params ...PROCEDURE_PARAM) error {
	//	执行存储过程时，如果存在inout,则必须是先设置参数形式，不能用？替换符
	//	即 set @param1=1;set @param2='2';set @param4=4;set @param5=5; call p1(@param1,@param2,@param3,@param4,@param5);
	//	   select @param3,@param4,@param5
	//	以上格式
	//	不允许 set @param1; call p1(?,?,@param3)
	//	注：只要有 inout类型，必须先对in，inout类型设置变量，然后call
	// 由于涉及到多个SQL 语句同时执行，需要在go-sql-driver/mysql/packets.go中writeAuthPacket函数增加clientMultiStatements标志，大约在244行
	//var inParams []interface{}
	var outParamsPoint []interface{}
	var _sb bytes.Buffer
	outParams := []PROCEDURE_PARAM{}

	for i := range params {
		if (strings.EqualFold(params[i].Direct, "inout") ||
			strings.EqualFold(params[i].Direct, "in")) &&
			params[i].Value != nil {
			_sb.WriteString("set @")
			_sb.WriteString(params[i].Name)
			_sb.WriteString("=")
			switch params[i].Value.(type) {
			case int, int8, int16, int32, int64:
				_sb.WriteString(fmt.Sprintf("%d", params[i].Value))
			case float32, float64:
				_sb.WriteString(fmt.Sprintf("%f", params[i].Value))
			case string:
				_sb.WriteString("\"")
				//由于在执行存储过程的方式是先set变量值，再执行的CALL
				//设置变量时需要将转义符替换成双数，这样在mysql中才不至于出错
				_sb.WriteString(fmt.Sprintf("%s", strings.Replace(params[i].Value.(string), "\\", "\\\\", -1)))
				_sb.WriteString("\"")
			}
			_sb.WriteString("; ")
		}
	}
	_sb.WriteString(" call ")
	_sb.WriteString(sqlProcedure)
	_sb.WriteString("(")
	for i := range params {
		if i != 0 {
			_sb.WriteString(",")
		}
		if strings.EqualFold(params[i].Direct, "in") {
			//_sb.WriteString("?")
			//inParams = append(inParams, params[i].Value)
			_sb.WriteString("@")
			_sb.WriteString(params[i].Name)
		} else if strings.EqualFold(params[i].Direct, "out") ||
			strings.EqualFold(params[i].Direct, "inout") {
			_sb.WriteString("@")
			_sb.WriteString(params[i].Name)
			outParamsPoint = append(outParamsPoint, params[i].ParamPoint)
			outParams = append(outParams, params[i])
		}
	}
	_sb.WriteString(");")

	//fmt.Println(_sb.String())
	//for i := range inParams {
	//	fmt.Println(inParams[i])
	//}
	//_, err := database.Exec(_sb.String(), inParams...)
	_, err := database.Exec(_sb.String())
	if err != nil {
		errorinfo("ExecProcedureNoQuery, %s ,err:%s", _sb.String(), err)
		return err
	}

	var _sb1 bytes.Buffer
	if len(outParams) > 0 {
		_sb1.WriteString("select ")
		for i := range outParams {
			if i != 0 {
				_sb1.WriteString(",")
			}
			_sb1.WriteString("@")
			_sb1.WriteString(outParams[i].Name)
			_sb1.WriteString(" as ")
			_sb1.WriteString(outParams[i].Name)
		}
		rows1 := database.QueryRow(_sb1.String())
		if err = rows1.Scan(outParamsPoint...); err != nil {
			errorinfo("ExecProcedureNoQuery,%s  %s,err:%s", sqlProcedure, _sb1.String(), err)
			return err
		}
	}
	loginfo("ExecProcedureNoQuery, %s  %s", _sb.String(), _sb1.String())
	return nil
}
