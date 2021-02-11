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
 * Copyright 2021 Red Hat, Inc.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"flag"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	"github.com/spf13/pflag"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"

	"github.com/fromanirh/cpumgrx/pkg/cpumgrx"
	"github.com/fromanirh/cpumgrx/pkg/tmutils"
)

func main() {
	// Add klog flags
	klog.InitFlags(flag.CommandLine)
	// Add flags registered by imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	var policyName string
	var rawHint string
	var rawReservedCPUs string
	var machineInfoPath string
	pflag.StringVarP(&rawReservedCPUs, "reserved-cpus", "R", "0", "set reserved CPUs")
	pflag.StringVarP(&rawHint, "hint", "H", "", "set topology manager hint")
	pflag.StringVarP(&machineInfoPath, "machine-info", "M", "", "machine info path")
	pflag.StringVarP(&policyName, "policy", "P", "static", "set CPU manager Policy")
	pflag.Parse()

	podSpecPaths := pflag.Args()

	if machineInfoPath == "" {
		klog.Errorf("missing machine info JSON path")
		os.Exit(1)
	}
	if len(podSpecPaths) == 0 {
		klog.Errorf("missing pod spec path(s)")
		os.Exit(1)
	}

	reservedCPUSet := parseReservedCPUsOrDie(rawReservedCPUs)
	params := cpumgrx.Params{
		PolicyName:     policyName,
		ReservedCPUSet: reservedCPUSet,
		ReservedCPUQty: resource.MustParse(fmt.Sprintf("%d", reservedCPUSet.Size())),
		MachineInfo:    readMachineInfoOrDie(machineInfoPath),
	}
	if rawHint != "" {
		params.Hint = parseHintOrDie(rawHint)
	}

	mgrx, err := cpumgrx.NewFromParams(params)
	if err != nil {
		klog.Errorf("cpumanager creation failed: %v", err)
		os.Exit(1)
	}

	for _, podSpecPath := range podSpecPaths {
		pod := readPodSpecOrDie(podSpecPath)
		cpus, err := mgrx.Run(pod)
		if err != nil {
			klog.Errorf("cpumanager allocation failed: %v", err)
			continue
		}
		fmt.Printf("%s\n", cpus.String())
	}
}

func parseHintOrDie(rawHint string) topologymanager.TopologyHint {
	hints, err := tmutils.ParseGOHints([]string{rawHint})
	if err != nil {
		klog.Errorf("error parsing hint: %v", err)
		os.Exit(1)
	}
	if len(hints) != 1 {
		klog.Errorf("wrong hints given: %#v", hints)
		os.Exit(1)
	}
	// slightly abuse because we cannot predict the key (do we really care about the key?)
	for _, hint := range hints {
		return hint[0]
	}

	// can't happen
	return topologymanager.TopologyHint{}
}

func parseReservedCPUsOrDie(rawReservedCPUs string) cpuset.CPUSet {
	reservedCPUs, err := cpuset.Parse(rawReservedCPUs)
	if err != nil {
		klog.Errorf("bad format for reserved CPU set: %v", err)
		os.Exit(1)
	}
	return reservedCPUs
}

func readMachineInfoOrDie(machineInfoPath string) *cadvisorapi.MachineInfo {
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

func readPodSpecOrDie(podSpecPath string) *v1.Pod {
	var pod v1.Pod

	src, err := os.Open(podSpecPath)
	if err != nil {
		klog.Errorf("error opening %q: %v", podSpecPath, err)
		os.Exit(1)
	}
	defer src.Close()

	dec := k8syaml.NewYAMLOrJSONDecoder(src, 1024)
	if err := dec.Decode(&pod); err != nil {
		klog.Errorf("error decoding %q: %v", podSpecPath, err)
		os.Exit(1)
	}

	return &pod
}
