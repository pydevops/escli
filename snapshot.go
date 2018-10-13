/*
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/modules-snapshots.html

check indices recovery
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/indices-recovery.html
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/cat-recovery.html

repository-gcs
   https://www.elastic.co/guide/en/elasticsearch/plugins/master/repository-gcs-usage.html
*/
package main

import (
	"fmt"

	logrus "github.com/Sirupsen/logrus"
)

type SnapShot struct {
	backup   string
	snapshot string
}

func createS3Repo(esServer string, backupName string, bucketName string, region string) {
	payload := fmt.Sprintf(`{"type": "s3", "settings": 
		{ "bucket": "%s",
			"region": "%s"}}`, bucketName, region)
	snapshot := RestApiRequest{method: "put", api: "_snapshot", pathinfo: backupName, payload: payload}
	doRestApi(esServer, &snapshot)
}

func createGCSRepo(esServer string, backupName string, bucketName string, region string) {
	payload := fmt.Sprintf(`{"type": "gcs", "settings": 
		{ "bucket": "%s",
			"region": "%s"}}`, bucketName, region)
	logrus.Infof("%s", payload)
	snapshot := RestApiRequest{method: "put", api: "_snapshot", pathinfo: backupName, payload: payload}
	doRestApi(esServer, &snapshot)
}

func listRepos(esServer string) {
	snapshot := RestApiRequest{method: "get", api: "_snapshot"}
	doRestApi(esServer, &snapshot)
}

func snapshot(esServer string, ss *SnapShot) {
	pathinfo := fmt.Sprintf("%s/%s", ss.backup, ss.snapshot)
	snapshot := RestApiRequest{method: "put", api: "_snapshot", pathinfo: pathinfo}
	doRestApi(esServer, &snapshot)
}

func listSnapshots(esServer string, ss *SnapShot) {
	pathinfo := fmt.Sprintf("%s/%s", ss.backup, "_all")
	snapshot := RestApiRequest{method: "get", api: "_snapshot", pathinfo: pathinfo}
	doRestApi(esServer, &snapshot)
}

func getSnapshotStatus(esServer string, ss *SnapShot) {

	pathinfo := fmt.Sprintf("%s/%s/_status", ss.backup, ss.snapshot)
	snapshot := RestApiRequest{method: "get", api: "_snapshot", pathinfo: pathinfo}
	doRestApi(esServer, &snapshot)
}

func getIndicesRecovery(esServer string) {
	//snapshot := RestApiRequest{method: "get", api: "_recovery?human"}
	snapshot := RestApiRequest{method: "get", api: "_cat", pathinfo: "recovery?v"}
	doRestApi(esServer, &snapshot)
}

// close all indices, required by the restoreSnapshot
func closeAllIndices(esServer string) {

	pathinfo := fmt.Sprintf("_close")
	snapshot := RestApiRequest{method: "post", api: "_all", pathinfo: pathinfo}
	doRestApi(esServer, &snapshot)
}

// restore
func restoreSnapshot(esServer string, ss *SnapShot) {
	closeAllIndices(esServer)
	pathinfo := fmt.Sprintf("%s/%s/_restore", ss.backup, ss.snapshot)
	snapshot := RestApiRequest{method: "post", api: "_snapshot", pathinfo: pathinfo}
	doRestApi(esServer, &snapshot)
}
