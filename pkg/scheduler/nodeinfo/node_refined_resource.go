package nodeinfo

import (
	"fmt"

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

	return refinedNodeInfo
}
