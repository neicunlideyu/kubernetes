/*
Copyright 2019 The Kubernetes Authors.

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

package nodedeletion

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/kubernetes"

	api "k8s.io/kubernetes/pkg/apis/core"
)

const (
	// PluginName is the name of the plugin.
	PluginName = "NodeDeletion"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewPlugin(), nil
	})
}

// NewPlugin creates a new NodeDeletion admission plugin.
// This plugin identifies node deletion requests
func NewPlugin() *Plugin {
	return &Plugin{
		Handler: admission.NewHandler(admission.Delete),
	}
}

// Plugin holds state for and implements the admission plugin.
type Plugin struct {
	*admission.Handler
	client kubernetes.Interface
}

var (
	_ admission.ValidationInterface = &Plugin{}
	_                               = genericadmissioninitializer.WantsExternalKubeClientSet(&Plugin{})
)

var (
	nodeResource = api.Resource("nodes")
)

// Validate makes sure that node status is not ready
func (p *Plugin) Validate(ctx context.Context, attributes admission.Attributes, o admission.ObjectInterfaces) (err error) {
	// Our job is just to check delete nodes.
	if attributes.GetResource().GroupResource() != nodeResource {
		return nil
	}
	node, err := p.client.CoreV1().Nodes().Get(ctx, attributes.GetName(), metav1.GetOptions{ResourceVersion: "0"})
	if err == nil {
		// check whether node status is ready
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				return admission.NewForbidden(attributes, fmt.Errorf("node %s is still ready, can't delete", node.Name))
			}
		}
	}
	return nil
}

// SetExternalKubeClientSet implements the WantsExternalKubeClientSet interface.
func (p *Plugin) SetExternalKubeClientSet(client kubernetes.Interface) {
	p.client = client
}

// ValidateInitialization implements the InitializationValidator interface.
func (p *Plugin) ValidateInitialization() error {
	if p.client == nil {
		return fmt.Errorf("missing client")
	}
	return nil
}
