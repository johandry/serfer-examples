#!/bin/sh

echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$(hostname)] Event: ${SERF_EVENT} - ${SERF_SELF_NAME} - ${SERF_USER_EVENT} - ${SERF_QUERY_NAME}" >> /data/serf_event.log