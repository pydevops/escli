#!/usr/bin/env bash

# https://www.elastic.co/guide/en/elasticsearch/reference/current/rolling-updates.html

# Stop immediately if something goes wrong
set -euo pipefail
#set -x

# disable_shard_allocation() - Sets the cluster.routing.allocation.enable
# setting to "none".  Prevents shards from being migrated from an upgrading
# Data Node to another active Data Node.
function disable_shard_allocation() {
  local SERVER=$1
  echo "Disable shard allocation..."
  curl -X PUT "$SERVER:9200/_cluster/settings" -H 'Content-Type: application/json' -d'
    {
    "persistent": {
        "cluster.routing.allocation.enable": "none"
    }
    }
   '
   echo ""
}

# enable_shard_allocation() - sets the cluster.routing.allocation.enable to the
# default value ("all")
function enable_shard_allocation() {
  local SERVER=$1
  echo "Re-enabling shard allocation"
  curl -X PUT "$SERVER:9200/_cluster/settings" -H 'Content-Type: application/json' -d'
  {
    "persistent": {
      "cluster.routing.allocation.enable": "all"
    }
   }
   '
   echo ""
}


# https://www.elastic.co/guide/en/elasticsearch/reference/5.6/allocation-filtering.html
# ecommission a node, and you would like to move the shards from that node to other nodes in the cluster before shutting it down.
function move_shards_off_node() {
  local NODE_IP=$1
  echo "moving shards off $NODE_IP"
  PAYLOAD="{\"transient\" : {\"cluster.routing.allocation.exclude._ip\" : \"$NODE_IP\"}}"
  curl -XPUT 'localhost:9200/_cluster/settings'  -H 'Content-Type: application/json' -d"$PAYLOAD"
}



function update_replicas() {
  local INDEX=$1
  local NUMBER_OF_REPLICAS=$2
  curl -XPUT "localhost:9200/$INDEX/_settings" -d"
  {
    \"index\" : {
        \"number_of_replicas\" : \"$NUMBER_OF_REPLICAS\"
    }
  }"
}

# define NODES (TODO)
# reroute replicas ONLY
function reroute_replica_to_node() {

  local INDEX=$1
  local NUMBER_OF_REPLICAS=4

  # disable the shard allocation first so that we can manually move the shards
  disable_shard_allocation

  # increase the # of replicas, the replicas becomes unassigned due to disabled shard allocation
  update_replicas $INDEX $NUMBER_OF_REPLICAS

  # iterate over NEW nodes (in GCP)
  # and iterator over each shard, move the unassigned replicas to the new nodes
  for node in "${NODES[@]}" 
    do
     for shard in {0..4} 
     do
      curl -XPOST "localhost:9200/_cluster/reroute" -H 'Content-Type: application/json' -d"
     {
        \"commands\" : [
            {
              \"allocate_replica\" : {
                    \"index\" : \"$INDEX\", \"shard\" : \"$shard\",
                    \"node\" : \"$node\"
              }
            }
        ]
    }
    "
   done
  done 
}


# Stop indexing and sync flush as recommended by the Elasticsearch documentation
# function stop_indexing() {
#   local SERVER=$1
#   echo "Stop non-essential indexing and perform a sync flush..."
#   curl -X POST "http://$SERVER:9200/_flush/synced"
#   echo ""
# }


# function init_shutdown() {
#  local SERVER=$1
#   disable_shard_allocation $SERVER
#   stop_indexing $SERVER
# }

# post_after_shutdown() - checks the cluster health endpoint and looks for a 'green'
# status response in a loop.
# post_after_shutdown() {
#   local SERVER=$1
#   echo "Checking cluster status"
#   # enable shard allocaiton if not
#   enable_shard_allocation $SERVER
#   # Wait for the shards to re-initialize and balanced
#   wait_for_allocations $SERVER
#   while true; do
#     STATUS=$(curl "http://$SERVER:9200/_cluster/health" 2>/dev/null \
#       | jq -r '.status')
#     if [[ "${STATUS}" == "green" ]]; then
#       echo "Cluster health is now ${STATUS}, continuing shutdown..."
#       disable_shard_allocation $SERVER
#       stop_indexing $SERVER
#       echo "****Please proceed to shut down the next node from the cluster****"
#       return 0
#     fi
#     echo "Cluster status: ${STATUS}"
#     sleep 5
#   done
# }



# wait_for_allocations() - Checks cluster health in a loop waiting for
# unassianged_shards to return to 0.
function wait_for_allocations() {
  local SERVER=$1
  echo "Checking shard allocations"
  while true; do
    UNASSIGNED=$(curl "http://$SERVER:9200/_cluster/health" 2>/dev/null \
      | jq -r '.unassigned_shards')
    if [[ "${UNASSIGNED}" == "0" ]]; then
      echo "All shards-reallocated"
      return 0
    else
      echo "Number of unassigned shards: ${UNASSIGNED}"
      sleep 3
    fi
  done
}


# wait for cluster green status
# the difference with `prep_for_shutdown`: it is read only
wait_for_green() {
 local SERVER=$1
 echo "Checking cluster status"
 # Wait for the shards to re-initialize and balanced
 wait_for_allocations $SERVER
 while true; do
    STATUS=$(curl "http://$SERVER:9200/_cluster/health" 2>/dev/null \
      | jq -r '.status')
    if [[ "${STATUS}" == "green" ]]; then
      echo "Cluster health is now ${STATUS}"
      echo "Cluster setting: ${STATUS}"
      curl "http://$SERVER:9200/_cluster/settings"
      return 0
    fi
   
    sleep 5
  done
}


# usage on the script
usage() {
    echo "$0 status|decomission|enable|disable <one_elasticserver_server>"
    echo "   <one_elasticserver_server> defaults to localhost"
    echo "   status: wait for cluster status as \"green\" "
    echo "   decomission: decommision a node by internal ip"
    echo "   enable-shard: enable shard allocation"
    echo "   disable-shard: disable shard allocation"
    echo "   update_replicas: update_replicas"
}

main () {

    if [[ "$#" -ge 1 ]]; then
      local CMD=${1}
      local SERVER=${2:-localhost}

      case "$CMD" in
        status)
          wait_for_green $SERVER
          ;;
        enable-shard)
          enable_shard_allocation $SERVER
          ;;
        disable-shard)
          disable_shard_allocation $SERVER
          ;;
        reroute-replica)
          reroute_replica_to_node bank
          ;;
        *)
          usage
          ;;
      esac
    else 
       usage
    fi

}

main "$@"