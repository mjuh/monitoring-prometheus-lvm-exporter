# LVM thin pool Prometheus exporter

This program gets information about LVM thin pools with `lvs` binary and
provides LV's size and LV's free space as a Prometheus metrics.

The go binary will listen to port 9080 and serve metrics on the /metrics path.

http://example.org:9080/metrics
```
...
lvm_lv_data_percent{lv_name="p510348"} 6.96
lvm_lv_data_percent{lv_name="pool"} 86.37
...
lvm_lv_size{lv_name="p510348"} 6.442450944e+10
lvm_lv_size{lv_name="pool"} 3.790002454528e+12
...
```
