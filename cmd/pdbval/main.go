package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"sigs.k8s.io/yaml"
)

var (
	pdb    = flag.String("pdb", "", "A filepath to PodDisruptionBudget manifest")
	deploy = flag.String("deploy", "", "A filepath to Deployment manifest")
)

func main() {
	flag.Parse()

	if *pdb == "" || *deploy == "" {
		fmt.Fprintf(os.Stderr, "-pdb and -deploy must be specified\n")
		os.Exit(1)
	}

	if err := run(*pdb, *deploy); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(pdbManifest, deployManifest string) error {
	pdb, err := pdbFromManifest(pdbManifest)
	if err != nil {
		return err
	}
	deploy, err := deployFromManifest(deployManifest)
	if err != nil {
		return err
	}

	if !isTarget(pdb, deploy) {
		fmt.Fprintln(os.Stdout, "namespace or label selector not matched")
		return nil
	}

	return validate(pdb, deploy)
}

func pdbFromManifest(manifest string) (*policyv1beta1.PodDisruptionBudget, error) {
	b, err := ioutil.ReadFile(manifest)
	if err != nil {
		return nil, err
	}

	var pdb policyv1beta1.PodDisruptionBudget
	if err := yaml.Unmarshal(b, &pdb); err != nil {
		return nil, err
	}
	return &pdb, nil
}

func deployFromManifest(manifest string) (*appsv1.Deployment, error) {
	b, err := ioutil.ReadFile(manifest)
	if err != nil {
		return nil, err
	}

	var deploy appsv1.Deployment
	if err := yaml.Unmarshal(b, &deploy); err != nil {
		return nil, err
	}
	return &deploy, nil
}
