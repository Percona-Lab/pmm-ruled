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
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unknwon/goconfig"
)

// ConfigStr configure string
var ConfigStr map[string]string

// ConfigInt configure int
var ConfigInt map[string]int

// Cfg config
var Cfg *goconfig.ConfigFile

// Log logger object
var Log Logger

// RecRawName Record raw metric prefix
var RecRawName string

// RecStatName Record stat metric prefix
var RecStatName string

// LoadConfig load config
func LoadConfig() {

	ConfigStr = make(map[string]string)
	ConfigInt = make(map[string]int)

	// parameter
	var config string
	flag.StringVar(&config, "config", "config.ini", "configuration")
	flag.Parse()

	var err error
	Cfg, err = goconfig.LoadConfigFile(config)
	if err != nil {
		panic("Load confguration failed")
	}

	// Load string configure
	ConfigStr["abs"], err = filepath.Abs(filepath.Dir(os.Args[0]))

	ConfigStr["glob.base"] = Cfg.MustValue("global", "base", "/pmm-ruled")
	ConfigStr["glob.exp_listen_port"] = Cfg.MustValue("global", "exp_listen_port", ":9104")
	ConfigStr["glob.adm_listen_port"] = Cfg.MustValue("global", "adm_listen_port", ":3333")
	ConfigInt["glob.log_level"] = Cfg.MustInt("global", "log_level", 2)

	ConfigInt["snapshot.interval"] = Cfg.MustInt("snapshot", "interval", 3)
	ConfigInt["snapshot.tombstone_sec"] = Cfg.MustInt("snapshot", "tombstone_sec", 600)
	ConfigStr["snapshot.row_key"] = Cfg.MustValue("snapshot", "row_key", "instance")

	ConfigStr["prom.api"] = Cfg.MustValue("prometheus", "api", "http://127.0.0.1:9090/prometheus")
	ConfigStr["prom.rule_path"] = Cfg.MustValue("prometheus", "rule_path", "prom-rule")
	ConfigStr["prom.work_path"] = fmt.Sprintf("%s/%s", ConfigStr["prom.rule_path"], "work")
	ConfigStr["prom.promtool"] = Cfg.MustValue("prometheus", "promtool", "promtool")
	ConfigInt["prom.timeout"] = Cfg.MustInt("prometheus", "timeout", 500)

	ConfigStr["db.host"] = Cfg.MustValue("database", "host", "127.0.0.1")
	ConfigStr["db.user"] = Cfg.MustValue("database", "user", "root")
	ConfigStr["db.pass"] = Cfg.MustValue("database", "pass", "pass")
	ConfigStr["db.db"] = Cfg.MustValue("database", "db", "db")
	ConfigInt["db.show_sql"] = Cfg.MustInt("database", "show_sql", 0)

	RecRawName = "raw"
	RecStatName = "stat"

	// Work path create
	if os.MkdirAll(ConfigStr["prom.work_path"], os.ModePerm); err != nil {
		PanicIf(err)
	}
	Log.Info("work path", ConfigStr["prom.work_path"], "ok")

	// Rule path create
	if os.MkdirAll(ConfigStr["prom.rule_path"], os.ModePerm); err != nil {
		PanicIf(err)
	}
	Log.Info("work path", ConfigStr["prom.rule_path"], "ok")

	// Log Setting
	Log.SetLogLevel(ConfigInt["glob.log_level"])

	// Prometheus http client setting
	SetPrometheus()

	// Load messages
	LoadMSG()
}
