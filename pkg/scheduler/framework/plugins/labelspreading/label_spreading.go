/*
Copyright 2020 The Kubernetes Authors.

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

package labelspreading

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/defaultpodtopologyspread"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
	utilnode "k8s.io/kubernetes/pkg/util/node"
)

// CalculateSpreadPriority spreads pods across hosts and zones, considering pods belonging to the same service or replication controller.
// When a pod is scheduled, it looks for services, RCs or RSs that match the pod, then finds existing pods that match those selectors.
// It favors nodes that have fewer existing matching pods.
// i.e. it pushes the scheduler towards a node where there's the smallest number of
// pods which match the same service, RC or RS selectors as the pod being scheduled.
// Where zone information is included on the nodes, it favors nodes in zones with fewer existing matching pods.
type LabelSpread struct {
	handle framework.FrameworkHandle
}

var _ = framework.ScorePlugin(&LabelSpread{})

// LabelSpread is the name of the plugin used in the plugin registry and configurations.
const Name = "LabelSpread"

// Name returns name of the plugin. It is used in logs, etc.
func (s *LabelSpread) Name() string {
	return Name
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &LabelSpread{handle: h}, nil
}

// Score invoked at the Score extension point.
func (s *LabelSpread) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil || nodeInfo.Node() == nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", nodeName, err, nodeInfo.Node() == nil))
	}

	return s.score(pod, nodeInfo)
}

// ScoreExtensions of the Score plugin.
func (s *LabelSpread) ScoreExtensions() framework.ScoreExtensions {
	return s
}

// NormalizeScore calculates the source of each node
// based on the number of existing matching pods on the node
// where zone information is included on the nodes, it favors nodes
// in zones with fewer existing matching pods.
func (s *LabelSpread) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	countsByZone := make(map[string]int64, 10)
	maxCountByNodeName := int64(0)
	maxCountByZone := int64(0)
	for i := range scores {
		if scores[i].Score > maxCountByNodeName {
			maxCountByNodeName = scores[i].Score
		}
		nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(scores[i].Name)
		if err != nil || nodeInfo.Node() == nil {
			return framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", scores[i].Name, err, nodeInfo.Node() == nil))
		}
		zoneID := utilnode.GetZoneKey(nodeInfo.Node())
		if zoneID == "" {
			continue
		}
		countsByZone[zoneID] += scores[i].Score

	}

	for zoneID := range countsByZone {
		if countsByZone[zoneID] > maxCountByZone {
			maxCountByZone = countsByZone[zoneID]
		}
	}

	// Aggregate by-zone information
	// Compute the maximum number of pods hosted in any zone
	haveZones := len(countsByZone) != 0

	maxCountByNodeNameFloat64 := float64(maxCountByNodeName)
	maxCountByZoneFloat64 := float64(maxCountByZone)
	MaxPriorityFloat64 := float64(framework.MaxNodeScore)

	//score int - scale of 0-maxPriority
	// 0 being the lowest priority and maxPriority being the highest
	for i := range scores {
		// initializing to the default/max node score of maxPriority
		fScore := MaxPriorityFloat64
		if maxCountByNodeName > 0 {
			fScore = MaxPriorityFloat64 * (float64(maxCountByNodeName-scores[i].Score) / maxCountByNodeNameFloat64)
		}

		// If there is zone information present, incorporate it
		if haveZones {
			nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(scores[i].Name)
			if err != nil || nodeInfo.Node() == nil {
				return framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", scores[i].Name, err, nodeInfo.Node() == nil))
			}
			zoneID := utilnode.GetZoneKey(nodeInfo.Node())
			if zoneID != "" {
				zoneScore := MaxPriorityFloat64
				if maxCountByZone > 0 {
					zoneScore = MaxPriorityFloat64 * (float64(maxCountByZone-countsByZone[zoneID]) / maxCountByZoneFloat64)

				}
				fScore = (fScore * (1.0 - defaultpodtopologyspread.ZoneWeighting)) + (defaultpodtopologyspread.ZoneWeighting * zoneScore)
			}
		}

		scores[i].Score = int64(fScore)
		if klog.V(10) {
			// We explicitly don't do klog.V(10).Infof() to avoid computing all the parameters if this is
			// not logged. There is visible performance gain from it.
			klog.V(10).Infof(
				"%v -> %v: LabelSpreadPriority, Score: (%d)", pod.Name, scores[i].Name, scores[i].Score,
			)
		}
	}
	return nil
}

func (s *LabelSpread) getSelectors(pod *v1.Pod) []labels.Selector {
	return []labels.Selector{labels.Set(pod.Labels).AsSelector()}
}

// CalculateSpreadPriority spreads pods across hosts and zones, considering pods belonging to the same service or replication controller.
// When a pod is scheduled, it looks for services, RCs or RSs that match the pod, then finds existing pods that match those selectors.
// It favors nodes that have fewer existing matching pods.
// i.e. it pushes the scheduler towards a node where there's the smallest number of
// pods which match the same service, RC or RS selectors as the pod being scheduled.
// Where zone information is included on the nodes, it favors nodes in zones with fewer existing matching pods.
func (s *LabelSpread) score(pod *v1.Pod, nodeInfo *schedulernodeinfo.NodeInfo) (int64, *framework.Status) {
	node := nodeInfo.Node()
	if node == nil {
		return 0, framework.NewStatus(framework.Error, "node not found")
	}

	selectors := s.getSelectors(pod)
	if len(selectors) == 0 {
		return 0, nil
	}
	count := int64(0)
	for _, nodePod := range nodeInfo.Pods() {
		if pod.Namespace != nodePod.Namespace {
			continue
		}
		// When we are replacing a failed pod, we often see the previous
		// deleted version while scheduling the replacement.
		// Ignore the previous deleted version for spreading purposes
		// (it can still be considered for resource restrictions etc.)
		if nodePod.DeletionTimestamp != nil {
			klog.V(4).Infof("skipping pending-deleted pod: %s/%s", nodePod.Namespace, nodePod.Name)
			continue
		}
		matches := false
		for _, selector := range selectors {
			if selector.Matches(labels.Set(nodePod.ObjectMeta.Labels)) {
				matches = true
				break
			}
		}
		if matches {
			count++
		}
	}
	return count, nil
}
