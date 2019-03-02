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
	"pmm-ruled/common"
	"reflect"
	"time"
)

// SnapshotMetric metric snapshot result
type SnapshotMetric struct {
	RuleID    int       `json:"rule_id" xorm:"rule_id pk not null "`
	Instance  string    `json:"instance" xorm:"instance char(32) pk not null index(01) "`
	Name      string    `json:"name" xorm:"name varchar(32) not null "`
	Job       string    `json:"job" xorm:"varchar(20) not null "`
	NumValue  float64   `json:"num_value" xorm:"double not null default 0 "`
	StrValue  string    `json:"str_value" xorm:"varchar(255) not null default '' "`
	CreatedAt time.Time `json:"created_at" xorm:"datetime not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"datetime not null updated"`
}

// GetList get all rows
func (o *SnapshotMetric) GetList(sort ...string) ([]SnapshotMetric, error) {
	var err error
	var arr []SnapshotMetric
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

// Replace delete and insert row
func (o *SnapshotMetric) Replace() error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	cond := SnapshotMetric{RuleID: o.RuleID, Instance: o.Instance}
	if affected, err = session.Delete(cond); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "deleted!")

	if _, err = session.Insert(o); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "rows inserted!")

	return err
}

// ReplaceBulk delete all and insert bulk rows
func (o *SnapshotMetric) ReplaceBulk(metrics []SnapshotMetric) error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	if affected, err = session.Delete(o); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "deleted!")

	var rows int
	for _, metric := range metrics {
		if _, err = session.Insert(metric); err == nil {
			rows = rows + 1
		}
	}
	common.Log.Info(reflect.TypeOf(o), rows, "rows inserted!")

	return err
}

// Sweep sweep used data
func (o *SnapshotMetric) Sweep(sec int) {
	var err error
	// delete snapshot metric not in snapshot rule
	_, err = orm.Exec(`
		delete snapshot_metric
		from snapshot_metric
		left join alert_instance on alert_instance.name = snapshot_metric.instance
		left join snapshot_rule on snapshot_rule.id = snapshot_metric.rule_id
		where snapshot_rule.id is null`)
	if err != nil {
		common.Log.Error(err)
	}

	// delete snapshot metric not gathered too long
	orm.Exec(`
		delete snapshot_metric
		from (
		select instance, job, max(updated_at) updated_at
		from snapshot_metric
		group by instance, job
		) m
		inner join snapshot_metric on snapshot_metric.instance = m.instance and snapshot_metric.job = m.job
		where unix_timestamp(m.updated_at) - unix_timestamp(snapshot_metric.updated_at) > ?
	`, sec)
	if err != nil {
		common.Log.Error(err)
	}
}
