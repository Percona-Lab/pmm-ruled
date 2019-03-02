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

// SnapshotRule snapshot rule
type SnapshotRule struct {
	ID        int       `json:"rule_id" xorm:"id int(11) pk not null autoincr"`
	Name      string    `form:"name" json:"name" xorm:"varchar(32) unique(01) not null "`
	Label     *string   `form:"label" json:"label" xorm:"varchar(64) not null default ''"` // resettable
	Query     string    `form:"query" json:"query" xorm:"text not null "`
	CreatedAt time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// Exist check exists
func (o *SnapshotRule) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *SnapshotRule) GetFirst() (SnapshotRule, error) {
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
func (o *SnapshotRule) GetList(sort ...string) ([]SnapshotRule, error) {
	var err error
	var arr []SnapshotRule
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
func (o *SnapshotRule) Insert() error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = o.InsertCheck(); err != nil {
		common.Log.Error(err)
		return err
	}

	if affected, err = session.Insert(o); err != nil {
		common.Log.Error(err)
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "rows inserted!")

	return err
}

// Update update row (partitial column)
func (o *SnapshotRule) Update(to *SnapshotRule) (int64, error) {
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
func (o *SnapshotRule) Delete() (int64, error) {
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
func (o *SnapshotRule) InsertCheck() error {
	var err error

	// Empty check
	if o.Name == "" {
		return fmt.Errorf(common.MSG["err.name_empty"])
	}

	if o.Query == "" {
		return fmt.Errorf(common.MSG["err.query_empty"])
	}

	if o.Label == nil {
		o.Label = new(string)
	}

	// prometheus syntax check
	common.Log.Info("Prometheus rule check - ", o.Query)
	if _, err = common.Prom.Exec(o.Query); err != nil {
		return err
	}

	return err
}

// UpdateCheck validation check
func (o *SnapshotRule) UpdateCheck() error {
	var err error

	// prometheus syntax check
	if o.Query != "" {
		common.Log.Info("Prometheus rule check - ", o.Query)
		if _, err = common.Prom.Exec(o.Query); err != nil {
			return err
		}
	}

	return err
}

// DeleteCheck validation check
func (o *SnapshotRule) DeleteCheck() error {
	var err error

	return err
}

// rewriteCols rewrite column value
func (o *SnapshotRule) rewriteCols() {
	o.Name = regexp.MustCompile(`\s`).ReplaceAllString(o.Name, "_")
}

// SnapshotRuleMetric get metrics for rule
type SnapshotRuleMetric struct {
	Rule    SnapshotRule     `json:"rule"`
	Metrics []SnapshotMetric `json:"metrics"`
}

// Get get metrics for rule
func (o *SnapshotRuleMetric) Get(id int) error {
	var err error
	if o.Rule, err = (&SnapshotRule{ID: id}).GetFirst(); err != nil {
		return err
	}
	o.Metrics, _ = (&SnapshotMetric{RuleID: o.Rule.ID}).GetList("name")
	return err
}
