// pmm-ruled
// Copyright (C) 2019 gywndi@gmail.com in kakaoBank
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package model

import (
	_ "github.com/go-sql-driver/mysql" // for xorm
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"pmm-ruled/common"
	"time"
)

var orm *xorm.Engine

// NewDatabase new database connection
func NewDatabase() {
	common.Log.Info("db initializing...")

	var err error
	host := common.ConfigStr["db.host"]
	user := common.ConfigStr["db.user"]
	pass := common.ConfigStr["db.pass"]
	db := common.ConfigStr["db.db"]

	orm, err = xorm.NewEngine("mysql", user+":"+pass+"@tcp("+host+")/"+db+"?charset=utf8")
	common.PanicIf(err)

	orm.TZLocation = time.Local
	orm.SetMaxIdleConns(10)
	orm.SetConnMaxLifetime(1 * time.Hour)
	orm.SetMapper(core.GonicMapper{})

	if common.ConfigInt["db.show_sql"] == 1 {
		orm.ShowSQL(true)
	}

	// Initialize
	syncTable()
}

// syncTable sync table and data
func syncTable() {
	err := orm.Sync(new(AlertGroupRule))
	common.PanicIf(err)
	orm.Sync(new(AlertGroup))
	orm.Sync(new(AlertInstanceRule))
	orm.Sync(new(AlertInstanceSkipRule))
	orm.Sync(new(AlertInstance))
	orm.Sync(new(AlertRule))
	orm.Sync(new(RecordRule))
	orm.Sync(new(SnapshotRule))
	orm.Sync(new(SnapshotMetric))
}

// GetDatabase get database orm object to use another package (not common)
func GetDatabase() *xorm.Engine {
	return orm
}

// DBUtil util for db
type DBUtil struct{}
