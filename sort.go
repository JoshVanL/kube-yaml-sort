package main

import (
	"bytes"
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	yamlsep   = []byte("---\n")
	yamlsepnl = []byte("\n---\n")
)

type object struct {
	i   int
	obj *unstructured.Unstructured
}

func SortYAMLObjects(yamlBytes []byte) ([]byte, error) {
	// Split on '---' as per the yaml spec
	split := bytes.Split(yamlBytes, yamlsep)

	var objs []object
	for i, s := range split {
		json, err := yaml.ToJSON(s)
		if err != nil {
			return nil, err
		}

		// If json returns null then we can ignore as this is not a valid yaml object
		if bytes.Equal(json, []byte("null")) {
			continue
		}

		runObj, _, err := unstructured.UnstructuredJSONScheme.Decode(json, nil, nil)
		if err != nil {
			return nil, err
		}

		obj, ok := runObj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("failed to convert runtime object to unstructured: %+v", runObj)
		}

		objs = append(objs, object{
			i:   i,
			obj: obj,
		})
	}

	sort.SliceStable(objs, func(i, j int) bool {
		if objs[i].obj.GetAPIVersion() != objs[j].obj.GetAPIVersion() {
			return objs[i].obj.GetAPIVersion() < objs[j].obj.GetAPIVersion()
		}

		if objs[i].obj.GetKind() != objs[j].obj.GetKind() {
			return objs[i].obj.GetKind() < objs[j].obj.GetKind()
		}

		if objs[i].obj.GetNamespace() != objs[j].obj.GetNamespace() {
			return objs[i].obj.GetNamespace() < objs[j].obj.GetNamespace()
		}

		return objs[i].obj.GetName() < objs[j].obj.GetName()
	})

	var sorted [][]byte
	for _, obj := range objs {
		sorted = append(sorted, bytes.TrimSpace(split[obj.i]))
	}

	return append(bytes.Join(sorted, yamlsepnl), '\n'), nil
}
