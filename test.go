// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_helloworld_get]

// Package helloworld provides a set of Cloud Functions samples.
package ocspTest

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/crypto/ocsp"
)

// HelloGet is an HTTP Cloud Function.
func HelloGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

// [END functions_helloworld_get]

func getResponse(cert, issuer *x509.Certificate) ([]byte, error) {
	if len(cert.OCSPServer) == 0 {
		return nil, errors.New("no OCSPServer provided")
	}

	req, err := ocsp.CreateRequest(cert, issuer, nil)
	if err != nil {
		return nil, err
	}

	for _, server := range cert.OCSPServer {
		if u, err := url.Parse(server); err == nil {
			u.Path = base64.StdEncoding.EncodeToString(req)
			if resp, err := http.Get(u.String()); err == nil && resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if err == nil {
					return body, nil
				}
				log.Println(err)
			} else {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}
	return nil, errors.New("no OCSP response")
}
