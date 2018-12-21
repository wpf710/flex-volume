package main

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/kubernetes-incubator/external-storage/lib/util"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/exec"
	"os"
	"time"
)

const (
	deleteCmd = "delete"
)

// Provisioner struct
type yrfsProvisioner struct {
	// kubernetes api client
	client kubernetes.Interface
	// Provisioner identity so you can
	// run multiple provisioners. Each
	// is responsible for separate sets of PVs
	id     string
	driver string

	execCommand string
	runner      exec.Interface
	timeout     time.Duration
}

var ErrorTimeout = fmt.Errorf("Timeout")

// Constructor
func NewYrfsProvisioner(c kubernetes.Interface, id string, driver string, execCommand string, timeout time.Duration) controller.Provisioner {
	return &yrfsProvisioner{client: c, id: id, driver: driver, execCommand: execCommand, timeout: timeout, runner: exec.New()}
}

// Calling other APIs to provision a persistent volume
// and return the volume ID
func createYrfsVolume(size int64) string {
	return "111111111-9-98-3"
}

// Calling other APIs to delete persistent
// volume by given id
func (p *yrfsProvisioner) deleteYrfsVolume(id string) error {
	var args = []string{deleteCmd, id}
	cmd := p.runner.Command(p.execCommand, args...)

	timeout := false
	if p.timeout > 0 {
		timer := time.AfterFunc(p.timeout, func() {
			timeout = true
			cmd.Stop()
		})
		defer timer.Stop()
	}

	output, execErr := cmd.CombinedOutput()
	if execErr != nil {
		if timeout {
			return ErrorTimeout
		}

		glog.Warningf("FlexVolume: driver call failed: executable: %s, args: %s, error: %s, output: %q", p.execCommand, args, execErr.Error(), output)
		return execErr
	}

	return nil
}

// Implementing external-storage controller required interface
func (p *yrfsProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	// Get volume size requirement from PVC
	capacity := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	volumeSize := util.RoundUpSize(capacity.Value(), util.GiB)
	// Provisioning the physical volume
	volumeId := createYrfsVolume(volumeSize)
	volumeId = options.PVName
	// This is the drive your created in "Part 1"

	// You can change to any file type as long as
	// your drive can deal with it
	fsType := "yrfs"
	// PV spec, similar to a PV manifest
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				"yrfsProvisionerIdentity": p.id,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexPersistentVolumeSource{
					Driver: p.driver,
					FSType: fsType,
					// Provide the name of the secret
					// if you are using one
					SecretRef: &v1.SecretReference{Name: os.Getenv("secret-name"),
						Namespace: options.PVC.Namespace,
					},
					ReadOnly: false,
					Options:  map[string]string{"volumeId": volumeId, "path": options.PVName, "storage": fmt.Sprintf("%d%s", volumeSize, "Gi")},
				},
			},
		},
	}
	return pv, nil
}

// Implementing external-storage controller required interface
func (p *yrfsProvisioner) Delete(volume *v1.PersistentVolume) error {
	// Check if requested PV should be processed by this provisioner
	ann, ok := volume.Annotations["yrfsProvisionerIdentity"]
	if !ok {
		return errors.New("identity annotation not found on PV")
	}
	if ann != p.id {
		return &controller.IgnoredError{"identity annotation on PV does not match"}
	}
	// Get volumeId from PV that created by this provisioner
	volumeId := volume.Spec.PersistentVolumeSource.FlexVolume.Options["volumeId"]
	glog.Infof("Received request to delete yrfs volume (ID: %s)\n", volumeId)
	// delete the physical volume
	if volumeId != "" {
		err := p.deleteYrfsVolume(volumeId)
		if err == nil {
			glog.Infof("Successfully deleted yrfs volume (ID: %s)\n", volumeId)
		}

		return err
	}

	return nil
}
