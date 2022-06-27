package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"strconv"
	"os/exec"
)

type lvmCollector struct {
	vgSizeMetric *prometheus.Desc
}

// LVM Collector contains LV size and VG free in MB
func newLvmCollector() *lvmCollector {
	return &lvmCollector{
		vgSizeMetric: prometheus.NewDesc("lvm_lv_size",
			"Shows LVM LV size size in Bytes",
			[]string{"lv_name"}, nil,
		),
	}
}

func (collector *lvmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.vgSizeMetric
}

// LVM Collect, call OS command and set values
func (collector *lvmCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("lvs", "--separator", ",", "--units", "B", "--noheadings", "-o", "lv_name,lv_size").Output()
	if err != nil {
		log.Print(err)
	}
	lines := strings.Split(string(out),"\n")
	for _, line := range lines {
		values := strings.Split(line,",")
		if len(values)==2 {
			lv_size, err := strconv.ParseFloat(strings.Trim(values[1],"B"), 64)
			if err!= nil {
				log.Print(err)
			} else {
				lv_name := strings.Trim(values[0], " ")
				ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.GaugeValue, lv_size, lv_name)
			}
		}
	}

}
