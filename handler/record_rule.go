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
	"os"
	"os/exec"
	"pmm-ruled/common"
	"pmm-ruled/model"

	"github.com/gin-gonic/gin"
)

// startRecordRuleAPI record rule API
func startRecordRuleAPI(r *gin.RouterGroup) {

	// new
	r.POST("/record/rule", func(c *gin.Context) {
		var err error
		var params model.RecordRule

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
	r.PUT("/record/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.RecordRule

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
		target := model.RecordRule{ID: params.ID}

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
	r.DELETE("/record/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.RecordRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.RecordRule{ID: params.ID}

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
	r.GET("/record/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.RecordRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.RecordRule{ID: params.ID}

		// get first one
		if params, err = target.GetFirst(); err != nil {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_not_exists"]))
			return
		}

		Success(c, params)
	})

	// get list
	r.GET("/record/rules", func(c *gin.Context) {
		var err error
		var params model.RecordRule

		// bind params (form params)
		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		list, err := params.GetList("name")
		Success(c, list)
	})

	// rule flush
	r.GET("/record/flush", func(c *gin.Context) {
		err := FlushRecord()
		if ErrorIf(c, err) {
			return
		}
		Success(c, "ok")
	})
}

// FlushRecord flush record rule
func FlushRecord() error {
	var err error

	// get all alert list
	RecordRules, _ := (&model.RecordRule{}).GetList()

	// set interval and static type
	recordInterval := []string{"5s", "1m", "5m", "1h"}
	recordStaticType := []string{"min", "max", "avg"}

	// generate rule for each interval
	for seq, interval := range recordInterval {

		// set source(temporary) rule file
		ruleFileName := common.RecStatName + "_" + fmt.Sprintf("%02d", seq) + ".rule.yml"
		ruleFileWorkPath := fmt.Sprintf("%s/%s", common.PromWorkPath, ruleFileName) // source

		// create source(temporary) rule file
		var file *os.File
		if file, err = os.Create(ruleFileWorkPath); err != nil {
			return err
		}
		defer file.Close()
		common.Log.Info("temporary rule file ", ruleFileWorkPath, " created")

		// generate rule (header)
		ruleString := fmt.Sprintln("groups:")
		ruleString += fmt.Sprintln("- name: " + common.RecStatName + "_" + interval + "_rules")
		ruleString += fmt.Sprintln("  interval: " + interval)
		ruleString += fmt.Sprintln("  rules:")

		// generate rule (body)
		for _, rule := range RecordRules {

			if seq == 0 {
				// raw record
				ruleString += fmt.Sprintln("  - record: " + common.RecRawName + ":" + rule.Name)
				ruleString += fmt.Sprintln("    expr:  ")
				ruleString += fmt.Sprintln("      " + rule.Query)

			} else if rule.StatYn == "Y" {
				// static rules
				for _, opr := range recordStaticType {

					var metricName string
					if seq > 2 {
						metricName = common.RecStatName + ":" + rule.Name + ":" + opr + ":" + recordInterval[seq-2]
					} else {
						metricName = common.RecRawName + ":" + rule.Name
					}

					ruleString += fmt.Sprintln("  - record: " + common.RecStatName + ":" + rule.Name + ":" + opr + ":" + interval)
					ruleString += fmt.Sprintln("    expr:  ")
					ruleString += fmt.Sprintln("      " + opr + "_over_time(" + metricName + "[" + interval + "])")
				}
			} else {
				common.Log.Info("Skip static rule -", rule.Name)
				ruleString += fmt.Sprintln("  - record: " + common.RecStatName + ":" + rule.Name + ":snap:" + interval)
				ruleString += fmt.Sprintln("    expr:  ")
				ruleString += fmt.Sprintln("      " + common.RecRawName + ":" + rule.Name)
			}
		}

		// write to temporary rule file
		if _, err = fmt.Fprintln(file, ruleString); err != nil {
			return err
		}
		common.Log.Info("temporary rule file ", ruleFileWorkPath, " writed")

		// Check Rules
		if _, err = exec.Command(common.Prom.Promtool, "check", "rules", ruleFileWorkPath).Output(); err != nil {
			return err
		}
		common.Log.Info(common.Prom.Promtool, "check", "rules", ruleFileWorkPath)

	}

	// move to prometheus rule directory
	errored := false
	for seq := range recordInterval {

		// set source(temporary) and target(prometheus) rule file
		ruleFileName := common.RecStatName + "_" + fmt.Sprintf("%02d", seq) + ".rule.yml"
		ruleFileWorkPath := fmt.Sprintf("%s/%s", common.PromWorkPath, ruleFileName)  // source
		ruleFilePromPath := fmt.Sprintf("%s/%s", common.Prom.RulePath, ruleFileName) // destination

		// move to rule file directory
		if err = os.Rename(ruleFileWorkPath, ruleFilePromPath); err != nil {
			common.Log.Error("Rule file move failed", ruleFileWorkPath, "->", ruleFilePromPath)
			errored = true
		} else {
			common.Log.Info("Rule file moved", ruleFileWorkPath, "->", ruleFilePromPath)
		}
	}

	// reload if no error
	if errored {
		return fmt.Errorf("move file failed")
	}

	// reload rule
	if err = common.Prom.Reload(); err != nil {
		return err
	}

	return err
}
