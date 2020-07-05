package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_validate(t *testing.T) {
	fakeDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deploy",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
		},
	}

	type args struct {
		pdb *policyv1beta1.PodDisruptionBudget
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "PodDisruptionBudget is valid",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MaxUnavailable: intOrStringPtr(intstr.FromInt(1)),
					},
				},
			},
			wantErr: "",
		},
		{
			name: "PodDisruptionBudget has too many maxUnavailable integer",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MaxUnavailable: intOrStringPtr(intstr.FromInt(5)),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): maxUnavailable(5) is greater than or equal to Deployment(deploy) replicas(3)",
		},
		{
			name: "PodDisruptionBudget has too few maxUnavailable integer",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MaxUnavailable: intOrStringPtr(intstr.FromInt(0)),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): maxUnavailable(0) is less than 1",
		},
		{
			name: "PodDisruptionBudget has too much maxUnavailable percentage",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MaxUnavailable: intOrStringPtr(intstr.FromString("100%")),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): maxUnavailable(100%) is greater than or equal to Deployment(deploy) replicas(3)",
		},
		{
			name: "PodDisruptionBudget has too few minAvailable integer",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MinAvailable: intOrStringPtr(intstr.FromInt(0)),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): minAvailable(0) is less than 1",
		},
		{
			name: "PodDisruptionBudget has too many minAvailable integer",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MinAvailable: intOrStringPtr(intstr.FromInt(5)),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): minAvailable(5) is greater than or equal to Deployment(deploy) replicas(3)",
		},
		{
			name: "PodDisruptionBudget has too much minAvailable percentage",
			args: args{
				pdb: &policyv1beta1.PodDisruptionBudget{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pdb",
					},
					Spec: policyv1beta1.PodDisruptionBudgetSpec{
						MinAvailable: intOrStringPtr(intstr.FromString("90%")),
					},
				},
			},
			wantErr: "PodDisruptionBudget(pdb): minAvailable(90%) is greater than or equal to Deployment(deploy) replicas(3)",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate(tt.args.pdb, fakeDeploy)
			if err != nil && err.Error() != tt.wantErr {
				if diff := cmp.Diff(tt.wantErr, err.Error()); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
			}
		})
	}
}

func int32Ptr(v int32) *int32 {
	return &v
}

func intOrStringPtr(v intstr.IntOrString) *intstr.IntOrString {
	return &v
}
