package nodeinfo

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"

	"k8s.io/kubernetes/pkg/scheduler/util"
)

type NodeRefinedResourceInfo struct {
	nodeName        string
	cpuProperty     CPUProperty
	gpuProperty     GPUProperty
	diskProperty    DiskProperty
	memoryProperty  MemoryProperty
	networkProperty NetworkProperty
	otherProperty   OtherProperty

	// numeric resource status
	// key is numeric resource name, value is the node capacity of that specific resource
	// the unit of value is ignored here. we need to keep the unit same between u8s and k8s
	// the units of request and capacity must be equal too.
	// we need to cache CPU, Memory and Numa Capacity too, in order for numa use cases
	// TODO: support more units in k8s Quantity, for example: bps, Bps...
	numericResourceProperty map[string]int64
	// allocatableResources v1.ResourceList
	// requestedResources   v1.ResourceList
	numaTopologyStatus NumaTopologyStatus
}

type NumaTopologyStatus struct {
	freeNumasInSockets map[int]sets.Int
	numasOfPods        map[string]sets.Int
	socketsOfNumas     map[int]int
}

const (
	cpuPrefix     = "cpu"
	gpuPrefix     = "gpu"
	diskPrefix    = "disk"
	memoryPrefix  = "memory"
	networkPrefix = "network"
	otherPrefix   = "other"
)

func NewNodeRefinedResourceInfo(nodeName string) *NodeRefinedResourceInfo {
	return &NodeRefinedResourceInfo{
		nodeName:        nodeName,
		cpuProperty:     make(map[string]sets.String),
		gpuProperty:     make(map[string]sets.String),
		diskProperty:    make(map[string]sets.String),
		memoryProperty:  make(map[string]sets.String),
		networkProperty: make(map[string]sets.String),
		otherProperty:   make(map[string]sets.String),

		numericResourceProperty: make(map[string]int64),
		numaTopologyStatus: NumaTopologyStatus{
			freeNumasInSockets: make(map[int]sets.Int),
			numasOfPods:        make(map[string]sets.Int),
			socketsOfNumas:     make(map[int]int),
		},
	}
}

func (nrri *NodeRefinedResourceInfo) AddDiscreteResourceProperty(prefix, propertyKey, propertyValue string) {
	switch prefix {
	case cpuPrefix:
		if nrri.cpuProperty[propertyKey] == nil {
			nrri.cpuProperty[propertyKey] = sets.NewString()
		}
		nrri.cpuProperty[propertyKey].Insert(propertyValue)
	case gpuPrefix:
		if nrri.gpuProperty[propertyKey] == nil {
			nrri.gpuProperty[propertyKey] = sets.NewString()
		}
		nrri.gpuProperty[propertyKey].Insert(propertyValue)
	case diskPrefix:
		if nrri.diskProperty[propertyKey] == nil {
			nrri.diskProperty[propertyKey] = sets.NewString()
		}
		nrri.diskProperty[propertyKey].Insert(propertyValue)
	case memoryPrefix:
		if nrri.memoryProperty[propertyKey] == nil {
			nrri.memoryProperty[propertyKey] = sets.NewString()
		}
		nrri.memoryProperty[propertyKey].Insert(propertyValue)
	case networkPrefix:
		if nrri.networkProperty[propertyKey] == nil {
			nrri.networkProperty[propertyKey] = sets.NewString()
		}
		nrri.networkProperty[propertyKey].Insert(propertyValue)
	case otherPrefix:
		if nrri.otherProperty[propertyKey] == nil {
			nrri.otherProperty[propertyKey] = sets.NewString()
		}
		nrri.otherProperty[propertyKey].Insert(propertyValue)
	default:
		fmt.Println("error prefix, should not occur")
	}
}

func (nrri *NodeRefinedResourceInfo) AddNumericResourceProperties(propertyKey string, propertyValue resource.Quantity) {
	/*if propertyKey == "cpu" {
		nrri.cpuCapacity = propertyValue
	} else if propertyKey == "memory" {
		nrri.memoryCapacity = propertyValue
	} else {
		nrri.numericResourceProperty[propertyKey] = propertyValue.Value()
	}*/
	if propertyKey == util.CPURefinedResourceKey {
		nrri.numericResourceProperty[propertyKey] = propertyValue.MilliValue()
	} else {
		nrri.numericResourceProperty[propertyKey] = propertyValue.Value()
	}
}

func (nrri *NodeRefinedResourceInfo) AddNumaTopologyStatus(socketID, numaID int, podName string) {
	if nrri.numaTopologyStatus.freeNumasInSockets[socketID] == nil {
		nrri.numaTopologyStatus.freeNumasInSockets[socketID] = sets.NewInt()
	}
	nrri.numaTopologyStatus.freeNumasInSockets[socketID].Insert(numaID)
	nrri.numaTopologyStatus.socketsOfNumas[numaID] = socketID
	if podName != "" {
		if nrri.numaTopologyStatus.numasOfPods[podName] == nil {
			nrri.numaTopologyStatus.numasOfPods[podName] = sets.NewInt()
		}
		nrri.numaTopologyStatus.numasOfPods[podName].Insert(numaID)
		nrri.numaTopologyStatus.freeNumasInSockets[socketID].Delete(numaID)
	}
}

func (nrri *NodeRefinedResourceInfo) RemovePod(pod *v1.Pod) error {
	if nrri == nil {
		return nil
	}
	numaTopology := nrri.GetNumaTopologyStatus()
	numaTopologyCopy := numaTopology.Clone()
	numas := numaTopologyCopy.GetNumasOfPods()[pod.GetName()]
	socketsOfNumas := numaTopologyCopy.GetSocketsOfNumas()
	for _, numa := range numas.List() {
		socket := socketsOfNumas[numa]
		numaTopologyCopy.freeNumasInSockets[socket].Insert(numa)
	}
	nrri.numaTopologyStatus = numaTopologyCopy
	return nil
}

