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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// PromMetric prometheus metric result struct
type PromMetric struct {
	Status    string
	ErrorType string
	Error     string
	Data      struct {
		ResultType string
		Result     []struct {
			Metric map[string]string
			Value  []interface{}
		}
	}
}

// PromAPI prometheus http client
type PromAPI struct {
	API      string
	Timeout  int
	RulePath string
	Promtool string
	client   *http.Client
}

// Prom prometheus api
var Prom PromAPI

// PromWorkPath prometheus temporary work path
var PromWorkPath string

// SetPrometheus new prometheus api
func SetPrometheus() {
	Prom = PromAPI{
		API:      ConfigStr["prom.api"],
		RulePath: ConfigStr["prom.rule_path"],
		Promtool: ConfigStr["prom.promtool"],
		Timeout:  ConfigInt["prom.timeout"],
	}

	PromWorkPath = ConfigStr["prom.work_path"]
	if err := os.Mkdir(PromWorkPath, os.ModePerm); !os.IsExist(err) {
		PanicIf(err)
	}
}

// Exec execute
func (o *PromAPI) Exec(s string) (PromMetric, error) {

	// Timeout setting
	o.client = &http.Client{
		Timeout: time.Duration(o.Timeout) * time.Millisecond,
	}

	var promMetric PromMetric
	api := fmt.Sprintf("%s/api/v1/query?query=%s", o.API, url.QueryEscape(s))
	resp, err := o.client.Get(api)
	if err != nil {
		return promMetric, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return promMetric, err
	}

	json.Unmarshal(body, &promMetric)
	if promMetric.Status != "success" {
		return promMetric, fmt.Errorf("prometheus query failed")
	}

	return promMetric, nil
}

// Reload rule and conf reload
func (o *PromAPI) Reload() error {
	resp, err := http.Post(fmt.Sprintf("%s/-/reload", o.API), "text/plain", bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

// PromRecordRule prometheus record rules
type PromRecordRule struct {
	Status string
	Data   struct {
		Groups []struct {
			Rules []struct {
				Name   string
				Query  string
				Type   string
				Health string
			}
		}
	}
}
