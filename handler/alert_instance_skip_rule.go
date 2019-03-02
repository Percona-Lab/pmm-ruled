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

// startAlertInstanceSkipRuleAPI alert instance skip rule API
func startAlertInstanceSkipRuleAPI(r *gin.RouterGroup) {

	// new
	r.POST("/alert/skip_rule/:instance_id/:rule_id", func(c *gin.Context) {
		var err error
		var params model.AlertInstanceSkipRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		// get id
		params.InstanceID = common.ParseInt(c.Param("instance_id"))
		params.RuleID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.InstanceID == 0 || params.RuleID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertInstanceSkipRule{InstanceID: params.InstanceID, RuleID: params.RuleID}

		// check exists
		if target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_exists"]))
			return
		}

		// insert
		err = params.Insert()
		if ErrorIf(c, err) {
			return
		}

		Success(c, "ok")
	})

	// delete
	r.DELETE("/alert/skip_rule/:instance_id/:rule_id", func(c *gin.Context) {
		var err error
		var params model.AlertInstanceSkipRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		// get id
		params.InstanceID = common.ParseInt(c.Param("instance_id"))
		params.RuleID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.InstanceID == 0 || params.RuleID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertInstanceSkipRule{InstanceID: params.InstanceID, RuleID: params.RuleID}

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

		Success(c, "ok")
	})

}
