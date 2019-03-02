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

// AlertInstance AlertInstances motored
type AlertInstance struct {
	ID        int       `json:"instance_id" xorm:"id int(11) pk not null autoincr"`
	Name      string    `form:"name" json:"name" xorm:"name varchar(32) unique(01) not null "`
	GroupID   *int      `form:"group_id" json:"group_id" xorm:"group_id not null index(02) default 0"`
	CreatedAt time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// AlertInstanceExt alert instance with additional info
type AlertInstanceExt struct {
	AlertInstance `xorm:"extends"`
	GroupName     string `json:"group_name"`
}

// Exist check exists
func (o *AlertInstance) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *AlertInstance) GetFirst() (AlertInstance, error) {
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
func (o *AlertInstance) GetList(sort ...string) ([]AlertInstance, error) {
	var err error
	var arr []AlertInstance
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
func (o *AlertInstance) Insert() error {
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
func (o *AlertInstance) Update(to *AlertInstance) (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = to.UpdateCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Update(to, o); err != nil {
		common.Log.Error(err)
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows updated!")
	return affected, err
}

// Delete delete row
func (o *AlertInstance) Delete() (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	if err = o.DeleteCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Delete(o); err != nil {
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows deleted!")

	// delete AlertInstance Alert
	orm.Delete(&AlertInstanceRule{InstanceID: o.ID})

	// delete AlertInstance Alert
	orm.Delete(&AlertInstanceSkipRule{InstanceID: o.ID})

	return affected, err
}

// InsertCheck validation check
func (o *AlertInstance) InsertCheck() error {
	var err error

	if o.GroupID == nil {
		o.GroupID = new(int)
	}

	if *o.GroupID != 0 && !(&AlertGroup{ID: *o.GroupID}).Exist() {
		return fmt.Errorf(common.MSG["err.group_not_exists"])
	}

	return err
}

// UpdateCheck validation check
func (o *AlertInstance) UpdateCheck() error {
	var err error

	if o.GroupID != nil && *o.GroupID != 0 {
		if !(&AlertGroup{ID: *o.GroupID}).Exist() {
			return fmt.Errorf(common.MSG["err.group_not_exists"])
		}
	}

	return err
}

// DeleteCheck validation check
func (o *AlertInstance) DeleteCheck() error {
	var err error

	return err
}

// rewriteCols rewrite column value
func (o *AlertInstance) rewriteCols() {
	o.Name = regexp.MustCompile(`\s`).ReplaceAllString(o.Name, "_")
}

// GetInstanceExt get instance list with group name
func (o *AlertInstance) GetInstanceExt() AlertInstanceExt {
	var r AlertInstanceExt
	if o.Exist() {
		var rows []AlertInstanceExt
		orm.Sql(`
			select alert_instance.*, ifnull(alert_group.name, '') group_name
			from alert_instance
			left join alert_group on alert_group.id = alert_instance.group_id
			where alert_instance.id = ?
		`, o.ID).Find(&rows)
		r = rows[0]
	}
	return r
}

// GetInstanceExtList get instance list with group name
func (o *AlertInstance) GetInstanceExtList() []AlertInstanceExt {
	var r []AlertInstanceExt
	orm.Sql(`
		select alert_instance.*, ifnull(alert_group.name, '') group_name
		from alert_instance
		left join alert_group on alert_group.id = alert_instance.group_id
	`).Find(&r)
	return r
}

// // GetRules get rules in alert group
// func (o *AlertInstance) GetRules() []AlertRuleInstanceVal {
// 	var r []AlertRuleInstanceVal
// 	if o.Exist() {
// 		orm.Sql(`
// 			select
// 				alert_rule.*,
// 				alert_instance_rule.val as instance_val
// 			from alert_instance_rule
// 			inner join alert_rule on alert_rule.id = alert_instance_rule.rule_id
// 			where alert_instance_rule.instance_id = ?
// 			order by alert_rule.name
// 		`, o.ID).Find(&r)
// 	}
// 	return r
// }

// GetRules get rules in alert group
func (o *AlertInstance) GetRules() []map[string]string {
	results, _ := orm.QueryString(`
	select 
		alert_rule.id       as rule_id,
		alert_rule.name     as name,
		alert_rule.level    as level,
		alert_rule.val      as val,
		alert_instance_rule.val as instance_val
	from alert_instance_rule
	inner join alert_rule on alert_rule.id = alert_instance_rule.rule_id
	where alert_instance_rule.instance_id = ?
	order by alert_rule.name`, o.ID)
	return results
}

// // GetSkipRules get skip rules in instance
// func (o *AlertInstance) GetSkipRules() []AlertRuleMini {
// 	var r []AlertRuleMini
// 	if o.Exist() {
// 		orm.Sql(`
// 			select
// 				alert_rule.*
// 			from alert_instance_skip_rule
// 			inner join alert_rule on alert_rule.id = alert_instance_skip_rule.rule_id
// 			where alert_instance_skip_rule.instance_id = ?
// 			order by alert_rule.name
// 		`, o.ID).Find(&r)
// 	}
// 	return r
// }

// GetSkipRules get skip rules in instance
func (o *AlertInstance) GetSkipRules() []map[string]string {
	results, _ := orm.QueryString(`
	select 
		alert_rule.*
	from alert_instance_skip_rule
	inner join alert_rule on alert_rule.id = alert_instance_skip_rule.rule_id
	where alert_instance_skip_rule.instance_id = ?
	order by alert_rule.name`, o.ID)
	return results
}
