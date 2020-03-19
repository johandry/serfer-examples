#!/bin/sh

hostname=$(hostname)
config=""

config=${config}"{"
config=${config}"  \"node_name\": \"${hostname}\","
# config=${config}"  \"tags\": {"
# config=${config}"    \"role\": \"serfer\""
# config=${config}"  },"
config=${config}"  \"bind\": \"${hostname}\","
# config=${config}"  \"keyring_file\": \"keys.json\","
config=${config}"  \"rpc_addr\": \"0.0.0.0:7373\","
config=${config}"  \"event_handlers\": ["
config=${config}"    \"/data/logger.sh\""
config=${config}"  ],"
config=${config}"  \"interface\": \"eth0\""

# If this is not the leader, join to the leader
if [[ "${hostname}" != "serf_leader.docker.local" ]]; then
config=${config}","
config=${config}"  \"retry_join\": ["
config=${config}"    \"serf_leader\""
config=${config}"  ],"
config=${config}"  \"retry_interval\": \"5s\""
fi

config=${config}"}"

mkdir -p /etc/serf
echo ${config} > /etc/serf/conf.json

echo "Serf configuration: '$(cat /serf.json)'"

/usr/local/bin/serf agent -config-file /etc/serf/conf.json