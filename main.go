/*
Copyright 2021.

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
	"context"
	"flag"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	helmv1 "github.com/k3s-io/helm-controller/pkg/apis/helm.cattle.io/v1"
	upgraderesponder "github.com/longhorn/upgrade-responder/client"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	demov1alpha1 "github.com/rancher/opni/apis/demo/v1alpha1"
	opniiov1beta1 "github.com/rancher/opni/apis/v1beta1"
	"github.com/rancher/opni/controllers"
	"github.com/rancher/opni/controllers/demo"
	"github.com/rancher/opni/pkg/metrics"
	"github.com/rancher/opni/pkg/util/manager"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
	Version  = "dev"
)

const (
	upgradeResponderAddress = "https://opni-usage.danbason.dev/v1/checkupgrade"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opniiov1beta1.AddToScheme(scheme))
	utilruntime.Must(demov1alpha1.AddToScheme(scheme))
	utilruntime.Must(helmv1.AddToScheme(scheme))
	utilruntime.Must(apiextv1beta1.AddToScheme(scheme))
	utilruntime.Must(loggingv1beta1.AddToScheme(scheme))
	utilruntime.Must(monitoringv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var disableUsage bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&disableUsage, "disable-usage", false, "Disable anonymous Opni usage tracking.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "98e737d4.opni.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if !disableUsage && Version != "dev" {
		upgradeRequester := manager.UpgradeRequester{Version: Version}
		upgradeRequester.SetupLoggerWithManager(mgr)
		setupLog.Info("Usage tracking enabled", "current-version", Version)
		upgradeChecker := upgraderesponder.NewUpgradeChecker(upgradeResponderAddress, &upgradeRequester)
		upgradeChecker.Start()
		defer upgradeChecker.Stop()
	}

	if err = (&controllers.OpniClusterReconciler{}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpniCluster")
		os.Exit(1)
	}
	if err = demo.NewOpniDemoReconciler().SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpniDemo")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	mgr.Add(ctrlmgr.RunnableFunc(func(ctx context.Context) error {
		return wait.PollImmediate(5*time.Second, 30*time.Second, func() (bool, error) {
			return metrics.ReconcileServiceMonitor(ctx, mgr, setupLog)
		})
	}))

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
