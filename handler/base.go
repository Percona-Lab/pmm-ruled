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
	"pmm-ruled/common"

	"net/http"

	"github.com/gin-gonic/gin"
)

var routerGroup *gin.RouterGroup

// StartAdmin admin server
func StartAdmin() error {
	router := gin.Default()

	// Base router
	base := router.Group(common.ConfigStr["glob.base"])

	// API route
	StartAPI(base.Group("/api/v1"))

	return router.Run(common.ConfigStr["glob.adm_listen_port"])
}

// StartAPI start api server
func StartAPI(r *gin.RouterGroup) {
	startAlertGroupRuleAPI(r)
	startAlertGroupAPI(r)
	startAlertInstanceRuleAPI(r)
	startAlertInstanceSkipRuleAPI(r)
	startAlertInstanceAPI(r)
	startAlertRuleAPI(r)
	startRecordRuleAPI(r)
	startSnapshotRuleAPI(r)

        // Flush rule on start
        FlushAlert()
        FlushRecord()
}

// ErrorIf return boolean if error
func ErrorIf(c *gin.Context, err error) bool {
	if err != nil {
		common.Log.Error(err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"status": "fail",
			"result": err.Error(),
		})
		c.Abort()
		return true
	}
	return false
}

// Success normal message if success
func Success(c *gin.Context, result interface{}) {
	common.Log.Info("Success", result)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": result,
	})
	c.Abort()
}
