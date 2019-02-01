#!/bin/bash 


exec /sample-raft/node -ca ${CLUSTER_ADDRS} $@
