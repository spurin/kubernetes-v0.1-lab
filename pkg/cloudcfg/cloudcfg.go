/*
Copyright 2014 Google Inc. All rights reserved.

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

// Package cloudcfg is ...
package cloudcfg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"gopkg.in/yaml.v1"
)

func promptForString(field string) string {
	fmt.Printf("Please enter %s: ", field)
	var result string
	fmt.Scan(&result)
	return result
}

// Parse an AuthInfo object from a file path
func LoadAuthInfo(path string) (client.AuthInfo, error) {
	var auth client.AuthInfo
	if _, err := os.Stat(path); os.IsNotExist(err) {
		auth.User = promptForString("Username")
		auth.Password = promptForString("Password")
		data, err := json.Marshal(auth)
		if err != nil {
			return auth, err
		}
		err = ioutil.WriteFile(path, data, 0600)
		return auth, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return auth, err
	}
	err = json.Unmarshal(data, &auth)
	return auth, err
}

// Perform a rolling update of a collection of tasks.
// 'name' points to a replication controller.
// 'client' is used for updating tasks.
// 'updatePeriod' is the time between task updates.
func Update(name string, client client.ClientInterface, updatePeriod time.Duration) error {
	controller, err := client.GetReplicationController(name)
	if err != nil {
		return err
	}
	labels := controller.DesiredState.ReplicasInSet

	taskList, err := client.ListTasks(labels)
	if err != nil {
		return err
	}
	for _, task := range taskList.Items {
		_, err = client.UpdateTask(task)
		if err != nil {
			return err
		}
		time.Sleep(updatePeriod)
	}
	return nil
}

// RequestWithBody is a helper method that creates an HTTP request with the specified url, method
// and a body read from 'configFile'
// FIXME: need to be public API?
func RequestWithBody(configFile, url, method string) (*http.Request, error) {
	if len(configFile) == 0 {
		return nil, fmt.Errorf("empty config file.")
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	return RequestWithBodyData(data, url, method)
}

// RequestWithBodyData is a helper method that creates an HTTP request with the specified url, method
// and body data
// FIXME: need to be public API?
func RequestWithBodyData(data []byte, url, method string) (*http.Request, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	request.ContentLength = int64(len(data))
	return request, err
}

// Execute a request, adds authentication, and HTTPS cert ignoring.
// TODO: Make this stuff optional
// FIXME: need to be public API?
func DoRequest(request *http.Request, user, password string) (string, error) {
	request.SetBasicAuth(user, password)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return string(body), err
}

// StopController stops a controller named 'name' by setting replicas to zero
func StopController(name string, client client.ClientInterface) error {
	controller, err := client.GetReplicationController(name)
	if err != nil {
		return err
	}
	controller.DesiredState.Replicas = 0
	controllerOut, err := client.UpdateReplicationController(controller)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(controllerOut)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

func makePorts(spec string) []api.Port {
	parts := strings.Split(spec, ",")
	var result []api.Port
	for _, part := range parts {
		pieces := strings.Split(part, ":")
		if len(pieces) != 2 {
			log.Printf("Bad port spec: %s", part)
			continue
		}
		host, err := strconv.Atoi(pieces[0])
		if err != nil {
			log.Printf("Host part is not integer: %s %v", pieces[0], err)
			continue
		}
		container, err := strconv.Atoi(pieces[1])
		if err != nil {
			log.Printf("Container part is not integer: %s %v", pieces[1], err)
			continue
		}
		result = append(result, api.Port{ContainerPort: container, HostPort: host})
	}
	return result
}

// RunController creates a new replication controller named 'name' which creates 'replicas' tasks running 'image'
func RunController(image, name string, replicas int, client client.ClientInterface, portSpec string, servicePort int) error {
	controller := api.ReplicationController{
		JSONBase: api.JSONBase{
			ID: name,
		},
		DesiredState: api.ReplicationControllerState{
			Replicas: replicas,
			ReplicasInSet: map[string]string{
				"name": name,
			},
			TaskTemplate: api.TaskTemplate{
				DesiredState: api.TaskState{
					Manifest: api.ContainerManifest{
						Containers: []api.Container{
							api.Container{
								Image: image,
								Ports: makePorts(portSpec),
							},
						},
					},
				},
				Labels: map[string]string{
					"name": name,
				},
			},
		},
		Labels: map[string]string{
			"name": name,
		},
	}

	controllerOut, err := client.CreateReplicationController(controller)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(controllerOut)
	if err != nil {
		return err
	}
	fmt.Print(string(data))

	if servicePort > 0 {
		svc, err := createService(name, servicePort, client)
		if err != nil {
			return err
		}
		data, err = yaml.Marshal(svc)
		if err != nil {
			return err
		}
		fmt.Printf(string(data))
	}
	return nil
}

func createService(name string, port int, client client.ClientInterface) (api.Service, error) {
	svc := api.Service{
		JSONBase: api.JSONBase{ID: name},
		Port:     port,
		Labels: map[string]string{
			"name": name,
		},
	}
	svc, err := client.CreateService(svc)
	return svc, err
}

// DeleteController deletes a replication controller named 'name', requires that the controller
// already be stopped
func DeleteController(name string, client client.ClientInterface) error {
	controller, err := client.GetReplicationController(name)
	if err != nil {
		return err
	}
	if controller.DesiredState.Replicas != 0 {
		return fmt.Errorf("controller has non-zero replicas (%d)", controller.DesiredState.Replicas)
	}
	return client.DeleteReplicationController(name)
}
