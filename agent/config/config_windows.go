// +build windows
// Copyright 2014-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/aws/amazon-ecs-agent/agent/dockerclient"
	"github.com/aws/amazon-ecs-agent/agent/utils"
)

const (
	// AgentCredentialsAddress is used to serve the credentials for tasks.
	AgentCredentialsAddress = "127.0.0.1"
	// defaultAuditLogFile specifies the default audit log filename
	defaultCredentialsAuditLogFile = `log\audit.log`
	// When using IAM roles for tasks on Windows, the credential proxy consumes port 80
	httpPort = 80
	// Remote Desktop / Terminal Services
	rdpPort = 3389
	// RPC client
	rpcPort = 135
	// Server Message Block (SMB) over TCP
	smbPort = 445
	// HTTP port for Windows Remote Management (WinRM) listener
	winRMPortHTTP = 5985
	// HTTPS port for Windows Remote Management (WinRM) listener
	winRMPortHTTPS = 5986
	// DNS client
	dnsPort = 53
	// NetBIOS over TCP/IP
	netBIOSPort = 139
	// defaultContainerStartTimeout specifies the value for container start timeout duration
	defaultContainerStartTimeout = 8 * time.Minute
	// minimumContainerStartTimeout specifies the minimum value for starting a container
	minimumContainerStartTimeout = 2 * time.Minute
	// default image pull inactivity time is extra time needed on container extraction
	defaultImagePullInactivityTimeout = 3 * time.Minute
)

// DefaultConfig returns the default configuration for Windows
func DefaultConfig() Config {
	programData := utils.DefaultIfBlank(os.Getenv("ProgramData"), `C:\ProgramData`)
	ecsRoot := filepath.Join(programData, "Amazon", "ECS")
	dataDir := filepath.Join(ecsRoot, "data")
	platformVariables := PlatformVariables{
		CPUUnbounded:    false,
		MemoryUnbounded: false,
	}
	return Config{
		DockerEndpoint: "npipe:////./pipe/docker_engine",
		ReservedPorts: []uint16{
			DockerReservedPort,
			DockerReservedSSLPort,
			AgentIntrospectionPort,
			AgentCredentialsPort,
			rdpPort,
			rpcPort,
			smbPort,
			winRMPortHTTP,
			winRMPortHTTPS,
			dnsPort,
			netBIOSPort,
		},
		ReservedPortsUDP: []uint16{},
		DataDir:          dataDir,
		// DataDirOnHost is identical to DataDir for Windows because we do not
		// run as a container
		DataDirOnHost:                       dataDir,
		ReservedMemory:                      0,
		AvailableLoggingDrivers:             []dockerclient.LoggingDriver{dockerclient.JSONFileDriver, dockerclient.NoneDriver, dockerclient.AWSLogsDriver},
		TaskCleanupWaitDuration:             DefaultTaskCleanupWaitDuration,
		DockerStopTimeout:                   defaultDockerStopTimeout,
		ContainerStartTimeout:               defaultContainerStartTimeout,
		ImagePullInactivityTimeout:          defaultImagePullInactivityTimeout,
		CredentialsAuditLogFile:             filepath.Join(ecsRoot, defaultCredentialsAuditLogFile),
		CredentialsAuditLogDisabled:         false,
		ImageCleanupDisabled:                false,
		MinimumImageDeletionAge:             DefaultImageDeletionAge,
		ImageCleanupInterval:                DefaultImageCleanupTimeInterval,
		NumImagesToDeletePerCycle:           DefaultNumImagesToDeletePerCycle,
		NumNonECSContainersToDeletePerCycle: DefaultNumNonECSContainersToDeletePerCycle,
		ContainerMetadataEnabled:            false,
		TaskCPUMemLimit:                     ExplicitlyDisabled,
		PlatformVariables:                   platformVariables,
		TaskMetadataSteadyStateRate:         DefaultTaskMetadataSteadyStateRate,
		TaskMetadataBurstRate:               DefaultTaskMetadataBurstRate,
		SharedVolumeMatchFullConfig:         false, //only requiring shared volumes to match on name, which is default docker behavior
		PollMetrics:                         false,
		PollingMetricsWaitDuration:          DefaultPollingMetricsWaitDuration,
	}
}

func (cfg *Config) platformOverrides() {
	// Enabling task IAM roles for Windows requires the credential proxy to run on port 80,
	// so we reserve this port by default when that happens.
	if cfg.TaskIAMRoleEnabled {
		if cfg.ReservedPorts == nil {
			cfg.ReservedPorts = []uint16{}
		}
		cfg.ReservedPorts = append(cfg.ReservedPorts, httpPort)
	}

	// ensure TaskResourceLimit is disabled
	cfg.TaskCPUMemLimit = ExplicitlyDisabled

	cpuUnbounded := utils.ParseBool(os.Getenv("ECS_ENABLE_CPU_UNBOUNDED_WINDOWS_WORKAROUND"), false)
	memoryUnbounded := utils.ParseBool(os.Getenv("ECS_ENABLE_MEMORY_UNBOUNDED_WINDOWS_WORKAROUND"), false)

	platformVariables := PlatformVariables{
		CPUUnbounded:    cpuUnbounded,
		MemoryUnbounded: memoryUnbounded,
	}
	cfg.PlatformVariables = platformVariables
}

// platformString returns platform-specific config data that can be serialized
// to string for debugging
func (cfg *Config) platformString() string {
	return ""
}
