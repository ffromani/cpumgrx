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
	"os"

	"github.com/spf13/pflag"

	"k8s.io/klog/v2"

	infov1 "github.com/google/cadvisor/info/v1"

	"github.com/fromanirh/cpumgrx/pkg/machineinformer"
)

func main() {
	var rootDir string
	var rawOutput bool

	// Add klog flags
	klog.InitFlags(flag.CommandLine)
	// Add flags registered by imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.StringVarP(&rootDir, "root-dir", "r", "", "use <arg> as root prefix - use this if run inside a container")
	pflag.BoolVarP(&rawOutput, "raw-output", "R", false, "emit full output - including machine-identifiable parts")
	pflag.Parse()

	var err error
	var info *infov1.MachineInfo
	if rawOutput {
		info, err = machineinformer.GetRaw(rootDir)
	} else {
		info, err = machineinformer.Get(rootDir)
	}

	if err != nil {
		klog.Fatalf("Cannot get machine info: %v")
	}

	json.NewEncoder(os.Stdout).Encode(info)
}
