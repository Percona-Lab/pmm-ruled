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

import "github.com/go-easylog/el"

// Logger login
type Logger struct{}

// Fatal fatal log
func (log *Logger) Fatal(vals ...interface{}) {
	el.Fatal(vals)
}

// Error error  log
func (log *Logger) Error(vals ...interface{}) {
	el.Error(vals)
}

// Warn warn log
func (log *Logger) Warn(vals ...interface{}) {
	el.Warn(vals)
}

// Info info log
func (log *Logger) Info(vals ...interface{}) {
	el.Info(vals)
}

// Trace trace log
func (log *Logger) Trace(vals ...interface{}) {
	el.Trace(vals)
}

// SetLogLevel change log level
func (log *Logger) SetLogLevel(level int) {
	el.SetLogLevel(el.LogLevel(level))
}
