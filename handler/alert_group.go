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

// startAlertGroupAPI alert group API
func startAlertGroupAPI(r *gin.RouterGroup) {

	// new
	r.POST("/alert/group", func(c *gin.Context) {
		var err error
		var params model.AlertGroup

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
	r.PUT("/alert/group/:group_id", func(c *gin.Context) {
		var err error
		var params model.AlertGroup

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		// get id
		params.ID = common.ParseInt(c.Param("group_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertGroup{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.group_not_exists"]))
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
	r.DELETE("/alert/group/:group_id", func(c *gin.Context) {
		var err error
		var params model.AlertGroup

		// get id
		params.ID = common.ParseInt(c.Param("group_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertGroup{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.group_not_exists"]))
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
	r.GET("/alert/group/:group_id", func(c *gin.Context) {
		var err error
		var params model.AlertGroup

		// get id
		params.ID = common.ParseInt(c.Param("group_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertGroup{ID: params.ID}

		// get first one
		if params, err = target.GetFirst(); err != nil {
			ErrorIf(c, fmt.Errorf(common.MSG["err.group_not_exists"]))
			return
		}

		Success(c, params)

	})

	// get list
	r.GET("/alert/groups", func(c *gin.Context) {
		var err error
		var params model.AlertGroup

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		list, err := params.GetList("name")
		Success(c, list)
	})

	// get alert rules
	r.GET("/alert/group/:group_id/rules", func(c *gin.Context) {
		var params model.AlertGroup

		// get id
		params.ID = common.ParseInt(c.Param("group_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertGroup{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.group_not_exists"]))
			return
		}

		Success(c, params.GetRules())
	})

	// get alert rules
	r.GET("/alert/group/:group_id/instances", func(c *gin.Context) {
		var params model.AlertGroup

		// get id
		params.ID = common.ParseInt(c.Param("group_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertGroup{ID: params.ID}

		// check exists
		if !target.Exist() {
			ErrorIf(c, fmt.Errorf(common.MSG["err.group_not_exists"]))
			return
		}

		Success(c, params.GetInstances())
	})
}
