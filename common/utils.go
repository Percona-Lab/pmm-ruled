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
	"crypto/md5"
	"fmt"
	"strconv"
)

// PanicIf panic if error
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseInt string -> int
func ParseInt(value string) int {
	if value == "" {
		return 0
	}
	val, _ := strconv.Atoi(value)
	return val
}

// IntString int -> string
func IntString(value int) string {
	return strconv.Itoa(value)
}

// MD5 get md5
func MD5(value string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(value)))
}
