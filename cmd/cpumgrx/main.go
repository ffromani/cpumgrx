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
	"path/filepath"
	"regexp"

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
	var podTemplateMode bool
	var keepState bool
	var stateFileDirectory string
	pflag.StringVarP(&rawReservedCPUs, "reserved-cpus", "R", "0", "set reserved CPUs")
	pflag.StringVarP(&rawHint, "hint", "H", "", "set topology manager hint")
	pflag.StringVarP(&machineInfoPath, "machine-info", "M", "", "machine info path")
	pflag.BoolVarP(&podTemplateMode, "pod-template-mode", "T", false, "pod template mode")
	pflag.StringVarP(&policyName, "policy", "P", "static", "set CPU manager Policy")
	pflag.BoolVarP(&keepState, "keep-state", "k", false, "keep the cpu_manager_state file")
	pflag.StringVarP(&stateFileDirectory, "state-dir", "s", ".", "directory to store the cpu_manager_state_file")
	pflag.Parse()

	args := pflag.Args()

	if machineInfoPath == "" {
		klog.Errorf("missing machine info JSON path")
		os.Exit(1)
	}
	if len(args) == 0 {
		klog.Errorf("missing args")
		os.Exit(1)
	}

	reservedCPUSet := mustParseReservedCPUs(rawReservedCPUs)
	params := cpumgrx.Params{
		PolicyName:         policyName,
		StateFileDirectory: stateFileDirectory,
		ReservedCPUSet:     reservedCPUSet,
		ReservedCPUQty:     resource.MustParse(fmt.Sprintf("%d", reservedCPUSet.Size())),
		MachineInfo:        mustReadMachineInfo(machineInfoPath),
	}
	if rawHint != "" {
		params.Hint = mustParseHint(rawHint)
	}

	var pods []*v1.Pod

	if podTemplateMode {
		cpuReqs := parseCpuReqs(args)
		for _, cpuReq := range cpuReqs {
			pods = append(pods, makePod(cpuReq))
		}
	} else {
		podSpecPaths := args
		for _, podSpecPath := range podSpecPaths {
			pods = append(pods, mustReadPodSpec(podSpecPath))
		}
	}

	mgrx, err := cpumgrx.NewFromParams(params)
	if err != nil {
		klog.Errorf("cpumanager creation failed: %v", err)
		os.Exit(1)
	}

	defer func() {
		if keepState {
			return
		}
		fullPath := filepath.Join(stateFileDirectory, "cpu_manager_state")
		klog.V(3).Infof("removing cpu_manager_state file on %q", fullPath)
		err := os.Remove(fullPath)
		if err != nil {
			klog.Warning("error removing %q: %v", fullPath, err)
		}
	}()

	for _, pod := range pods {
		if blob, err := json.Marshal(pod); err == nil {
			klog.V(5).Infof("handling pod: %s", string(blob))
		}

		cpus, err := mgrx.Run(pod)
		if err != nil {
			klog.Errorf("cpumanager allocation failed: %v", err)
			continue
		}
		fmt.Printf("%s\n", cpus.String())
	}
}

type cpuReqSpec struct {
	Name     string
	Limits   resource.Quantity
	Requests resource.Quantity
}

// name=request/limit
func parseCpuReqs(args []string) []cpuReqSpec {
	var reqsRE = regexp.MustCompile(`^(\S*)=(\S*)/(\S*)$`)
	var cpuReqs []cpuReqSpec
	for _, arg := range args {
		items := reqsRE.FindAllStringSubmatch(arg, -1)
		// items[0] is the full match
		if items == nil || len(items[0]) != 4 {
			klog.Warningf("cannot parse cpu req spec %q - skipped", arg)
			continue
		}
		cpuReqs = append(cpuReqs, cpuReqSpec{
			Name:     items[0][1],
			Limits:   resource.MustParse(items[0][2]),
			Requests: resource.MustParse(items[0][3]),
		})
	}
	return cpuReqs
}

func makePod(cpuReq cpuReqSpec) *v1.Pod {
	pod := v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Resources: v1.ResourceRequirements{
						Limits:   make(v1.ResourceList),
						Requests: make(v1.ResourceList),
					},
				},
			},
		},
	}
	// yep, that's the lazy way
	pod.Name = fmt.Sprintf("%s-pod", cpuReq.Name)
	pod.Spec.Containers[0].Name = fmt.Sprintf("%s-cnt", cpuReq.Name)
	pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU] = cpuReq.Requests
	pod.Spec.Containers[0].Resources.Limits[v1.ResourceCPU] = cpuReq.Limits
	pod.Spec.Containers[0].Resources.Requests[v1.ResourceMemory] = resource.MustParse("1Gi")
	pod.Spec.Containers[0].Resources.Limits[v1.ResourceMemory] = resource.MustParse("1Gi")
	return &pod
}

func mustParseHint(rawHint string) topologymanager.TopologyHint {
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

func mustParseReservedCPUs(rawReservedCPUs string) cpuset.CPUSet {
	reservedCPUs, err := cpuset.Parse(rawReservedCPUs)
	if err != nil {
		klog.Errorf("bad format for reserved CPU set: %v", err)
		os.Exit(1)
	}
	return reservedCPUs
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

func mustReadPodSpec(podSpecPath string) *v1.Pod {
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
