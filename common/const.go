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

package common

// MSG message map
var MSG map[string]string

// LoadMSG load messages
func LoadMSG() {
	MSG = make(map[string]string)
	MSG["err.invalid_zero_id"] = "Invalid ID"
	MSG["err.rule_not_exists"] = "Rule not exists"
	MSG["err.group_not_exists"] = "Group not exists"
	MSG["err.instance_not_exists"] = "Instance not exists"
	MSG["err.get_metric_fail"] = "Get metric fail"
	MSG["err.val_not_digit"] = "Val must be digit"
	MSG["err.row_not_found"] = "Row not found"
	MSG["err.rule_exists"] = "Rule exists"
	MSG["err.invalid_operator"] = "Invalid operator"

	MSG["err.name_empty"] = "Name can not be empty"
	MSG["err.level_empty"] = "Level can not be empty"
	MSG["err.rule_empty"] = "Rule can not be empty"
	MSG["err.opr_empty"] = "Operator can not be empty"
	MSG["err.subj_empty"] = "Subject can not be empty"
	MSG["err.desc_empty"] = "Descrption can not be empty"
	MSG["err.val_empty"] = "Val can not be empty"
	MSG["err.query_empty"] = "Query can not be empty"
	MSG["err.statyn_empty"] = "Stat_YN can not be empty"
	MSG["err.invalid_statyn"] = "Stat must be 'Y' or 'N'"

	MSG["err.snapshot_bactch_fail"] = "Snapshotshot batch execution failed -"
}
