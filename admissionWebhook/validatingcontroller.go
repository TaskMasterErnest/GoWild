package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (whs *WebhookServer) serveValidatingRequest(req *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	// verify that we are handling only Pod resources
	if req.Resource != podResource {
		return nil, fmt.Errorf("unsupported Resource: %v", req.Resource)
	}

	// parse incoming Pod specification
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Pod: %w", err)
	}

	// perform validations
	var errors field.ErrorList

	// Example validation; Require "team" annotation
	if pod.Annotations["team"] == "" {
		errors = append(errors, field.Required(
			field.NewPath("metadata", "annotations", "team"),
			"team annotation is required",
		))
	}

	// check for validation failures
	if len(errors) > 0 {
		return &admissionv1.AdmissionResponse{
			UID:     req.UID,
			Allowed: false, // deny the request
			Result: &metav1.Status{
				Message: errors.ToAggregate().Error(),
				Code:    http.StatusForbidden,
				Reason:  metav1.StatusReasonInvalid,
			},
		}, nil
	}

	// return successful validation
	return &admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
		Result: &metav1.Status{
			Code:   http.StatusOK,
			Reason: "Validation Passed",
		},
	}, nil
}
