/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"log"
	"os"
	"time"

	//
	// Uncomment to load all auth plugins
	//
	// Or uncomment to load specific auth plugins

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func main() {
	// // creates the in-cluster config
	// config, err := rest.InClusterConfig()

	// if err != nil {
	// 	panic(err.Error())
	// }
	// // creates the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// pods, err := clientset.CoreV1().Pods("poker").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// _, err = clientset.CoreV1().Pods("poker").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
	// if errors.IsNotFound(err) {
	// 	fmt.Printf("Pod example-xxxxx not found in poker namespace\n")
	// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	// } else if err != nil {
	// 	panic(err.Error())
	// } else {
	// 	fmt.Printf("Found example-xxxxx pod in poker namespace\n")
	// }

	settings := cli.New()

	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	clientRead := action.NewList(actionConfig)

	clientRead.Deployed = true
	results, err := clientRead.Run()
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	if len(results) > 0 {
		for _, room := range results[0].Config["rooms"].([]interface{}) {
			if room.(string) == "TeamRed" {
				log.Printf("found %+v value", room)

			}
		}

		log.Printf("helm results: %+v", results[0].Config["rooms"])
	}
	modifyRooms(results[0])

	time.Sleep(10 * time.Second)
}

func modifyRooms(currentRelease *release.Release) {
	settings := cli.New()
	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	cfg.KubeClient.IsReachable()
	client := action.NewUpgrade(cfg)

	client.Namespace = settings.Namespace()
	name := "poker"
	if err := chartutil.ValidateReleaseName(name); err != nil {
		log.Printf("release name is invalid: %s", name)
		os.Exit(1)
	}
	ctx := context.Background()

	newValues := make(map[string]interface{})
	for k, v := range currentRelease.Config {
		newValues[k] = v
		// replace rooms
		if k == "rooms" {
			newValues[k] = []string{
				"John",
				"Paul",
				"George",
				"Ringo"}
		}
	}
	// newValues["rooms"] = append(newValues["rooms"].([]string), "OOOoooOOO")
	// var values = map[string]interface{}{"rooms": []string{
	// 	"John",
	// 	"Paul",
	// 	"George",
	// 	"Ringo"},
	// }
	if _, err := client.RunWithContext(ctx, name, currentRelease.Chart, newValues); err != nil {
		log.Printf("run upgrade gives error: %s", err)
		os.Exit(1)
	}
}
