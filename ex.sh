clickhouse-client --host 79.175.157.227 \
  --port 9000 \
  --user carriot \
  --password ZwFjX39FFSDYrecjEJdxAruQ \
  --format CSVWithNames \
  --query "select * from car_devicelog_new where device_id in ('860697059573628') order by datetime limit 10" \
  >> devices.csv