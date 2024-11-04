/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2022 Red Hat, Inc.
 */

package main

import (
	"encoding/json"
	"os"

	"flag"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	"github.com/sanity-io/litter"
	"github.com/spf13/pflag"

	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/topology"
)

func main() {
	// Add klog flags
	klog.InitFlags(flag.CommandLine)
	// Add flags registered by imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	var machineInfoPath string
	pflag.StringVarP(&machineInfoPath, "machine-info", "M", "", "machine info path")
	pflag.Parse()

	if machineInfoPath == "" {
		klog.Errorf("missing machine info JSON path")
		os.Exit(1)
	}

	machineInfo := mustReadMachineInfo(machineInfoPath)

	topo, err := topology.Discover(machineInfo)
	if err != nil {
		klog.Errorf("topology discovery failed: %v", err)
		os.Exit(2)
	}

	litter.Dump(topo)
}

func mustReadMachineInfo(machineInfoPath string) *cadvisorapi.MachineInfo {
	var machineInfo cadvisorapi.MachineInfo
	src, err := os.Open(machineInfoPath)
	if err != nil {
		klog.Errorf("error opening %q: %v", machineInfoPath, err)
		os.Exit(1)
	}
	defer src.Close()

	dec := json.NewDecoder(src)
	if err := dec.Decode(&machineInfo); err != nil {
		klog.Errorf("error decoding %q: %v", machineInfoPath, err)
		os.Exit(1)
	}

	return &machineInfo
}
