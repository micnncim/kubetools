package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	annotationKeyRestart = "kubectl.kubernetes.io/restartedAt"
)

var (
	kubeconfig = flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Absolute path to a kubeconfig file")
	namespace  = flag.String("n", "", "Namespace where a resource to be patched belongs to")
)

var (
	schm   = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(schm)
)

func init() {
	_ = scheme.AddToScheme(schm)
}

func main() {
	flag.Parse()

	name := flag.Arg(0)

	c, err := client()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := restart(ctx, c, name, *namespace); err != nil {
		log.Fatal(err)
	}
}

func client() (*kubernetes.Clientset, error) {
	c, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func restart(ctx context.Context, client *kubernetes.Clientset, name, namespace string) error {
	deploy, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get Deployment: %w", err)
	}

	before, err := encode(deploy)
	if err != nil {
		return fmt.Errorf("failed to encode Deployment: %w", err)
	}

	if deploy.Spec.Template.Annotations == nil {
		deploy.Spec.Template.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.Annotations[annotationKeyRestart] = time.Now().Format(time.RFC3339)

	after, err := encode(deploy)
	if err != nil {
		return fmt.Errorf("failed to encode Deployment with patch: %w", err)
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(before, after, deploy)
	if err != nil {
		return fmt.Errorf("failed to create patch: %w", err)
	}

	_, err = client.AppsV1().Deployments(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch %s/%s: %w", namespace, name, err)
	}

	return nil
}

func encode(obj runtime.Object) ([]byte, error) {
	enc := unstructured.NewJSONFallbackEncoder(codecs.LegacyCodec(schm.PrioritizedVersionsAllGroups()...))

	return runtime.Encode(enc, obj)
}
