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

package handler

import (
	"fmt"
	"pmm-ruled/common"
	"pmm-ruled/model"

	"github.com/gin-gonic/gin"
)

// startSnapshotRuleAPI snapshot rule API
func startSnapshotRuleAPI(r *gin.RouterGroup) {

	// new
	r.POST("/snapshot/rule", func(c *gin.Context) {
		var err error
		var params model.SnapshotRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		// insert
		err = params.Insert()
		if ErrorIf(c, err) {
			return
		}

		Success(c, params.ID)
	})

	// update
	r.PUT("/snapshot/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.SnapshotRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.SnapshotRule{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_not_exists"]))
			return
		}

		// update
		_, err = target.Update(&params)
		if ErrorIf(c, err) {
			return
		}

		Success(c, params.ID)

	})

	// delete
	r.DELETE("/snapshot/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.SnapshotRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.SnapshotRule{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_not_exists"]))
			return
		}

		// delete
		_, err = target.Delete()
		if ErrorIf(c, err) {
			return
		}

		Success(c, params.ID)

	})

	// get one
	r.GET("/snapshot/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.SnapshotRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.SnapshotRule{ID: params.ID}

		// get first one
		if params, err = target.GetFirst(); err != nil {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_not_exists"]))
			return
		}

		Success(c, params)
	})

	// get list
	r.GET("/snapshot/rules", func(c *gin.Context) {
		var err error
		var params model.SnapshotRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}
		list, err := params.GetList("name")
		Success(c, list)
	})

	// get metrics
	r.GET("/snapshot/rule/:rule_id/metrics", func(c *gin.Context) {
		var err error
		var snapshotMetric model.SnapshotRuleMetric

		// get id
		id := common.ParseInt(c.Param("rule_id"))

		// check ID
		if id == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		if err = snapshotMetric.Get(id); err != nil {
			ErrorIf(c, fmt.Errorf(common.MSG["err.get_metric_fail"]))
			return
		}

		Success(c, snapshotMetric)
	})
}
