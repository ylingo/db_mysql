// exec_proc_total
package db_mysql

import (
	"bytes"
	"fmt"
	"strings"
)

func QueryWithFoundRows(dest interface{},
	table, fields, where, join, orderby string,
	pageSize, pageIndex int) (totalcount, pagecount, outPageIndex int, err error) {

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

	//结果集
	err = ExecProcedureQuery(dest, "sp_page", params...)
	if err != nil {
		return 0, 0, 0, err
	}

	//处理total，页count
	params = []PROCEDURE_PARAM{}
	params = append(params, PROCEDURE_PARAM{Name: "_table", Direct: "in", Value: table})
	params = append(params, PROCEDURE_PARAM{Name: "_join", Direct: "in", Value: join})
	params = append(params, PROCEDURE_PARAM{Name: "_where", Direct: "in", Value: where})
	params = append(params, PROCEDURE_PARAM{Name: "_pageSize", Direct: "in", Value: pageSize})
	params = append(params, PROCEDURE_PARAM{Name: "found_rows_count", Direct: "out", Value: 0, ParamPoint: &totalcount_p})
	params = append(params, PROCEDURE_PARAM{Name: "calc_page_count", Direct: "out", Value: 0, ParamPoint: &pagecount_p})
	err = ExecProcedureCalcTotal("calc_page", params...)
	if err != nil {
		return 0, 0, 0, err
	}
	totalcount = totalcount_p
	pagecount = pagecount_p
	outPageIndex = pageIndex
	return totalcount, pagecount, outPageIndex, nil
}

func ExecProcedureCalcTotal(procedureName string, params ...PROCEDURE_PARAM) error {
	defer checkerr()
	var outParamsPoint []interface{}
	var _sb bytes.Buffer

	for i := range params {
		if (strings.EqualFold(params[i].Direct, "in") || strings.EqualFold(params[i].Direct, "inout")) &&
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
	_sb.WriteString(procedureName)
	_sb.WriteString("(")
	for i := range params {
		if strings.EqualFold(params[i].Direct, "in") {
			//只写入in参数的值
			if i != 0 {
				_sb.WriteString(",")
			}
			_sb.WriteString("@")
			_sb.WriteString(params[i].Name)
		} else if strings.EqualFold(params[i].Direct, "out") || strings.EqualFold(params[i].Direct, "inout") {
			outParamsPoint = append(outParamsPoint, params[i].ParamPoint)
		}
	}
	_sb.WriteString(");")

	rowCursor := database.QueryRow(_sb.String())
	if err := rowCursor.Scan(outParamsPoint...); err != nil {
		errorinfo("ExecProcedureCalcTotal,%s %s,err:%s", procedureName, _sb.String(), err)
		return err
	}
	loginfo("ExecProcedureCalcTotal, %s", _sb.String())

	return nil
}
