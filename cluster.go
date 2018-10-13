/*
https://www.elastic.co/guide/en/elasticsearch/guide/1.x/_rolling_restarts.html
Please note the rolling restart related API is only supported till ES 2.0.
*/
package main

func disableShardAllocation(esServer string) {
	payload := `{
    "transient" : {
        "cluster.routing.allocation.enable" : "none"
      }
	}`
	request := RestApiRequest{method: "put", api: "_cluster", pathinfo: "settings", payload: payload}
	doRestApi(esServer, &request)
}

func enableShardAllocation(esServer string) {
	payload := `{
    "transient" : {
        "cluster.routing.allocation.enable" : "all"
      }
	}`
	request := RestApiRequest{method: "put", api: "_cluster", pathinfo: "settings", payload: payload}
	doRestApi(esServer, &request)
}

func shutdownNode(esServer string) {

	pathinfo := "nodes/_local/_shutdown"
	request := RestApiRequest{method: "post", api: "_cluster", pathinfo: pathinfo}
	doRestApi(esServer, &request)
}

func clusterSettings(esServer string) {
	pathinfo := "settings"
	request := RestApiRequest{method: "get", api: "_cluster", pathinfo: pathinfo}
	doRestApi(esServer, &request)
}

func clusterState(esServer string) {
	pathinfo := "state"
	request := RestApiRequest{method: "get", api: "_cluster", pathinfo: pathinfo}
	doRestApi(esServer, &request)
}

func nodes(esServer string) {
	pathinfo := "nodes?v"
	request := RestApiRequest{method: "get", api: "_cat", pathinfo: pathinfo}
	doRestApi(esServer, &request)
}
