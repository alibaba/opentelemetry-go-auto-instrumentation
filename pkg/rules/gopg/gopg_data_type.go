package gopg

import "github.com/go-pg/pg/v10/orm"

type gopgRequest struct {
	QueryOp   orm.QueryOp
	System    string
	Statement string
	Addr      string
	User      string
	DbName    string
}