func (nrri *NodeRefinedResourceInfo) AddPod(pod *v1.Pod) {
	if nrri == nil {
		return
	}
	numaTopology := nrri.GetNumaTopologyStatus()
	numaTopologyCopy := numaTopology.Clone()
	numas := numaTopologyCopy.GetNumasOfPods()[pod.GetName()]
	socketsOfNumas := numaTopologyCopy.GetSocketsOfNumas()
	for _, numa := range numas.List() {
		socket := socketsOfNumas[numa]
		numaTopologyCopy.freeNumasInSockets[socket].Delete(numa)
	}
	nrri.numaTopologyStatus = numaTopologyCopy
}

type CPUProperty map[string]sets.String

type GPUProperty map[string]sets.String

type DiskProperty map[string]sets.String

type MemoryProperty map[string]sets.String

type NetworkProperty map[string]sets.String

type OtherProperty map[string]sets.String

func (nrri *NodeRefinedResourceInfo) GetNodeName() string {
	return nrri.nodeName
}

func (nrri *NodeRefinedResourceInfo) GetCPUProperties() map[string]sets.String {
	return nrri.cpuProperty
}

func (nrri *NodeRefinedResourceInfo) GetGPUProperties() map[string]sets.String {
	return nrri.gpuProperty
}

func (nrri *NodeRefinedResourceInfo) GetDiskProperties() map[string]sets.String {
	return nrri.diskProperty
}

func (nrri *NodeRefinedResourceInfo) GetMemoryProperties() map[string]sets.String {
	return nrri.memoryProperty
}

func (nrri *NodeRefinedResourceInfo) GetNetworkProperties() map[string]sets.String {
	return nrri.networkProperty
}

func (nrri *NodeRefinedResourceInfo) GetOtherProperties() map[string]sets.String {
	return nrri.otherProperty
}

func (nrri *NodeRefinedResourceInfo) GetNumericResourceProperties() map[string]int64 {
	return nrri.numericResourceProperty
}

func (nrri *NodeRefinedResourceInfo) Clone() *NodeRefinedResourceInfo {
	if nrri == nil {
		return nil
	}
	refinedNodeInfo := &NodeRefinedResourceInfo{
		nodeName:        nrri.nodeName,
		cpuProperty:     make(map[string]sets.String),
		gpuProperty:     make(map[string]sets.String),
		diskProperty:    make(map[string]sets.String),
		memoryProperty:  make(map[string]sets.String),
		networkProperty: make(map[string]sets.String),
		otherProperty:   make(map[string]sets.String),

		numericResourceProperty: make(map[string]int64),
	}

	for k, v := range nrri.GetCPUProperties() {
		if refinedNodeInfo.cpuProperty[k] == nil {
			refinedNodeInfo.cpuProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.cpuProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetGPUProperties() {
		if refinedNodeInfo.gpuProperty[k] == nil {
			refinedNodeInfo.gpuProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.gpuProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetDiskProperties() {
		if refinedNodeInfo.diskProperty[k] == nil {
			refinedNodeInfo.diskProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.diskProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetMemoryProperties() {
		if refinedNodeInfo.memoryProperty[k] == nil {
			refinedNodeInfo.memoryProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.memoryProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetNetworkProperties() {
		if refinedNodeInfo.networkProperty[k] == nil {
			refinedNodeInfo.networkProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.networkProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetOtherProperties() {
		if refinedNodeInfo.otherProperty[k] == nil {
			refinedNodeInfo.otherProperty[k] = sets.NewString(v.List()...)
		} else {
			refinedNodeInfo.otherProperty[k].Insert(v.List()...)
		}
	}

	for k, v := range nrri.GetNumericResourceProperties() {
		refinedNodeInfo.numericResourceProperty[k] = v
	}

	refinedNodeInfo.numaTopologyStatus = nrri.GetNumaTopologyStatus().Clone()

	return refinedNodeInfo
}

func (nrri *NodeRefinedResourceInfo) GetNumaTopologyStatus() NumaTopologyStatus {
	return nrri.numaTopologyStatus
}

func (nts NumaTopologyStatus) GetFreeNumasInSockets() map[int]sets.Int {
	return nts.freeNumasInSockets
}

func (nts NumaTopologyStatus) GetNumasOfPods() map[string]sets.Int {
	return nts.numasOfPods
}

func (nts NumaTopologyStatus) GetSocketsOfNumas() map[int]int {
	return nts.socketsOfNumas
}

func (nts NumaTopologyStatus) GetSocketNum() int {
	return len(nts.freeNumasInSockets)
}

func (nts NumaTopologyStatus) GetNumaNum() int {
	return len(nts.socketsOfNumas)
}

func (nts NumaTopologyStatus) Clone() NumaTopologyStatus {
	clone := NumaTopologyStatus{
		freeNumasInSockets: make(map[int]sets.Int),
		numasOfPods:        make(map[string]sets.Int),
		socketsOfNumas:     make(map[int]int),
	}
	for k, v := range nts.GetFreeNumasInSockets() {
		clone.freeNumasInSockets[k] = sets.NewInt(v.List()...)
	}
	for k, v := range nts.GetNumasOfPods() {
		clone.numasOfPods[k] = sets.NewInt(v.List()...)
	}
	for k, v := range nts.GetSocketsOfNumas() {
		clone.socketsOfNumas[k] = v
	}
	return clone
}
