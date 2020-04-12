package raw

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	info "github.com/google/cadvisor/info/v1"
)

const (
	memInfoPath     = "/proc/meminfo"
	vmStatPath      = "/proc/vmstat"
	maxUseBytesPath = "/sys/fs/cgroup/memory/memory.max_usage_in_bytes"
	failCntPath     = "/sys/fs/cgroup/memory/memory.failcnt"
)

func getParamKeyValue(line string, sep string) (key, val string, err error) {
	items := strings.Split(line, sep)
	switch len(items) {
	case 2:
		return strings.TrimSpace(items[0]), strings.TrimSpace(items[1]), nil
	default:
		return "", "", fmt.Errorf("invalid format")
	}
}

func parseValue(valueStr string) (value uint64, err error) {
	if strings.HasSuffix(valueStr, "kB") {
		valueStr = valueStr[:len(valueStr)-3]
		value, err = strconv.ParseUint(valueStr, 10, 64)
		value = value << 10
	} else {
		value, err = strconv.ParseUint(valueStr, 10, 64)
	}
	return
}

// memory.usage_in_bytes ~= free.used + free.(buff/cache) - (buff);  Used memory (calculated as total - free - buffers - cache)
// so the last is  memory.usage_in_bytes  = total - free - buffers
// memory.max_usage_in_bytes  read from cgroup
// memory.cache = meminfo Cached
// memory.Swap =  meminfo SwapTotal
// memory.MappedFile = meminfo Mapped
// memory.Workingset = usage_in_bytes -  inactive(file)
// memory.failcnt read from cgroup
func getMemoryStats() (info.MemoryStats, error) {

	result := info.MemoryStats{
		Usage:            0,
		MaxUsage:         0,
		Cache:            0,
		RSS:              0,
		Swap:             0,
		MappedFile:       0,
		WorkingSet:       0,
		Failcnt:          0,
		ContainerData:    info.MemoryStatsMemoryData{},
		HierarchicalData: info.MemoryStatsMemoryData{},
	}

	memInfo, err := loadKeyValueFromFile(memInfoPath, ":")
	if err != nil {
		return result, err
	}

	vmStat, err := loadKeyValueFromFile(vmStatPath, " ")
	if err != nil {
		return result, err
	}
	result.Usage = memInfo["MemTotal"] - memInfo["MemFree"] - memInfo["Buffers"]
	result.Cache = memInfo["Cached"]
	result.RSS = memInfo["Active"]
	result.MappedFile = memInfo["Mapped"]
	result.Swap = memInfo["SwapTotal"]
	result.WorkingSet = result.Usage - memInfo["Inactive(file)"]

	result.ContainerData.Pgfault = vmStat["pgfault"]
	result.ContainerData.Pgmajfault = vmStat["pgmajfault"]
	maxUsage, err := loadValueFromFile(maxUseBytesPath)
	if err != nil {
		return result, err
	}
	result.MaxUsage = maxUsage
	if result.MaxUsage < result.Usage {
		result.MaxUsage = result.Usage
	}

	failCnt, err := loadValueFromFile(failCntPath)
	if err != nil {
		return result, err
	}
	result.Failcnt = failCnt
	return result, nil
}

func loadValueFromFile(filePath string) (uint64, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(bytes.TrimSpace(content)), 10, 64)
}

func loadKeyValueFromFile(filePath, sep string) (map[string]uint64, error) {
	statsFile, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func() {
		_ = statsFile.Close()
	}()
	memInfo := map[string]uint64{}
	sc := bufio.NewScanner(statsFile)
	for sc.Scan() {
		t, v, err := getParamKeyValue(sc.Text(), sep)
		if err != nil {
			return nil, fmt.Errorf("failed to parse. (%q) - %v", sc.Text(), err)
		}
		val, err := parseValue(v)

		if err != nil {
			return nil, fmt.Errorf("failed to parse. (%q) - %v", sc.Text(), err)
		}
		memInfo[t] = val
	}
	return memInfo, nil
}
