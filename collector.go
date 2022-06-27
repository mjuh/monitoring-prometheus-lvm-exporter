package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type lvmCollector struct {
	lvSizeMetric        *prometheus.Desc
	lvDataPercentMetric *prometheus.Desc
}

// LVM Collector contains LV size and VG free in MB
func newLvmCollector() *lvmCollector {
	return &lvmCollector{
		lvSizeMetric: prometheus.NewDesc("lvm_lv_size",
			"Shows LVM LV size size in Bytes",
			[]string{"lv_name"}, nil,
		),
		lvDataPercentMetric: prometheus.NewDesc("lvm_lv_data_percent",
			"Shows LVM LV data percent",
			[]string{"lv_name"}, nil,
		),
	}
}

func (collector *lvmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.lvSizeMetric
	ch <- collector.lvDataPercentMetric
}

// LVM Collect, call OS command and set values
func (collector *lvmCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("lvs", "--separator", ",", "--units", "B", "--noheadings", "-o", "lv_name,lv_size,data_percent").Output()
	if err != nil {
		log.Print(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		values := strings.Split(line, ",")
		if len(values) == 3 {
			lv_size, err := strconv.ParseFloat(strings.Trim(values[1], "B"), 64)
			if err != nil {
				log.Print(err)
			} else {
				data_percent, err := strconv.ParseFloat(strings.Trim(values[2], "B"), 64)
				if err != nil {
					log.Print(err)
				} else {
					vg_name := strings.Trim(values[0], " ")
					ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.GaugeValue, lv_size, vg_name)
					ch <- prometheus.MustNewConstMetric(collector.lvDataPercentMetric, prometheus.GaugeValue, data_percent, vg_name)
				}
			}
		}
	}

}
