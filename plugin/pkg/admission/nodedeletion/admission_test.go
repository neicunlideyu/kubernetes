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

package nodedeletion

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"
	genericadmissioninitializer "k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"

	api "k8s.io/kubernetes/pkg/apis/core"
)

// newHandlerForTest returns the admission controller configured for testing.
func newHandlerForTest(c kubernetes.Interface) (admission.ValidationInterface, error) {
	handler := NewPlugin()
	pluginInitializer := genericadmissioninitializer.New(c, nil, nil, nil)
	pluginInitializer.Initialize(handler)
	err := admission.ValidateInitialization(handler)
	return handler, err
}

// newMockClientForTest creates a mock client that returns a client configured for the specified get of nodes.
func newMockClientForTest(node *api.Node, conditions []corev1.NodeCondition) *fake.Clientset {
	mockClient := &fake.Clientset{}
	mockClient.AddReactor("get", "nodes", func(action core.Action) (bool, runtime.Object, error) {
		return true, &corev1.Node{
			ObjectMeta: node.ObjectMeta,
			Status: corev1.NodeStatus{
				Conditions: conditions,
			},
		}, nil
	})
	return mockClient
}

// TestAdmissionNodeDeletion verifies node is allowed to be deleted only if status not ready.
func TestAdmissionNodeDeletion(t *testing.T) {
	var (
		resource = api.Resource("nodes").WithVersion("v1")

		myNodeObjMeta     = metav1.ObjectMeta{Name: "mynode"}
		myNodeObj         = api.Node{ObjectMeta: myNodeObjMeta}
		readyCondition    = corev1.NodeCondition{Type: corev1.NodeReady, Status: corev1.ConditionTrue}
		notReadyCondition = corev1.NodeCondition{Type: corev1.NodeReady, Status: corev1.ConditionFalse}
		nodeKind          = api.Kind("Node").WithVersion("v1")
	)

	tests := []struct {
		name          string
		node          api.Node
		oldNode       api.Node
		operation     admission.Operation
		conditions    []corev1.NodeCondition
		expectSuccess bool
	}{
		{
			name:          "Node with no status can be deleted",
			node:          myNodeObj,
			operation:     admission.Delete,
			conditions:    []corev1.NodeCondition{},
			expectSuccess: true,
		},
		{
			name:          "Ready node can't be deleted",
			node:          myNodeObj,
			operation:     admission.Delete,
			conditions:    []corev1.NodeCondition{readyCondition},
			expectSuccess: false,
		},
		{
			name:          "NotReady node can be deleted",
			node:          myNodeObj,
			operation:     admission.Delete,
			conditions:    []corev1.NodeCondition{notReadyCondition},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockClientForTest(&tt.node, tt.conditions)
			handler, err := newHandlerForTest(mockClient)
			if err != nil {
				t.Errorf("unexpected error initializing handler: %v", err)
			}
			attributes := admission.NewAttributesRecord(&tt.node, &tt.oldNode, nodeKind, myNodeObj.Namespace, myNodeObj.Name, resource, "", tt.operation, &metav1.DeleteOptions{}, false, nil)
			err = handler.Validate(context.Background(), attributes, nil)
			if tt.expectSuccess && err != nil {
				t.Errorf("unexpected error returned from admission handler")
			} else if !tt.expectSuccess && err == nil {
				t.Errorf("expect error returned from admission handler")
			}
		})
	}
}
