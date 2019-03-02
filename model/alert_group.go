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
	"fmt"
	"pmm-ruled/common"
	"reflect"
	"regexp"
	"time"
)

// AlertGroup alert group (ex, PROD, DEV)
type AlertGroup struct {
	ID        int       `json:"group_id" xorm:"id int(11) pk not null autoincr"`
	Name      string    `form:"name" json:"name" xorm:"name varchar(32) unique(01) not null "`
	CreatedAt time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// Exist check exists
func (o *AlertGroup) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *AlertGroup) GetFirst() (AlertGroup, error) {
	var err error

	ret := *o
	boolean, err := orm.Get(&ret)
	if err != nil {
		return ret, err
	}

	if !boolean {
		return ret, fmt.Errorf(common.MSG["err.row_not_found"])
	}

	return ret, err
}

// GetList get rows
func (o *AlertGroup) GetList(sort ...string) ([]AlertGroup, error) {
	var err error
	var arr []AlertGroup
	var order string

	for i, s := range sort {
		if i > 0 {
			order += ","
		}
		order += s
	}
	err = orm.OrderBy(order).Find(&arr, o)
	common.Log.Info(reflect.TypeOf(o), len(arr), " rows selected")
	return arr, err
}

// Insert new row
func (o *AlertGroup) Insert() error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = o.InsertCheck(); err != nil {
		return err
	}

	if affected, err = session.Insert(o); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "rows inserted!")

	return err
}

// Update update row (partitial column)
func (o *AlertGroup) Update(to *AlertGroup) (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = to.UpdateCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Update(to, o); err != nil {
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows updated!")
	return affected, err
}

// Delete delete row
func (o *AlertGroup) Delete() (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	if err = o.DeleteCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Delete(o); err != nil {
		return 0, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows deleted!")

	// delete AlertGroup Alert
	orm.Delete(&AlertGroupRule{GroupID: o.ID})

	// update instance "group_id" to zero
	orm.Update(&AlertInstance{GroupID: new(int)})

	return affected, err
}

// InsertCheck validation check
func (o *AlertGroup) InsertCheck() error {
	var err error

	// Empty check
	if o.Name == "" {
		return fmt.Errorf(common.MSG["err.name_empty"])
	}

	return err
}

// UpdateCheck validation check
func (o *AlertGroup) UpdateCheck() error {
	var err error

	// Empty check
	if o.Name == "" {
		return fmt.Errorf(common.MSG["err.name_empty"])
	}

	return err
}

// DeleteCheck validation check
func (o *AlertGroup) DeleteCheck() error {
	var err error

	return err
}

// rewriteCols rewrite column value
func (o *AlertGroup) rewriteCols() {
	o.Name = regexp.MustCompile(`\s`).ReplaceAllString(o.Name, "_")
}

// // GetRules get rules in alert group
// func (o *AlertGroup) GetRules() []AlertRuleGroupVal {
// 	var r []AlertRuleGroupVal

// 	// Check group  exists
// 	if o.Exist() {
// 		orm.Sql(`
// 			select
// 			  alert_rule.*,
// 			  alert_group_rule.val as group_val
// 			from alert_group_rule
// 			inner join alert_rule on alert_rule.id = alert_group_rule.rule_id
// 			where alert_group_rule.group_id = ?
// 			order by alert_rule.name
// 		`, o.ID).Find(&r)
// 	}
// 	return r
// }

// GetRules get rules in alert group
func (o *AlertGroup) GetRules() []map[string]string {
	results, _ := orm.QueryString(`
	select 
		alert_rule.id       as rule_id,
		alert_rule.name     as name,
		alert_rule.level    as level,
		alert_rule.val      as val,
		alert_group_rule.val as group_val
	from alert_group_rule
	inner join alert_rule on alert_rule.id = alert_group_rule.rule_id
	where alert_group_rule.group_id = ?
	order by alert_rule.name`, o.ID)
	return results
}

// GetInstances get instances in alert group
func (o *AlertGroup) GetInstances() []AlertInstance {
	var r []AlertInstance
	if o.Exist() {
		r, _ = (&AlertInstance{GroupID: &o.ID}).GetList()
	}
	return r
}
