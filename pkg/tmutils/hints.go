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
	"encoding/json"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/bitmask"
)

type ResHints struct {
	Resource string `json:"R"`
	Hints    []Hint `json:"H"`
}

type Hint struct {
	Mask      string `json:"M"`
	Preferred bool   `json:"P"`
}

func (ht Hint) ToTM() topologymanager.TopologyHint {
	tmht := topologymanager.TopologyHint{
		Preferred:        ht.Preferred,
		NUMANodeAffinity: bitmask.NewEmptyBitMask(),
	}
	num := len(ht.Mask)
	for idx := 0; idx < num; idx++ {
		if ht.Mask[idx] == '1' {
			tmht.NUMANodeAffinity.Add(num - 1 - idx)
		}
	}
	return tmht
}

func addHint(allHints map[string][]topologymanager.TopologyHint, rh ResHints) {
	hints := allHints[rh.Resource]
	for _, ht := range rh.Hints {
		hints = append(hints, ht.ToTM())
	}
	allHints[rh.Resource] = hints
}

func ParseJSONHints(rawHints []string) (map[string][]topologymanager.TopologyHint, error) {
	var err error
	allHints := make(map[string][]topologymanager.TopologyHint)
	for _, rawHint := range rawHints {
		var rh ResHints
		err = json.Unmarshal([]byte(rawHint), &rh)
		if err != nil {
			return allHints, err
		}

		addHint(allHints, rh)
	}
	return allHints, nil
}

// cpu:[{01 true} {10 true} {11 false}]
func ParseGOHints(rawHints []string) (map[string][]topologymanager.TopologyHint, error) {
	allHints := make(map[string][]topologymanager.TopologyHint)
	for _, rawHint := range rawHints {
		data := strings.SplitN(rawHint, ":", 2)
		rh := ResHints{
			Resource: strings.TrimSpace(data[0]),
			Hints:    ParseGOProviderHints(data[1]),
		}
		addHint(allHints, rh)
	}
	return allHints, nil
}

func ParseGOProviderHints(rawHints string) []Hint {
	if len(rawHints) == 0 {
		return nil
	}
	items := strings.FieldsFunc(unquoteHints(rawHints), func(r rune) bool {
		return r == '{'
	})
	var hints []Hint
	for _, item := range items {
		hintData := strings.SplitN(unquoteHintItem(strings.TrimSpace(item)), " ", 2)
		hints = append(hints, Hint{
			Mask:      hintData[0],
			Preferred: hintData[1] == "true",
		})
	}
	return hints
}

func unquoteHints(s string) string {
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	return s
}

func unquoteHintItem(s string) string {
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	return s
}
