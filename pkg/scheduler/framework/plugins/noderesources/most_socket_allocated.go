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

package noderesources

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// MostSocketAllocated is a score plugin that favors nodes with high Socket allocation based on requested resources.
type MostSocketAllocated struct {
	handle framework.FrameworkHandle
}

var _ = framework.ScorePlugin(&MostSocketAllocated{})

// MostSocketAllocatedName is the name of the plugin used in the plugin registry and configurations.
const MostSocketAllocatedName = "NodeResourcesMostSocketAllocated"

// Name returns name of the plugin. It is used in logs, etc.
func (ma *MostSocketAllocated) Name() string {
	return MostSocketAllocatedName
}

// Score invoked at the Score extension point.
func (ma *MostSocketAllocated) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := ma.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil || nodeInfo.Node() == nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", nodeName, err, nodeInfo.Node() == nil))
	}

	// MostSocketRequestedPriority is a priority function that favors nodes with most requested resources.
	// It calculates the percentage of bytedan.com/socket requested by pods scheduled on the node, and prioritizes
	// based on the maximum of the average of the fraction of requested to capacity.
	// Details: socket(10 * sum(requested) / capacity)
	nonZeroRequest := &schedulernodeinfo.Resource{
		ScalarResources: make(map[v1.ResourceName]int64),
	}
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		socket, ok := container.Resources.Limits[v1.ResourceBytedanceSocket]
		if ok {
			nonZeroRequest.ScalarResources[v1.ResourceBytedanceSocket] += socket.Value()
		}
	}
	return calculateSocketUsedPriority(pod, nonZeroRequest, nodeInfo)
}

// ScoreExtensions of the Score plugin.
func (ma *MostSocketAllocated) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// NewMostSocketAllocated initializes a new plugin and returns it.
func NewMostSocketAllocated(_ *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &MostSocketAllocated{
		handle: h,
	}, nil
}

// The used capacity is calculated on a scale of 0-10
// 0 being the lowest priority and 10 being the highest.
// The more resources are used the higher the score is. This function
// is almost a reversed version of least_requested_priority.calculatUnusedScore
// (10 - calculateUnusedScore). The main difference is in rounding. It was added to
// keep the final formula clean and not to modify the widely used (by users
// in their default scheduling policies) calculateSocketUsedScore.
func calculateSocketUsedScore(requested int64, capacity int64, node string) int64 {
	if capacity == 0 {
		return 0
	}
	if requested > capacity {
		klog.V(10).Infof("Combined requested resources %d from existing pods exceeds capacity %d on node %s",
			requested, capacity, node)
		return 0
	}
	return (requested * 10) / capacity
}

func getSocketFromResource(resource schedulernodeinfo.Resource) int64 {
	if resource.ScalarResources == nil {
		return 0
	}
	return resource.ScalarResources[v1.ResourceBytedanceSocket]
}

// Calculate the resource used on a node.  'node' has information about the resources on the node.
// 'pods' is a list of pods currently scheduled on the node.
func calculateSocketUsedPriority(pod *v1.Pod, podRequests *schedulernodeinfo.Resource, nodeInfo *schedulernodeinfo.NodeInfo) (int64, *framework.Status) {
	node := nodeInfo.Node()
	if node == nil {
		return 0, framework.NewStatus(framework.Error, "node not found")
	}

	allocatableResources := nodeInfo.AllocatableResource()
	totalResources := getSocketFromResource(*podRequests)
	totalResources += getSocketFromResource(nodeInfo.RequestedResource())
	socketScore := calculateSocketUsedScore(totalResources, getSocketFromResource(allocatableResources), node.Name)
	if klog.V(10) {
		// We explicitly don't do klog.V(10).Infof() to avoid computing all the parameters if this is
		// not logged. There is visible performance gain from it.
		klog.V(10).Infof(
			"%v -> %v: Most Socket Requested Priority, capacity %d bytedance.com/socket, total request %d bytedance.com/socket, score %d Socket",
			pod.Name, node.Name,
			getSocketFromResource(allocatableResources),
			totalResources,
			socketScore,
		)
	}

	return socketScore, nil
}
