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
	"flag"
	"os"

	"github.com/spf13/pflag"

	"k8s.io/klog/v2"

	"github.com/ffromani/cpumgrx/pkg/machineinformer"
)

func main() {
	handle := machineinformer.Handle{
		Out: os.Stdout,
	}

	// Add klog flags
	klog.InitFlags(flag.CommandLine)
	// Add flags registered by imported packages
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.StringVarP(&handle.RootDirectory, "root-dir", "r", "", "use <arg> as root prefix - use this if run inside a container")
	pflag.BoolVarP(&handle.RawOutput, "raw-output", "R", false, "emit full output - including machine-identifiable parts")
	pflag.Parse()

	handle.Run()
}
