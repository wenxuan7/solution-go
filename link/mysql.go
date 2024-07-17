package link

import "strings"

type mysqlCon struct {
	username string
	password string
}

func dsn() string {
	mysqlCon := mysqlCon{username: "root", password: "wenxuan101314"}
	sb := strings.Builder{}
	// root:wenxuan101314@tcp(127.0.0.1:3306)/raycloud_erp?charset=utf8&parseTime=True&loc=Local
	sb.WriteString(mysqlCon.username)
	sb.WriteString(":")
	sb.WriteString(mysqlCon.password)
	sb.WriteString("@tcp(127.0.0.1:3306)/raycloud_erp?charset=utf8&parseTime=True&loc=Local")
	return sb.String()
}
