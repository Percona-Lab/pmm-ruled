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
	"strconv"
	"strings"
	"time"
)

// AlertRule alert rule
type AlertRule struct {
	ID          int       `json:"rule_id" xorm:"id int(11) pk not null autoincr"`
	Name        string    `form:"name" json:"name" xorm:"varchar(32) unique(01) not null "`
	Level       string    `form:"level" json:"level" xorm:"varchar(10) unique(01) not null "`
	Wait        *int      `form:"wait" json:"wait" xorm:"int(11) not null default 0"` // resettable
	Rule        string    `form:"rule" json:"rule" xorm:"text not null "`
	Opr         string    `form:"opr" json:"opr" xorm:"varchar(2) not null"`
	Val         *string   `form:"val" json:"val" xorm:"val varchar(30) not null default '' "` // resettable
	Subject     string    `form:"subject" json:"subject" xorm:"varchar(64) not null "`
	Description string    `form:"description" json:"description" xorm:"text not null "`
	CreatedAt   time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt   time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// // AlertRuleGroupVal  alert rule with group val
// type AlertRuleGroupVal struct {
// 	ID      int    `json:"rule_id" xorm:"id"`
// 	Name    string `json:"name"`
// 	Level   string `json:"level"`
// 	Val     string `json:"val"`
// 	GroupVal string `json:"group_val" xorm:"group_val"`
// }

// // AlertRuleInstanceVal alert rule with instance val
// type AlertRuleInstanceVal struct {
// 	ID      int    `json:"rule_id" xorm:"id"`
// 	Name    string `json:"name"`
// 	Level   string `json:"level"`
// 	Val     string `json:"val"`
// 	InstanceVal string `json:"instance_val"  xorm:"instance_val"`
// }

// // AlertRuleMini alert rule minimal info
// type AlertRuleMini struct {
// 	ID    int    `json:"rule_id" xorm:"id"`
// 	Name  string `json:"name"`
// 	Level string `json:"level"`
// }

// AlertThreshold alert thresolds info
type AlertThreshold struct {
	GroupID      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	InstanceID   int    `json:"instance_id"`
	InstanceName string `json:"instance_name"`
	RuleID       int    `json:"rule_id"`
	RuleName     string `json:"rule_name"`
	Level        string `json:"level"`
	Val          string `json:"val"`
	Activate     int    `json:"activate"`
}

// Exist check exists
func (o *AlertRule) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *AlertRule) GetFirst() (AlertRule, error) {
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
func (o *AlertRule) GetList(sort ...string) ([]AlertRule, error) {
	var err error
	var arr []AlertRule
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
func (o *AlertRule) Insert() error {
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
func (o *AlertRule) Update(to *AlertRule) (int64, error) {
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
func (o *AlertRule) Delete() (int64, error) {
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
	return affected, err
}

// InsertCheck validation check
func (o *AlertRule) InsertCheck() error {
	var err error

	// Empty check
	if o.Name == "" {
		return fmt.Errorf(common.MSG["err.name_empty"])
	}

	if o.Level == "" {
		return fmt.Errorf(common.MSG["err.level_empty"])
	}

	if o.Rule == "" {
		return fmt.Errorf(common.MSG["err.rule_empty"])
	}

	if o.Opr == "" {
		return fmt.Errorf(common.MSG["err.opr_empty"])
	}

	if o.Subject == "" {
		return fmt.Errorf(common.MSG["err.subj_empty"])
	}

	if o.Description == "" {
		return fmt.Errorf(common.MSG["err.desc_empty"])
	}

	if o.Wait == nil {
		o.Wait = new(int)
	}

	if o.Val == nil {
		return fmt.Errorf(common.MSG["err.val_empty"])
	}

	// Val digit check
	if o.Val != nil && *o.Val != "" {
		if _, err := strconv.ParseFloat(*o.Val, 64); err != nil {
			return fmt.Errorf(common.MSG["err.val_not_digit"])
		}
	}

	// prometheus syntax check
	q := fmt.Sprintf("(%s) < -1", o.Rule)
	common.Log.Info("rule check - ", q)
	if _, err = common.Prom.Exec(q); err != nil {
		return err
	}

	return err
}

// UpdateCheck validation check
func (o *AlertRule) UpdateCheck() error {
	var err error

	if o.Opr != "" && o.Opr != `<` && o.Opr != `<=` && o.Opr != `>` && o.Opr != `>=` && o.Opr != `!=` && o.Opr != `==` {
		return fmt.Errorf(common.MSG["err.invalid_operator"], o.Opr)
	}

	// Val digit check
	if o.Val != nil {
		if _, err := strconv.ParseFloat(*o.Val, 64); err != nil {
			return fmt.Errorf(common.MSG["err.val_not_digit"])
		}
	}

	// prometheus syntax check
	if o.Rule != "" {
		q := fmt.Sprintf("(%s) < -1", o.Rule)
		common.Log.Info("rule check - ", q)
		if _, err = common.Prom.Exec(q); err != nil {
			return err
		}
	}

	return err
}

// DeleteCheck validation check
func (o *AlertRule) DeleteCheck() error {
	var err error

	return err
}

// rewriteCols rewrite column value
func (o *AlertRule) rewriteCols() {
	o.Name = regexp.MustCompile(`\s`).ReplaceAllString(o.Name, "_")
	if o.Val != nil {
		*o.Val = strings.TrimSpace(*o.Val)
	}
}

// GetAlertThresoldList get alert thresold list
func (o *AlertRule) GetAlertThresoldList() []AlertThreshold {
	var rules []AlertThreshold
	orm.Sql(`
		select straight_join
			alert_group.id           as group_id
			,alert_group.name        as group_name
			,alert_instance.id          as instance_id
			,alert_instance.name        as instance_name
			,alert_rule.id          as rule_id
			,alert_rule.name        as rule_name
			,alert_rule.level       as level
			,case 
				when ifnull(e.val,'') != '' then e.val
				when ifnull(alert_group_rule.val,'') != '' then alert_group_rule.val
				else alert_rule.val
			end as val
			,if(f.instance_id is null, 1, 0)  as activate
		from(
			select 
				a.id instance_id,
				b.rule_id rule_id
			from alert_instance a
			inner join alert_group_rule b on b.group_id = a.group_id
			union
			select 
				instance_id, 
				rule_id
			from alert_instance_rule
		) t
		inner join alert_rule      on alert_rule.id = t.rule_id
		inner join alert_instance      on alert_instance.id = t.instance_id
		inner join alert_group      on alert_group.id = alert_instance.group_id
		inner join alert_group_rule on alert_group_rule.group_id = alert_group.id and alert_group_rule.rule_id = alert_rule.id
		left join alert_instance_rule      e on e.instance_id = alert_instance.id and e.rule_id = alert_rule.id
		left join alert_instance_skip_rule f on f.instance_id = alert_instance.id and f.rule_id = alert_rule.id
	`).Find(&rules)
	return rules
}
