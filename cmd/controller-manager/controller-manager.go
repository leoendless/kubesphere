/*

 Copyright 2019 The KubeSphere Authors.

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
	"flag"
	"k8s.io/klog"
	"kubesphere.io/kubesphere/cmd/controller-manager/app"
	"kubesphere.io/kubesphere/pkg/apis"
	"kubesphere.io/kubesphere/pkg/controller"
	"kubesphere.io/kubesphere/pkg/simple/client/k8s"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var (
	masterURL   string
	metricsAddr string
)

func init() {
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
}

func main() {
	flag.Parse()

	cfg, err := k8s.Config()
	if err != nil {
		klog.Error(err, "failed to build kubeconfig")
		os.Exit(1)
	}

	stopCh := signals.SetupSignalHandler()

	klog.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		klog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	klog.Info("setting up scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	klog.Info("Setting up controllers")
	if err := controller.AddToManager(mgr); err != nil {
		klog.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	if err := app.AddControllers(mgr, cfg, stopCh); err != nil {
		klog.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	klog.Info("Starting the Cmd.")
	if err := mgr.Start(stopCh); err != nil {
		klog.Error(err, "unable to run the manager")
		os.Exit(1)
	}

}
