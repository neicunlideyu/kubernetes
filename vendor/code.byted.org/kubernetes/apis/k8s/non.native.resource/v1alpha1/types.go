/*
Copyright 2020 The Bytedance Authors.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RefinedNodeResource captures refined resource information about a node
// RefinedNodeResource objects are non-namespaced.
// Its name is set to node name
type RefinedNodeResource struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object metadata.
	// metadata.Name indicates the name of the refined node resource that this object
	// refers to; it MUST be the same name to Node object
	// call for that driver.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the refined node resource.
	Spec RefinedNodeResourceSpec `json:"spec"`

	// Status if the refined node resource
	Status RefinedNodeResourceStatus `json:"status"`
}

// RefinedNodeResourceSpec is the specification of a RefinedNodeResource.
type RefinedNodeResourceSpec struct {
	// discrete resources
	DiscreteResource DiscreteResourceProperties `json:"discreteResource,omitempty"`
}

// DiscreteResourceProperties is the set of discrete properties
type DiscreteResourceProperties struct {
	// TODO: combining the properties below together

	// refined CPU properties
	CPUProperties []DiscreteResourcePropertyPattern `json:"cpuProperties,omitempty"`

	// refined memory properties
	MemoryProperties []DiscreteResourcePropertyPattern `json:"memoryProperties,omitempty"`

	// refined GPU properties
	GPUProperties []DiscreteResourcePropertyPattern `json:"gpuProperties,omitempty"`

	// refined disk properties
	DiskProperties []DiscreteResourcePropertyPattern `json:"diskProperties,omitempty"`

	// refined network properties
	NetworkProperties []DiscreteResourcePropertyPattern `json:"networkProperties,omitempty"`

	// other refined discrete resources properties
	OtherProperties []DiscreteResourcePropertyPattern `json:"otherProperties,omitempty"`
}

type DiscreteResourcePropertyPattern struct {
	// property name
	PropertyName string `json:"propertyName,omitempty"`

	// values of the specific property
	PropertyValues []string `json:"propertyValues,omitempty"`
}

// RefinedNodeResourceStatus is the status of a RefinedNodeResource.
type RefinedNodeResourceStatus struct {
	// numeric resources
	NumericResource NumericResourceProperties `json:"numericResource,omitempty"`

	NumaStatus NumaTopologyStatus `json:"numaStatus"`
}

type NumaTopologyStatus struct {
	Sockets []SocketStatus `json:"sockets,omitempty"`
}

type SocketStatus struct {
	Numas []NumaStatus `json:"numas,omitempty"`
}

type NumaStatus struct {
	// podName
	User string `json:"user,omitempty"`
}

// NumericResourceProperties is the set of numeric properties
// these properties can be combined into DiscreteResourceProperties
// TODO: integrate these properties into DiscreteResourceProperties
type NumericResourceProperties struct {
	// Generically, there are two types of numeric resources
	// one is: resource value can be consumed, and will be decreased if it is used, like: network bandwidth
	// the other is: can not be consumed, just a hard restriction, like MBM
	// do not differentiate them here
	// numeric resource properties
	NumericProperties []NumericResourcePropertyPattern `json:"numericProperties,omitempty"`
}

type NumericResourcePropertyPattern struct {
	// property name
	PropertyName string `json:"propertyName,omitempty"`

	// values of the specific property
	PropertyCapacityValue resource.Quantity `json:"propertyCapacityValue,omitempty"`

	// values of the specific property
	PropertyAllocatableValue resource.Quantity `json:"propertyAllocatableValue,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RefinedNodeResourceList is a collection of RefinedNodeResource objects.
type RefinedNodeResourceList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of RefinedNodeResource
	Items []RefinedNodeResource `json:"items"`
}
