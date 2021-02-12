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

package cpumgrx

import (
	"fmt"
	"time"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/cm/containermap"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

const (
	// caveat: we will make the reconcile loop a NOP anyway, so any random time interval is fine (being irrelevant)
	reconcilePeriod = 10 * time.Minute
)

type Params struct {
	PolicyName         string
	Hint               topologymanager.TopologyHint
	MachineInfo        *cadvisorapi.MachineInfo
	ReservedCPUQty     resource.Quantity
	ReservedCPUSet     cpuset.CPUSet
	StateFileDirectory string
}

func fakeActivePods() []*v1.Pod {
	return []*v1.Pod{}
}

type fakeRuntimeService struct{}

func (rs fakeRuntimeService) UpdateContainerResources(id string, resources *runtimeapi.LinuxContainerResources) error {
	return nil
}

type fakeTMStore struct {
	Hint topologymanager.TopologyHint
}

func (tm fakeTMStore) GetAffinity(podUID string, containerName string) topologymanager.TopologyHint {
	return tm.Hint
}

type fakeSourcesReady struct{}

func (s *fakeSourcesReady) AddSource(source string) {}
func (s *fakeSourcesReady) AllReady() bool {
	// this will disable the reconcile loop
	return false
}

// fakePodStatusProvider knows how to provide status for a pod. It's intended to be used by other components
// that need to introspect status.
type fakePodStatusProvider struct{}

// GetPodStatus returns the cached status for the provided pod UID, as well as whether it
// was a cache hit.
func (psp fakePodStatusProvider) GetPodStatus(uid types.UID) (v1.PodStatus, bool) {
	// returning false makes the caller skip
	return v1.PodStatus{}, false
}

type CpuMgrx struct {
	cpuMgr     cpumanager.Manager
	fakeTm     fakeTMStore
	fakeRs     fakeRuntimeService
	policyName string

	sourcesReady      *fakeSourcesReady
	podStatusProvider fakePodStatusProvider
	initialContainers containermap.ContainerMap
}

func (cmx *CpuMgrx) GetPolicyName() string {
	return cmx.policyName
}

func (cmx *CpuMgrx) String() string {
	return "N/A"
}

func (cmx *CpuMgrx) Run(pod *v1.Pod) (cpuset.CPUSet, error) {
	cnt := &pod.Spec.Containers[0]
	err := cmx.cpuMgr.Allocate(pod, cnt)
	if err != nil {
		return cpuset.CPUSet{}, err
	}

	state := cmx.cpuMgr.State()
	cpus, ok := state.GetCPUSet(string(pod.UID), cnt.Name)
	if !ok {
		return cpuset.CPUSet{}, fmt.Errorf("GetCPUSet returned false")
	}
	return cpus, nil
}

func NewFromParams(params Params) (*CpuMgrx, error) {
	nodeAllocatableReservation := v1.ResourceList{
		v1.ResourceCPU: params.ReservedCPUQty,
	}
	fakeTm := fakeTMStore{
		Hint: params.Hint,
	}

	mgr, err := cpumanager.NewManager(params.PolicyName, reconcilePeriod, params.MachineInfo, params.ReservedCPUSet, nodeAllocatableReservation, params.StateFileDirectory, fakeTm)
	if err != nil {
		return nil, err
	}

	fakeRs := fakeRuntimeService{}
	cpuMgrx := CpuMgrx{
		cpuMgr:     mgr,
		fakeRs:     fakeRs,
		fakeTm:     fakeTm,
		policyName: params.PolicyName,

		// TODO: always empty
		// TODO: allow to load state to check more complex allocations? is the state file sufficient?
		initialContainers: containermap.ContainerMap{},
		sourcesReady:      new(fakeSourcesReady),
		podStatusProvider: fakePodStatusProvider{},
	}

	if err := cpuMgrx.cpuMgr.Start(fakeActivePods, cpuMgrx.sourcesReady, cpuMgrx.podStatusProvider, fakeRs, cpuMgrx.initialContainers); err != nil {
		return nil, err
	}
	return &cpuMgrx, nil
}
