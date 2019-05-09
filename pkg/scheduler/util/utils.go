/*
Copyright 2017 The Kubernetes Authors.

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

package util

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

// GetPodFullName returns a name that uniquely identifies a pod.
func GetPodFullName(pod *v1.Pod) string {
	// Use underscore as the delimiter because it is not allowed in pod name
	// (DNS subdomain format).
	return pod.Name + "_" + pod.Namespace
}

// GetPodStartTime returns start time of the given pod or current timestamp
// if it hasn't started yet.
func GetPodStartTime(pod *v1.Pod) *metav1.Time {
	if pod.Status.StartTime != nil {
		return pod.Status.StartTime
	}
	// Assumed pods and bound pods that haven't started don't have a StartTime yet.
	return &metav1.Time{Time: time.Now()}
}

// GetEarliestPodStartTime returns the earliest start time of all pods that
// have the highest priority among all victims.
func GetEarliestPodStartTime(victims *extenderv1.Victims) *metav1.Time {
	if len(victims.Pods) == 0 {
		// should not reach here.
		klog.Errorf("victims.Pods is empty. Should not reach here.")
		return nil
	}

	earliestPodStartTime := GetPodStartTime(victims.Pods[0])
	maxPriority := podutil.GetPodPriority(victims.Pods[0])

	for _, pod := range victims.Pods {
		if podutil.GetPodPriority(pod) == maxPriority {
			if GetPodStartTime(pod).Before(earliestPodStartTime) {
				earliestPodStartTime = GetPodStartTime(pod)
			}
		} else if podutil.GetPodPriority(pod) > maxPriority {
			maxPriority = podutil.GetPodPriority(pod)
			earliestPodStartTime = GetPodStartTime(pod)
		}
	}

	return earliestPodStartTime
}

// MoreImportantPod return true when priority of the first pod is higher than
// the second one. If two pods' priorities are equal, compare their StartTime.
// It takes arguments of the type "interface{}" to be used with SortableList,
// but expects those arguments to be *v1.Pod.
func MoreImportantPod(pod1, pod2 *v1.Pod) bool {
	p1 := podutil.GetPodPriority(pod1)
	p2 := podutil.GetPodPriority(pod2)
	if p1 != p2 {
		return p1 > p2
	}
	return GetPodStartTime(pod1).Before(GetPodStartTime(pod2))
}

// GetPodAffinityTerms gets pod affinity terms by a pod affinity object.
func GetPodAffinityTerms(podAffinity *v1.PodAffinity) (terms []v1.PodAffinityTerm) {
	if podAffinity != nil {
		if len(podAffinity.RequiredDuringSchedulingIgnoredDuringExecution) != 0 {
			terms = podAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		}
		// TODO: Uncomment this block when implement RequiredDuringSchedulingRequiredDuringExecution.
		//if len(podAffinity.RequiredDuringSchedulingRequiredDuringExecution) != 0 {
		//	terms = append(terms, podAffinity.RequiredDuringSchedulingRequiredDuringExecution...)
		//}
	}
	return terms
}

// GetPodAntiAffinityTerms gets pod affinity terms by a pod anti-affinity.
func GetPodAntiAffinityTerms(podAntiAffinity *v1.PodAntiAffinity) (terms []v1.PodAffinityTerm) {
	if podAntiAffinity != nil {
		if len(podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution) != 0 {
			terms = podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		}
		// TODO: Uncomment this block when implement RequiredDuringSchedulingRequiredDuringExecution.
		//if len(podAntiAffinity.RequiredDuringSchedulingRequiredDuringExecution) != 0 {
		//	terms = append(terms, podAntiAffinity.RequiredDuringSchedulingRequiredDuringExecution...)
		//}
	}
	return terms
}

// CanPodBePreempted indicates whether the pod can be preempted
func CanPodBePreempted(pod *v1.Pod) bool {
	//if pod.Spec.CanBePreempted == nil {
	//	return false
	//}
	//
	//return *pod.Spec.CanBePreempted
	return false
}

// PreemptionScopeEqual compares the preemption scope label
// return true if they are equal
func PreemptionScopeEqual(pod1, pod2 *v1.Pod) bool {
	if pod1.Labels == nil || pod2.Labels == nil {
		return false
	}

	if len(pod1.Labels[PreemptionScopeKey]) == 0 || len(pod2.Labels[PreemptionScopeKey]) == 0 {
		return false
	}

	return pod1.Labels[PreemptionScopeKey] == pod2.Labels[PreemptionScopeKey]
}

// HigherPriorityPod return true when priority of the first pod is higher than
// the second one. It takes arguments of the type "interface{}" to be used with
// SortableList, but expects those arguments to be *v1.Pod.
func HigherPriorityPod(pod1, pod2 *v1.Pod) bool {
	p1 := podutil.GetPodPriority(pod1)
	p2 := podutil.GetPodPriority(pod2)
	if p1 != p2 {
		return p1 > p2
	}

	pod1HasGPU := HasResource(pod1, ResourceGPU)
	pod2HasGPU := HasResource(pod2, ResourceGPU)
	if pod1HasGPU != pod2HasGPU {
		if pod1HasGPU {
			return true
		} else {
			return false
		}
	}

	return true
}

func LessImportantPod(pod1, pod2 interface{}) bool {
	// compare priority first
	p1 := podutil.GetPodPriority(pod1.(*v1.Pod))
	p2 := podutil.GetPodPriority(pod2.(*v1.Pod))
	if p1 != p2 {
		return p1 < p2
	}
	// if the priorities are equal, compare resource types.
	// order is: GPU, Memory, CPU
	// Since memory and cpu is not that special, put them behind request size comparisons.
	pod1HasGPU := HasResource(pod1.(*v1.Pod), ResourceGPU)
	pod2HasGPU := HasResource(pod2.(*v1.Pod), ResourceGPU)
	if pod1HasGPU != pod2HasGPU {
		if pod1HasGPU {
			return false
		} else {
			return true
		}
	}

	// compare request resource
	// and the resources order is: GPU, Memory, CPU
	// the smaller, the quicker to reprieve
	if pod1HasGPU {
		pod1GPURequest := getPodRequest(pod1.(*v1.Pod), ResourceGPU, resource.DecimalSI)
		pod2GPURequest := getPodRequest(pod2.(*v1.Pod), ResourceGPU, resource.DecimalSI)
		result := pod1GPURequest.Cmp(*pod2GPURequest)
		if result < 0 {
			// pod2 request is greater than pod1,
			// since GPU resource is precious, pod1 is less important
			return true
		} else if result > 0 {
			return false
		}
	}

	pod1MemoryRequest := getPodRequest(pod1.(*v1.Pod), v1.ResourceMemory, resource.BinarySI)
	pod2MemoryRequest := getPodRequest(pod2.(*v1.Pod), v1.ResourceMemory, resource.BinarySI)
	result := pod1MemoryRequest.Cmp(*pod2MemoryRequest)
	if result < 0 {
		// pod2 request is greater than pod1
		return false
	} else if result > 0 {
		return true
	}

	pod1CPURequest := getPodRequest(pod1.(*v1.Pod), v1.ResourceCPU, resource.DecimalSI)
	pod2CPURequest := getPodRequest(pod2.(*v1.Pod), v1.ResourceCPU, resource.DecimalSI)
	result = pod1CPURequest.Cmp(*pod2CPURequest)
	if result < 0 {
		// pod2 request is greater than pod1
		return false
	} else if result > 0 {
		return true
	}

	return true
}

func HasResource(pod *v1.Pod, resourceType v1.ResourceName) bool {
	zeroQuantity := resource.MustParse("0")
	for _, container := range pod.Spec.Containers {
		for key, quantity := range container.Resources.Requests {
			if key == resourceType && quantity.Cmp(zeroQuantity) == 1 {
				return true
			}
		}
		for key, quantity := range container.Resources.Limits {
			if key == resourceType && quantity.Cmp(zeroQuantity) == 1 {
				return true
			}
		}
	}

	for _, container := range pod.Spec.InitContainers {
		for key, quantity := range container.Resources.Requests {
			if key == resourceType && quantity.Cmp(zeroQuantity) == 1 {
				return true
			}
		}
		for key, quantity := range container.Resources.Limits {
			if key == resourceType && quantity.Cmp(zeroQuantity) == 1 {
				return true
			}
		}
	}

	return false
}

const (
	// hardcode GPU name here
	// TODO: support more GPU names
	ResourceGPU v1.ResourceName = "nvidia.com/gpu"
)

func getPodRequest(pod *v1.Pod, resourceType v1.ResourceName, format resource.Format) *resource.Quantity {
	result := resource.NewQuantity(0, format)
	for _, container := range pod.Spec.Containers {
		for key, value := range container.Resources.Requests {
			if key == resourceType {
				result.Add(value)
			}
		}
	}

	for _, container := range pod.Spec.InitContainers {
		for key, value := range container.Resources.Requests {
			if key == resourceType {
				if result.Cmp(value) < 0 {
					result.SetMilli(value.MilliValue())
				}
			}
		}
	}

	return result
}

const PreemptionScopeKey = "PreemptionScopeKey"
