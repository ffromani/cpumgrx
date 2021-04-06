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
 * Copyright 2020 Red Hat, Inc.
 */

package cpuset

import (
	"sort"
	"strconv"
	"strings"
)

// Parse takes a string representing a cpuset definition, and returns as sorted slice of ints
func Parse(s string) ([]int, error) {
	cpus := Empty()
	if len(s) == 0 {
		return cpus, nil
	}

	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if !strings.Contains(item, "-") {
			// single cpu: "2"
			cpuid, err := strconv.Atoi(item)
			if err != nil {
				return cpus, err
			}
			cpus = append(cpus, cpuid)
			continue
		}

		// range of cpus: "0-3"
		cpuRange := strings.SplitN(item, "-", 2)
		cpuBegin, err := strconv.Atoi(cpuRange[0])
		if err != nil {
			return cpus, err
		}
		cpuEnd, err := strconv.Atoi(cpuRange[1])
		if err != nil {
			return cpus, err
		}
		for cpuid := cpuBegin; cpuid <= cpuEnd; cpuid++ {
			cpus = append(cpus, cpuid)
		}
	}

	sort.Ints(cpus)
	return cpus, nil
}
