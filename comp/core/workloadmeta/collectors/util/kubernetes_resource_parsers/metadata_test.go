// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver && test

package kubernetesresourceparsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/DataDog/datadog-agent/comp/core/workloadmeta/collectors/util"
	workloadmeta "github.com/DataDog/datadog-agent/comp/core/workloadmeta/def"
)

func TestParse_ParsePartialObjectMetadata(t *testing.T) {

	testcases := []struct {
		name                           string
		gvr                            schema.GroupVersionResource
		partialObjectMetadata          *metav1.PartialObjectMetadata
		expected                       *workloadmeta.KubernetesMetadata
		annotationsFilter              []string
		requireErrorInitailisingParser bool
	}{
		{
			name: "deployments [namespace scoped]",
			gvr: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
			partialObjectMetadata: &metav1.PartialObjectMetadata{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-app",
					Namespace:   "default",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
			},
			expected: &workloadmeta.KubernetesMetadata{
				EntityID: workloadmeta.EntityID{
					Kind: workloadmeta.KindKubernetesMetadata,
					ID:   string(util.GenerateKubeMetadataEntityID("apps", "deployments", "default", "test-app")),
				},
				EntityMeta: workloadmeta.EntityMeta{
					Name:        "test-app",
					Namespace:   "default",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
				GVR: &schema.GroupVersionResource{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				},
			},
			requireErrorInitailisingParser: false,
		},
		{
			name: "namespaces [cluster scoped]",
			gvr: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "namespaces",
			},
			partialObjectMetadata: &metav1.PartialObjectMetadata{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
			},
			expected: &workloadmeta.KubernetesMetadata{
				EntityID: workloadmeta.EntityID{
					Kind: workloadmeta.KindKubernetesMetadata,
					ID:   string(util.GenerateKubeMetadataEntityID("", "namespaces", "", "test-namespace")),
				},
				EntityMeta: workloadmeta.EntityMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
				GVR: &schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "namespaces",
				},
			},
			requireErrorInitailisingParser: false,
		},
		{
			name: "namespaces [cluster scoped], with well-formatted annotation filters",
			gvr: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "namespaces",
			},
			partialObjectMetadata: &metav1.PartialObjectMetadata{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "weird-annotation": "weird-value"},
				},
			},
			expected: &workloadmeta.KubernetesMetadata{
				EntityID: workloadmeta.EntityID{
					Kind: workloadmeta.KindKubernetesMetadata,
					ID:   string(util.GenerateKubeMetadataEntityID("", "namespaces", "", "test-namespace")),
				},
				EntityMeta: workloadmeta.EntityMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
				GVR: &schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "namespaces",
				},
			},
			annotationsFilter: []string{"weird-annotation"},

			requireErrorInitailisingParser: false,
		},
		{
			name: "namespaces [cluster scoped], with badly-formatted annotation filters",
			gvr: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "namespaces",
			},
			partialObjectMetadata: &metav1.PartialObjectMetadata{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "weird-annotation": "weird-value"},
				},
			},
			expected: &workloadmeta.KubernetesMetadata{
				EntityID: workloadmeta.EntityID{
					Kind: workloadmeta.KindKubernetesMetadata,
					ID:   string(util.GenerateKubeMetadataEntityID("", "namespaces", "", "test-namespace")),
				},
				EntityMeta: workloadmeta.EntityMeta{
					Name:        "test-namespace",
					Namespace:   "",
					Labels:      map[string]string{"l1": "v1", "l2": "v2", "l3": "v3"},
					Annotations: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
				},
				GVR: &schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "namespaces",
				},
			},
			annotationsFilter:              []string{"/foo1)("},
			requireErrorInitailisingParser: true,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(tt *testing.T) {
			parser, err := NewMetadataParser(test.gvr, test.annotationsFilter)

			if test.requireErrorInitailisingParser {
				require.Errorf(tt, err, "should have failed to create parser")
			} else {
				require.NoErrorf(tt, err, "should have not failed to create parser")
				entity := parser.Parse(test.partialObjectMetadata)
				storedMetadata, ok := entity.(*workloadmeta.KubernetesMetadata)
				require.True(t, ok)
				assert.Equal(t, test.expected, storedMetadata)
			}
		})
	}
}
