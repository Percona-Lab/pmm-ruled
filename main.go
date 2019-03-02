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

package main

import (
	"pmm-ruled/batch"
	"pmm-ruled/common"
	"pmm-ruled/exporter"
	"pmm-ruled/handler"
	"pmm-ruled/model"

	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func init() {
	common.LoadConfig()
	model.NewDatabase()
}

func main() {
	// =======================
	// Admin server
	// =======================
	g.Go(func() error {
		return handler.StartAdmin()
	})

	// // =======================
	// // Snapshotshot gather batch
	// // =======================
	g.Go(func() error {
		return batch.StartSnapshotBatch()
	})

	// =======================
	// Sync Instance Batch
	// =======================
	g.Go(func() error {
		return batch.StartInstanceBatch()
	})

	// // ================================
	// // Start Threshold Exporter
	// // ================================
	g.Go(func() error {
		return exporter.StartExporter()
	})

	if err := g.Wait(); err != nil {
		common.Log.Fatal(err)
	}
}
