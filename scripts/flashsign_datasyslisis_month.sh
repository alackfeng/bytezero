#!/bin/bash

## 每天记录.
#### /opt/bytezero/bin/bytezero flashsign --last-report-date "2022-12-05 00:00:00"
## 月度记录.
#### /opt/bytezero/bin/bytezero flashsign --table-field revenueMonth --last-report-date "2022-12-01 00:00:00" 

REPORT_DATE=`date "+%Y-%m-%d 00:00:00" -d '1 month ago'`
nohup /opt/bytezero/bin/bytezero flashsign --table-field revenueMonth --last-report-date "${REPORT_DATE}" >> /opt/bytezero/bz_month.log &
