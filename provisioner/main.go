package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	id              = "yrfs"
	driver          = "yr/yrfs"
	execCommandPath = flag.String("execCommandPath", "/usr/libexec/kubernetes/kubelet-plugins/volume/exec", "The provisioner executable path.")
	timeout         = flag.Duration("timeout", 10000, "Command execute timeout")
)

// Main func
func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	if execCommandPath == nil {
		glog.Fatalf("Invalid flags specified: must provide provisioner exec command")
	}

	// Create an InClusterConfig and use it to create a client for the controller
	// to use to communicate with Kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatalf("Failed to create config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}
	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("Error getting server version: %v", err)
	}
	provisioner := NewYrfsProvisioner(clientset, id, driver, (*execCommandPath)+id, *timeout)
	// Start the provision controller which will dynamically
	// provision awesome volume PVs
	pc := controller.NewProvisionController(
		clientset,
		"yrfs",
		provisioner,
		serverVersion.GitVersion,
	)
	pc.Run(wait.NeverStop)
}

