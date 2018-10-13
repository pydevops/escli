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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	logrus "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

const (
	GetMethod    string = "get"
	PostMethod   string = "post"
	PutMethod    string = "put"
	DeleteMethod string = "delete"
)

type RestApiRequest struct {
	method   string
	api      string
	pathinfo string
	payload  string
}

// func printBody(resp gorequest.Response, body string, errs []error) {
// 	fmt.Println(resp.Status)
// 	fmt.Println(body)
// }

func doRestApi(esServer string, apiRequest *RestApiRequest) {

	es_url := fmt.Sprintf("http://%s:%d/%s/%s", esServer, defaultESPort, apiRequest.api, apiRequest.pathinfo)
	logrus.Info(es_url)
	request := gorequest.New()

	var resp gorequest.Response
	var body string
	var errs []error

	switch apiRequest.method {
	case GetMethod:
		resp, body, errs = request.Get(es_url).End()
	case PostMethod:
		resp, body, errs = request.Post(es_url).End()
	case PutMethod:
		if apiRequest.payload == "" {
			resp, body, errs = request.Put(es_url).End()
		} else {
			resp, body, errs = request.Put(es_url).Send(apiRequest.payload).End()
		}
	case DeleteMethod:
		resp, body, errs = request.Delete(es_url).End()
	}
	if errs != nil {
		logrus.Error("Error attempting to doRestApi: ", errs)
	}

	// Non 2XX status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("Error creating snapshot [httpstatus: %d][url: %s] %s", resp.StatusCode, es_url, string(body))
		return
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, []byte(body), "", "  ")
	logrus.Println(buf)

}
