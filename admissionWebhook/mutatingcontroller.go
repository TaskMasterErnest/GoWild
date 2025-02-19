package main

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

func (whs *WebhookServer) serveMutatingRequest(req *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	// verify if the resource is a Pod
	if req.Resource != podResource {
		return nil, fmt.Errorf("unsupported Resource: %s", req.Resource)
	}

	// parse Pod with error handling
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Pod: %w", err)
	}

	// prepare patches
	var patches []patchOperation

	// add annotation patches
	if pod.Annotations == nil {
		patches = append(patches, patchOperation{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: make(map[string]string),
		})
	}
	// adding my custom annotation
	patches = append(patches, patchOperation{
		Op:    "add",
		Path:  "/metadata/annotations/PodModified",
		Value: "true",
	})

	// add pod resource limits and requst patches
	for i := range pod.Spec.Containers {
		resourcePatch := patchOperation{
			Op:   "replace", // this is for replacing
			Path: fmt.Sprintf("/spec/containers/%d/resources", i),
			Value: map[string]interface{}{
				"limits": map[string]string{
					"cpu":    "200m",
					"memory": "256Mi",
				},
				"requests": map[string]string{
					"cpu":    "100m",
					"memory": "128Mi",
				},
			},
		}
		// if the resources already exist, modify them to have these limits and requests
		// if the do not, add this resource block
		if pod.Spec.Containers[i].Resources.Limits == nil && pod.Spec.Containers[i].Resources.Requests == nil {
			resourcePatch.Op = "add"
		}

		patches = append(patches, resourcePatch)
	}

	// marshal JSON
	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patch: %w", err)
	}

	// return admissionResponse
	return &admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}, nil
}
