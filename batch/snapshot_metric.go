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

package batch

import (
	"pmm-ruled/common"
	"pmm-ruled/model"
	"strconv"
	"time"
)

// StartSnapshotBatch snapshot batch
func StartSnapshotBatch() error {
	rowKey := common.ConfigStr["snapshot.row_key"]
	interval := common.ConfigInt["snapshot.interval"]
	tombstoneSec := common.ConfigInt["snapshot.tombstone_sec"]
	for {
		go func() {

			var err error
			var rules []model.SnapshotRule
			var promMetric common.PromMetric

			// Get Snapshot Rules
			rules, _ = (&model.SnapshotRule{}).GetList()
			for _, rule := range rules {

				if promMetric, err = common.Prom.Exec(rule.Query); err != nil {
					common.Log.Error(common.MSG["err.snapshot_bactch_fail"], rule.Query)
					continue
				}

				var metrics []model.SnapshotMetric
				for _, json := range promMetric.Data.Result {

					// parse prometheus result
					ts, _ := json.Value[0].(float64)
					ins := json.Metric[rowKey]
					job := json.Metric["job"]
					num, _ := strconv.ParseFloat(json.Value[1].(string), 64)
					str := json.Metric[*rule.Label]

					snapshotMetric := model.SnapshotMetric{
						RuleID:    rule.ID,
						Instance:  ins,
						Name:      rule.Name,
						Job:       job,
						NumValue:  num,
						StrValue:  str,
						CreatedAt: time.Now(),
						UpdatedAt: time.Unix(int64(ts), 0),
					}
					metrics = append(metrics, snapshotMetric)
					// snapshotMetric.Replace()
				}
				(&model.SnapshotMetric{Name: rule.Name}).ReplaceBulk(metrics)

			}

			// REMOVE unused metrics
			(&model.SnapshotMetric{}).Sweep(tombstoneSec)

		}()
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
