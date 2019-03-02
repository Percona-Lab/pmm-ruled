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

package exporter

import (
	"net/http"
	"pmm-ruled/common"
	"pmm-ruled/model"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/version"
)

// Metric name parts.
const (
	namespace = "ruled"
	exporter  = "collector"
)

// Metric descriptors.
var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
	scrapeThresholdListDesc = prometheus.NewDesc(
		prometheus.BuildFQName("alert", "rule", "threshold"),
		"Threshold values for each instance and level",
		[]string{"instance", "level", "name", "group"}, nil)

	scrapeActivateListDesc = prometheus.NewDesc(
		prometheus.BuildFQName("alert", "rule", "activate"),
		"Alert activation",
		[]string{"instance", "level", "name", "group"}, nil)
)

// RuleExporter exporter struct
type RuleExporter struct {
	error        prometheus.Gauge
	totalScrapes prometheus.Counter
	scrapeErrors *prometheus.CounterVec
}

// Describe prometheus describe
func (e *RuleExporter) Describe(ch chan<- *prometheus.Desc) {

	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect prometheus collect
func (e *RuleExporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)

	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
}

// scrape monitory thresholds metric collect
func (e *RuleExporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	scrapeTime := time.Now()

	rows := (&model.AlertRule{}).GetAlertThresoldList()
	for _, thr := range rows {
		// Thresold values
		val, _ := strconv.ParseFloat(thr.Val, 64)
		ch <- prometheus.MustNewConstMetric(scrapeThresholdListDesc, prometheus.GaugeValue, val, thr.InstanceName, thr.Level, thr.RuleName, thr.GroupName)

		// Activate values
		ch <- prometheus.MustNewConstMetric(scrapeActivateListDesc, prometheus.GaugeValue, float64(thr.Activate), thr.InstanceName, thr.Level, thr.RuleName, thr.GroupName)
	}

	// Scrap time
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "connection")
}

// StartExporter start exporter server
func StartExporter() error {
	// Define namespace
	prometheus.MustRegister(version.NewCollector(namespace + "_" + exporter))

	// Rule Thresholds
	metricPathRule := "/metric-rule"
	http.HandleFunc(metricPathRule, prometheus.InstrumentHandlerFunc(metricPathRule, func(w http.ResponseWriter, r *http.Request) {

		registry := prometheus.NewRegistry()
		registry.MustRegister(&RuleExporter{
			totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: exporter,
				Name:      "scrapes_total",
				Help:      "Total number of times Alarm Schedule was scraped for metrics.",
			}),
			scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: exporter,
				Name:      "scrape_errors_total",
				Help:      "Total number of times an error occurred scraping a Alarm Schedule.",
			}, []string{"collector"}),
			error: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: exporter,
				Name:      "last_scrape_error",
				Help:      "Whether the last scrape of metrics from Alarm Schedule resulted in an error (1 for error, 0 for success).",
			}),
		})

		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}

		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)

	}))

	common.Log.Info("Start exporter, listening", common.ConfigStr["glob.exp_listen_port"])
	return http.ListenAndServe(common.ConfigStr["glob.exp_listen_port"], nil)
}
