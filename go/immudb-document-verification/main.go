/*
Copyright 2023 Codenotary Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codenotary/immudb/pkg/api/authorizationschema"
	"github.com/codenotary/immudb/pkg/api/documentschema"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/verification"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const serverHost = "localhost"
const serverPort = 8080

var sessionID string

func main() {
	// open session
	openSessionResp := &authorizationschema.OpenSessionResponseV2{}

	err := doHttpRequest(
		http.MethodPost,
		"authorization/session/open",
		[]byte(`{
			"username": "immudb",
			"password": "immudb",
			"database": "defaultdb"
	  	}`),
		openSessionResp,
	)
	if err != nil {
		log.Fatal(err)
	}

	sessionID = openSessionResp.Token

	// create a collection
	err = doHttpRequest(
		http.MethodPut,
		"collections/create",
		[]byte(`{
			"name": "mycollection",
			"indexKeys": {
			  "attribute1": {
				"type": "STRING",
				"isUnique": true
			  },
			  "attribute2": {
				"type": "INTEGER",
				"isUnique": false
			  },
			  "attribute3": {
				"type": "DOUBLE",
				"isUnique": false
			  }
			}
		  }`),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// insert a document
	err = doHttpRequest(
		http.MethodPut,
		"documents/insert",
		[]byte(`{
			"collection": "mycollection",
			"document": {
				"attribute1": "doc1",
				"attribute2": 10,
				"attribute3": 4.2,
				"attribute4": true,
				"attribute5": "additional"
		  	}
		}`),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// fetch a document
	documentSearchResp := &documentschema.DocumentSearchResponse{}

	err = doHttpRequest(
		http.MethodPost,
		"documents/search",
		[]byte(`{
			"collection": "mycollection",
			"query": [],
			"page": 1,
			"perPage": 1
		}`),
		documentSearchResp,
	)
	if err != nil {
		log.Fatal(err)
	}

	var knownState *schema.ImmutableState

	for _, doc := range documentSearchResp.Results {
		docID := doc.Fields["_id"].GetStringValue()

		req := []byte(fmt.Sprintf(`{
			"collection": "mycollection",
			"documentId": "%s"
		}`, docID))

		// request the proof for the document
		proofResp := &documentschema.DocumentProofResponse{}

		err := doHttpRequest(
			http.MethodPost,
			"documents/proof",
			req,
			proofResp,
		)
		if err != nil {
			log.Fatal(err)
		}

		// validate proof
		knownState, err = verification.VerifyDocument(context.Background(), proofResp, doc, knownState, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("new known state at transaction: ", knownState.TxId)
}

func doHttpRequest(method string, url string, jsonBody []byte, resp protoreflect.ProtoMessage) error {
	requestURL := fmt.Sprintf("http://%s:%d/api/v2/%s", serverHost, serverPort, url)
	req, err := http.NewRequest(method, requestURL, bytes.NewReader(jsonBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("grpc-metadata-sessionid", sessionID)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return err
	}

	if resp == nil {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return protojson.Unmarshal(body, resp)
}
