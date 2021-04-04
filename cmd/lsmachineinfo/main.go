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
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"k8s.io/klog/v2"

	"github.com/google/cadvisor/fs"
	infov1 "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/machine"
	"github.com/google/cadvisor/utils/sysfs"
)

func main() {
	var rawOutput bool

	// Add klog flags
	klog.InitFlags(flag.CommandLine)
	// Add flags registered by imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.BoolVarP(&rawOutput, "raw-output", "R", false, "emit full output - including machine-identifiable parts")
	pflag.Parse()

	fsInfo := newFakeFsInfo()
	sysFs := sysfs.NewRealSysFs()
	inHostNamespace := true

	info, err := machine.Info(sysFs, fsInfo, inHostNamespace)
	if err != nil {
		klog.Fatalf("Cannot get machine info: %v")
	}

	outInfo := info
	if !rawOutput {
		outInfo = cleanInfo(info)
	}

	json.NewEncoder(os.Stdout).Encode(outInfo)
}

func cleanInfo(in *infov1.MachineInfo) *infov1.MachineInfo {
	out := in.Clone()
	out.MachineID = ""
	out.SystemUUID = ""
	out.BootID = ""
	for i := 0; i < len(out.NetworkDevices); i++ {
		out.NetworkDevices[i].MacAddress = ""
	}
	out.CloudProvider = infov1.UnknownProvider
	out.InstanceType = infov1.UnknownInstance
	out.InstanceID = infov1.UnNamedInstance
	return out
}

type fakeFsInfo struct {
	notImplemented error
}

func newFakeFsInfo() fs.FsInfo {
	return fakeFsInfo{
		notImplemented: fmt.Errorf("not implemented"),
	}
}

func (ffi fakeFsInfo) GetGlobalFsInfo() ([]fs.Fs, error) {
	return nil, ffi.notImplemented
}

func (ffi fakeFsInfo) GetFsInfoForPath(mountSet map[string]struct{}) ([]fs.Fs, error) {
	return nil, ffi.notImplemented
}

func (ffi fakeFsInfo) GetDirUsage(dir string) (fs.UsageInfo, error) {
	return fs.UsageInfo{}, ffi.notImplemented
}

func (ffi fakeFsInfo) GetDeviceInfoByFsUUID(uuid string) (*fs.DeviceInfo, error) {
	return nil, ffi.notImplemented
}

func (ffi fakeFsInfo) GetDirFsDevice(dir string) (*fs.DeviceInfo, error) {
	return nil, ffi.notImplemented
}

func (ffi fakeFsInfo) GetDeviceForLabel(label string) (string, error) {
	return "", ffi.notImplemented
}

func (ffi fakeFsInfo) GetLabelsForDevice(device string) ([]string, error) {
	return nil, ffi.notImplemented
}

func (ffi fakeFsInfo) GetMountpointForDevice(device string) (string, error) {
	return "", ffi.notImplemented
}
