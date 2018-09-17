package main

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"shadowproject/common/errors"
)

type Status struct {
	Critical         bool  `json:"overall"`            // Something is very wrong and this node should be replaced
	Overloaded       bool  `json:"overloaded"`         // True if it's overloaded
	ContainerBackend bool  `json:"container_backend"`  // Status of the container backend
	CPUUtilization   uint8 `json:"cpu_utilization"`    // Last five minutes CPU utilization (based on load average) (0-100 %)
	Memory           uint8 `json:"memory_utilization"` // Memory utilization (0-100 %)
}

func numberOfCPUs() uint {
	processors, err := cpu.Counts(true)
	if err != nil {
		panic(errors.ShadowError{
			Origin:         err,
			VisibleMessage: "health check error",
		})
	}

	return uint(processors)
}

// Returns current load 5
func load5() float64 {
	load, err := load.Avg()
	if err != nil {
		panic(errors.ShadowError{
			Origin:         err,
			VisibleMessage: "health check error",
		})
	}

	return load.Load5
}

// Return utilization of memory
func memory() float64 {
	memory, err := mem.VirtualMemory()
	if err != nil {
		panic(errors.ShadowError{
			Origin:         err,
			VisibleMessage: "health check error",
		})
	}

	return float64(memory.Available) / float64(memory.Total)
}

// Returns structure with server status
func HealthCheck() *Status {
	var status Status

	status.ContainerBackend = dockerDriver.Status()
	status.CPUUtilization = uint8(load5() / float64(numberOfCPUs()) * 100)
	status.Memory = uint8(memory() * 100)

	status.Critical = !status.ContainerBackend
	status.Overloaded = status.CPUUtilization > 100 || status.Memory > 80

	return &status
}
