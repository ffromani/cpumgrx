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

package tmutils

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/bitmask"
)

var rawJSONHints []string = []string{
	// cpu:[{01 true} {10 true} {11 false}]
	`{"R":"cpu", "H":[{"M":"01","P":true},{"M":"10","P":true},{"M":"11","P":false}]}`,
	// nvidia.com/gpu:[{01 true} {11 false}]
	`{"R":"nvidia.com/gpu", "H":[{"M":"01","P":true},{"M":"11","P":false}]}`,
	// openshift.io/intelsriov:[{10 true} {11 false}]
	`{"R":"openshift.io/intelsriov", "H":[{"M":"10","P":true},{"M":"11","P":false}]}`,
}

var rawGOHints []string = []string{
	"cpu:[{01 true} {10 true} {11 false}]",
	"nvidia.com/gpu:[{01 true} {11 false}]",
	"openshift.io/intelsriov:[{10 true} {11 false}]",
}

func TestParseHintsCompare(t *testing.T) {
	jsonHints, err := ParseJSONHints(rawJSONHints)
	if err != nil {
		t.Errorf("failed to parse hints from JSON: %v", err)
	}
	goHints, err := ParseGOHints(rawGOHints)
	if err != nil {
		t.Errorf("failed to parse hints from GO: %v", err)
	}
	if !reflect.DeepEqual(jsonHints, goHints) {
		t.Errorf("parsed hints are different: json=%#v go=%#v", jsonHints, goHints)
	}
}

func TestParseGoHints(t *testing.T) {
	rawHints := []string{
		"hugepages-2Mi:[{11 true}]",
		"memory:[{11 true}]",
		"cpu:[{01 true} {10 true} {11 false}]",
	}
	expectedHints := map[string][]topologymanager.TopologyHint{
		"hugepages-2Mi": {
			{
				NUMANodeAffinity: MustNewBitMask(t, 0, 1),
				Preferred:        true,
			},
		},
		"memory": {
			{
				NUMANodeAffinity: MustNewBitMask(t, 0, 1),
				Preferred:        true,
			},
		},
		"cpu": {
			{
				NUMANodeAffinity: MustNewBitMask(t, 0),
				Preferred:        true,
			},
			{
				NUMANodeAffinity: MustNewBitMask(t, 1),
				Preferred:        true,
			},
			{
				NUMANodeAffinity: MustNewBitMask(t, 0, 1),
				Preferred:        false,
			},
		},
	}

	goHints, err := ParseGOHints(rawHints)
	if err != nil {
		t.Errorf("failed to parse hints from GO: %v", err)
	}
	if !reflect.DeepEqual(goHints, expectedHints) {
		t.Errorf("parsed hints are different: got=%#v expected=%#v", goHints, expectedHints)
	}
}

func MustNewBitMask(t *testing.T, bits ...int) bitmask.BitMask {
	bm, err := bitmask.NewBitMask(bits...)
	if err != nil {
		t.Fatalf("unexpected error in NewBitMask: %v", err)
	}
	return bm
}
