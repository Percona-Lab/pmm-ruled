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
	"regexp"

	"github.com/gin-gonic/gin"
)

// startAlertRuleAPI alert rule API
func startAlertRuleAPI(r *gin.RouterGroup) {

	// new
	r.POST("/alert/rule", func(c *gin.Context) {
		var err error
		var params model.AlertRule

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
	r.PUT("/alert/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.AlertRule

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
		target := model.AlertRule{ID: params.ID}

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
	r.DELETE("/alert/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.AlertRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertRule{ID: params.ID}

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
	r.GET("/alert/rule/:rule_id", func(c *gin.Context) {
		var err error
		var params model.AlertRule

		// get id
		params.ID = common.ParseInt(c.Param("rule_id"))

		// check ID
		if params.ID == 0 {
			ErrorIf(c, fmt.Errorf(common.MSG["err.invalid_zero_id"]))
			return
		}

		// target
		target := model.AlertRule{ID: params.ID}

		// get first one
		if params, err = target.GetFirst(); err != nil {
			ErrorIf(c, fmt.Errorf(common.MSG["err.rule_not_exists"]))
			return
		}

		Success(c, params)

	})

	// get list
	r.GET("/alert/rules", func(c *gin.Context) {
		var err error
		var params model.AlertRule

		err = c.Bind(&params)
		if ErrorIf(c, err) {
			return
		}

		list, err := params.GetList("name", "level='critical' desc", "level='warning' desc", "level='info' desc")
		Success(c, list)
	})

	// rule flush
	r.GET("/alert/flush", func(c *gin.Context) {
		err := FlushAlert()
		if ErrorIf(c, err) {
			return
		}
		Success(c, "ok")
	})
}

// FlushAlert flush alert rule
func FlushAlert() error {
	var err error

	// get all alert list
	AlertRules, _ := (&model.AlertRule{}).GetList()

	// set source(temporary) and target(prometheus) rule file
	ruleFileName := "alert.rule.yml"
	ruleFileWorkPath := fmt.Sprintf("%s/%s", common.PromWorkPath, ruleFileName)  // source
	ruleFilePromPath := fmt.Sprintf("%s/%s", common.Prom.RulePath, ruleFileName) // destination

	// create source(temporary) rule file
	var file *os.File
	if file, err = os.Create(ruleFileWorkPath); err != nil {
		return err
	}
	defer file.Close()
	common.Log.Info("temporary rule file ", ruleFileWorkPath, " created")

	var ruleString string

	// generate rule (header)
	ruleString += fmt.Sprintln("groups:")
	ruleString += fmt.Sprintln("- name: alert_rules")
	ruleString += fmt.Sprintln("  rules:")

	// generate rule (body)
	for _, rule := range AlertRules {

		expr := fmt.Sprintf(`round(%s, 0.01) %s on (instance) group_left (level, name, gname) (alert_rule_threshold{name="%s", level="%s"} and (alert_rule_activate{name="%s", level="%s"} == 1))`, rule.Rule, rule.Opr, rule.Name, rule.Level, rule.Name, rule.Level)

		rule.Subject = regexp.MustCompile(`"`).ReplaceAllString(rule.Subject, `\"`)
		rule.Description = regexp.MustCompile(`"`).ReplaceAllString(rule.Description, `\"`)

		ruleString += fmt.Sprintln("  - alert: " + rule.Name)
		ruleString += fmt.Sprintln("    expr:  " + expr)
		if *rule.Wait > 0 {
			ruleString += fmt.Sprintf("    for: %ds\n", *rule.Wait)
		}
		ruleString += fmt.Sprintln(`    labels: `)
		ruleString += fmt.Sprintln(`      level: ` + rule.Level)
		ruleString += fmt.Sprintln(`    annotations: `)
		ruleString += fmt.Sprintln(`      summary: "` + rule.Subject + `"`)
		ruleString += fmt.Sprintln(`      description: "` + rule.Description + `" `)
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

	// move to prometheus rule directory
	if err = os.Rename(ruleFileWorkPath, ruleFilePromPath); err != nil {
		return err
	}
	common.Log.Info("Rule file moved", ruleFileWorkPath, "->", ruleFilePromPath)

	// reload rule
	if err = common.Prom.Reload(); err != nil {
		return err
	}

	return err
}
