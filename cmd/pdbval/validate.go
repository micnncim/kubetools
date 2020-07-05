package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// isTarget determine the given Deployment is a target for the given PodDisruptionBudget
// using PodDisruptionBudget's label selector.
// Currently it support only matchLabels, not matchExpressions.
func isTarget(pdb *policyv1beta1.PodDisruptionBudget, deploy *appsv1.Deployment) bool {
	for k, v := range pdb.Spec.Selector.MatchLabels {
		value, ok := deploy.Spec.Template.ObjectMeta.Labels[k]
		if !ok || value != v {
			return false
		}
	}
	return true
}

func validate(pdb *policyv1beta1.PodDisruptionBudget, deploy *appsv1.Deployment) error {
	replicas := int(*deploy.Spec.Replicas)
	if replicas == 0 {
		replicas = 1 // default value is 1.
	}

	if pdb.Spec.MaxUnavailable != nil {
		maxUnavailable, err := intstr.GetValueFromIntOrPercent(pdb.Spec.MaxUnavailable, replicas, true)
		if err != nil {
			return err
		}
		healthyReplicas := replicas - maxUnavailable
		disruptionAllowedReplicas := maxUnavailable

		switch {
		case healthyReplicas < 1:
			return fmt.Errorf(
				"PodDisruptionBudget(%s): maxUnavailable(%s) is greater than or equal to Deployment(%s) replicas(%d)",
				pdb.Name, pdb.Spec.MaxUnavailable.String(), deploy.Name, replicas,
			)

		case disruptionAllowedReplicas < 1:
			return fmt.Errorf(
				"PodDisruptionBudget(%s): maxUnavailable(%s) is less than 1",
				pdb.Name, pdb.Spec.MaxUnavailable.String(),
			)
		}
	}

	if pdb.Spec.MinAvailable != nil {
		minAvailable, err := intstr.GetValueFromIntOrPercent(pdb.Spec.MinAvailable, replicas, true)
		if err != nil {
			return err
		}
		healthyReplicas := minAvailable
		disruptionAllowedReplicas := replicas - minAvailable

		switch {
		case healthyReplicas < 1:
			return fmt.Errorf(
				"PodDisruptionBudget(%s): minAvailable(%s) is less than 1",
				pdb.Name, pdb.Spec.MinAvailable.String(),
			)

		case disruptionAllowedReplicas < 1:
			return fmt.Errorf(
				"PodDisruptionBudget(%s): minAvailable(%s) is greater than or equal to Deployment(%s) replicas(%d)",
				pdb.Name, pdb.Spec.MinAvailable.String(), deploy.Name, replicas,
			)
		}
	}

	return nil
}
