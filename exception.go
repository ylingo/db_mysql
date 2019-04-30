package db_mysql

func checkerr() {
	if err := recover(); err != nil {
		errorinfo(err)
	}
}
