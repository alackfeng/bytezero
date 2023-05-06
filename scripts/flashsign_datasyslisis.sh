#!/bin/bash

## 每天记录.
#### /opt/bytezero/bin/bytezero flashsign --last-report-date "2022-12-05 00:00:00"
## 4月度记录.
#### /opt/bytezero/bin/bytezero flashsign --table-field revenueMonth --last-report-date "2023-04-01 00:00:00"

nohup /opt/bytezero/bin/bytezero flashsign >> /opt/bytezero/bz.log &
