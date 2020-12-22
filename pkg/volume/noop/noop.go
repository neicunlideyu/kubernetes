/*
Copyright 2014 The Kubernetes Authors.

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

// Package noop contains the internal representation of noopDir
// volumes.
package noop

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/volume"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"
	"k8s.io/utils/mount"
	utilstrings "k8s.io/utils/strings"
)

const perm os.FileMode = 0777

// ProbeVolumePlugins is the primary entrypoint for volume plugins.
func ProbeVolumePlugins() []volume.VolumePlugin {
	return []volume.VolumePlugin{
		&noopPlugin{nil},
	}
}

type noopPlugin struct {
	host volume.VolumeHost
}

var _ volume.VolumePlugin = &noopPlugin{}

const (
	noopPluginName = "kubernetes.io/noop"

	wrappedVolumeFileName = "volume.json"
)

func getPath(uid types.UID, volName string, host volume.VolumeHost) string {
	return host.GetPodVolumeDir(uid, utilstrings.EscapeQualifiedName(noopPluginName), volName)
}

func (plugin *noopPlugin) Init(host volume.VolumeHost) error {
	plugin.host = host

	return nil
}

func (plugin *noopPlugin) GetPluginName() string {
	return noopPluginName
}

func (plugin *noopPlugin) GetVolumeName(spec *volume.Spec) (string, error) {
	// Return user defined volume name and rand string, since this is an noop volume type
	return fmt.Sprintf("%s-%s", spec.Name(), rand.String(8)), nil
}

func (plugin *noopPlugin) CanSupport(spec *volume.Spec) bool {
	if spec.PersistentVolume != nil && spec.PersistentVolume.Spec.CephFS != nil {
		if volumeutil.IsPVManagedByNoopPlugin(spec.PersistentVolume) {
			return true
		}
	}

	return false
}

func (plugin *noopPlugin) IsMigratedToCSI() bool {
	return false
}

func (plugin *noopPlugin) RequiresRemount() bool {
	return false
}

func (plugin *noopPlugin) SupportsMountOption() bool {
	return false
}

func (plugin *noopPlugin) SupportsBulkVolumeVerification() bool {
	return false
}

func (plugin *noopPlugin) NewMounter(spec *volume.Spec, pod *v1.Pod, opts volume.VolumeOptions) (volume.Mounter, error) {
	volumeWrappped := wrappedVolume{}

	if spec.PersistentVolume.Spec.CephFS != nil {
		cephData, err := plugin.constructCephFS(spec, pod)
		if err != nil {
			return nil, err
		}

		volumeWrappped.CephFS = cephData
	} else {
		return nil, fmt.Errorf("spec does not reference an CephFS persistentvolume volume type")
	}

	return plugin.newMounterInternal(spec, pod, plugin.host.GetMounter(plugin.GetPluginName()), &volumeWrappped)
}

func (plugin *noopPlugin) constructCephFS(spec *volume.Spec, pod *v1.Pod) (*cephfs, error) {
	cephfsPv := spec.PersistentVolume.Spec.CephFS
	cephData := &cephfs{
		Monitors:     cephfsPv.Monitors,
		Path:         cephfsPv.Path,
		User:         cephfsPv.User,
		ReadOnly:     spec.ReadOnly,
		MountOptions: volumeutil.MountOptionFromSpec(spec),
	}

	if cephData.User == "" {
		cephData.User = "admin"
	}

	if cephData.Path == "" {
		cephData.Path = "/"
	}

	if !strings.HasPrefix(cephData.Path, "/") {
		cephData.Path = "/" + cephData.Path
	}

	secret, err := plugin.getCephSecret(spec, pod.Namespace, cephfsPv.SecretFile)
	if err != nil {
		return nil, err
	}

	cephData.Secret = secret

	return cephData, nil
}

func (plugin *noopPlugin) newMounterInternal(spec *volume.Spec, pod *v1.Pod, mounter mount.Interface, volume *wrappedVolume) (volume.Mounter, error) {
	return &noopDir{
		pod:     pod,
		volName: spec.Name(),
		volume:  volume,
		mounter: mounter,
		plugin:  plugin,
	}, nil
}

func (plugin *noopPlugin) NewUnmounter(volName string, podUID types.UID) (volume.Unmounter, error) {
	// Inject real implementations here, test through the internal function.
	return plugin.newUnmounterInternal(volName, podUID, plugin.host.GetMounter(plugin.GetPluginName()))
}

func (plugin *noopPlugin) newUnmounterInternal(volName string, podUID types.UID, mounter mount.Interface) (volume.Unmounter, error) {
	ed := &noopDir{
		pod:     &v1.Pod{ObjectMeta: metav1.ObjectMeta{UID: podUID}},
		volName: volName,
		mounter: mounter,
		plugin:  plugin,
	}
	return ed, nil
}

func (plugin *noopPlugin) ConstructVolumeSpec(volName, mountPath string) (*volume.Spec, error) {
	fsMode := v1.PersistentVolumeFilesystem

	// TODO: load volume spec from volume.json
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeSource: v1.PersistentVolumeSource{
				CephFS: &v1.CephFSPersistentVolumeSource{
					Path: mountPath,
				},
			},
			VolumeMode: &fsMode,
		},
	}
	return volume.NewSpecFromPersistentVolume(pv, false), nil
}

func (plugin *noopPlugin) getCephSecret(spec *volume.Spec, defaultNamespace string, secretFile string) (string, error) {
	secretRef := spec.PersistentVolume.Spec.CephFS.SecretRef
	var (
		secretNs   = defaultNamespace
		secretName = ""
		secret     = ""
	)

	if secretRef != nil {
		if len(secretRef.Namespace) != 0 {
			secretNs = secretRef.Namespace
		}
		secretName = secretRef.Name

		// if secret is provideded, retrieve it
		kubeClient := plugin.host.GetKubeClient()
		if kubeClient == nil {
			return "", fmt.Errorf("cannot get kube client")
		}
		secrets, err := kubeClient.CoreV1().Secrets(secretNs).Get(context.Background(), secretName, metav1.GetOptions{})
		if err != nil {
			err = fmt.Errorf("couldn't get secret %v/%v err: %v", secretNs, secretName, err)
			return "", err
		}
		for name, data := range secrets.Data {
			secret = string(data)
			klog.V(4).Infof("found ceph secret info: %s", name)
		}
	} else if secretFile != "" {
		bytes, err := ioutil.ReadFile(secretFile)
		if err != nil {
			return "", fmt.Errorf("read secret file: %s failed: %+v", secretFile, err)
		}
		secret = string(bytes)
	}

	return secret, nil
}

// EmptyDir volumes are temporary directories exposed to the pod.
// These do not persist beyond the lifetime of a pod.
type noopDir struct {
	pod     *v1.Pod
	volName string
	volume  *wrappedVolume
	mounter mount.Interface
	plugin  *noopPlugin
	volume.MetricsNil
}

type wrappedVolume struct {
	CephFS *cephfs `json:"cephfs,omitempty"`
}

type cephfs struct {
	Monitors     []string `json:"monitors"`
	Path         string   `json:"path"`
	User         string   `json:"user"`
	ReadOnly     bool     `json:"readonly"`
	Secret       string   `json:"secret"`
	MountOptions []string `json:"mountOptions"`
}

func (ed *noopDir) GetAttributes() volume.Attributes {
	var readOnly bool
	if ed.volume.CephFS != nil {
		readOnly = ed.volume.CephFS.ReadOnly
	}

	return volume.Attributes{
		ReadOnly:        readOnly,
		Managed:         false,
		SupportsSELinux: false,
	}
}

// Checks prior to mount operations to verify that the required components (binaries, etc.)
// to mount the volume are available on the underlying node.
// If not, it returns an error
func (ed *noopDir) CanMount() error {
	return nil
}

// SetUp creates new directory.
func (ed *noopDir) SetUp(mounterArgs volume.MounterArgs) error {
	return ed.SetUpAt(ed.GetPath(), mounterArgs)
}

// SetUpAt creates new directory.
func (ed *noopDir) SetUpAt(dir string, mounterArgs volume.MounterArgs) error {
	_, err := ed.mounter.IsLikelyNotMountPoint(dir)
	// Getting an os.IsNotExist err from is a contingency; the directory
	// may not exist yet, in which case, setup should run.
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// If the plugin readiness file is present for this volume, and the
	// storage medium is the default, then the volume is ready.
	if volumeutil.IsReady(ed.getMetaDir()) {
		return nil
	}

	// Create the directory if it doesn't already exist.
	if err = os.MkdirAll(dir, perm); err != nil {
		return err
	}

	//Write wrapped volume metadata into volume.json because downstream, such as kata-runtime need it.
	if err = writeWrappedVolumeSpec(dir, ed.volume); err != nil {
		return err
	}

	volumeutil.SetReady(ed.getMetaDir())
	return nil
}

func (ed *noopDir) GetPath() string {
	return getPath(ed.pod.UID, ed.volName, ed.plugin.host)
}

// TearDown simply discards everything in the directory.
func (ed *noopDir) TearDown() error {
	return ed.TearDownAt(ed.GetPath())
}

// TearDownAt simply discards everything in the directory.
func (ed *noopDir) TearDownAt(dir string) error {
	if pathExists, pathErr := mount.PathExists(dir); pathErr != nil {
		return fmt.Errorf("error checking if path exists: %v", pathErr)
	} else if !pathExists {
		klog.Warningf("Warning: Unmount skipped because path does not exist: %v", dir)
		return nil
	}

	// Renaming the directory is not required anymore because the operation executor
	// now handles duplicate operations on the same volume
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func (ed *noopDir) getMetaDir() string {
	return path.Join(ed.plugin.host.GetPodPluginDir(ed.pod.UID, utilstrings.EscapeQualifiedName(noopPluginName)), ed.volName)
}

func writeWrappedVolumeSpec(dir string, volume *wrappedVolume) error {
	bytes, err := json.Marshal(volume)
	if err != nil {
		return err
	}
	volumeFile := path.Join(dir, wrappedVolumeFileName)
	klog.Infof("Write wrapped volume spec: %s to file: %s", string(bytes), volumeFile)
	return ioutil.WriteFile(volumeFile, bytes, perm)
}
