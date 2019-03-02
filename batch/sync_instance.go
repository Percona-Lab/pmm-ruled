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
	"time"
)

// StartInstanceBatch Instanceshot batch
func StartInstanceBatch() error {
	for {
		go func() {

			// instance map
			instanceMap := make(map[string]string)
			instances, _ := (&model.AlertInstance{}).GetList()
			for _, instance := range instances {
				instanceMap[instance.Name] = instance.Name
			}

			// insert intance only if no instance in hash map
			if metric, err := common.Prom.Exec("count(up) by (instance)"); err == nil {
				for _, json := range metric.Data.Result {
					instance := json.Metric["instance"]
					if instanceMap[instance] == "" {
						(&model.AlertInstance{Name: instance}).Insert()
					}
				}
			}
		}()
		time.Sleep(1 * time.Minute)
	}
}
