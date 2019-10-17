// +build e2e

package e2e

import (
	"io/ioutil"
	"os"
	"testing"

	v1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	fiaasclientset "github.com/fiaas/fiaas-go-client/pkg/client/clientset/versioned"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreate(t *testing.T) {
	kubeconfigPath, ok := os.LookupEnv("KIND_KUBECONFIG")
	if !ok {
		t.Fatalf("$KIND_KUBECONFIG must be set for e2e test")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		t.Fatalf("unable to build kubeconfig from kubeconfigPath %s: %s", kubeconfigPath, err)
	}
	clientset, err := fiaasclientset.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("unable to create new clientset: %s", err)
	}

	expectedConfig := make(v1.Config)
	expectedConfig["version"] = float64(3)

	application := v1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testapplication",
			Namespace: "default",
			Labels: map[string]string{
				"app": "testapplication",
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fiaas.schibsted.io/v1",
			Kind:       "Application",
		},
		Spec: v1.ApplicationSpec{
			Application: "testapplication",
			Image:       "fiaas/fiaas-deploy-daemon:latest",
			Config:      expectedConfig,
			AdditionalAnnotations: &v1.AdditionalLabelsOrAnnotations{
				Status: map[string]string{
					"pipeline.finn.no/CallbackURL": "http://example.com/callback",
				},
			},
		},
	}

	expectedFile := "expected/application.yaml"
	yamlBytes, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read file %s", expectedFile)
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(yamlBytes, nil, &v1.Application{})
	if err != nil {
		t.Fatalf("failed to decode %s: %s", expectedFile, err)
	}
	expected := obj.(*v1.Application)

	applicationsClient := clientset.FiaasV1().Applications(application.Namespace)
	_, err = applicationsClient.Create(&application)
	if err != nil {
		t.Fatalf("failed to create application: %s", err)
	}

	actual, err := applicationsClient.Get(application.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get applicaiton: %s", err)
	}

	assert.Equal(t, expected.ObjectMeta.Name, actual.ObjectMeta.Name)
	assert.Equal(t, expected.ObjectMeta.Labels, actual.ObjectMeta.Labels)
	assert.Equal(t, expected.ObjectMeta.Annotations, actual.ObjectMeta.Annotations)
	assert.Equal(t, expected.Spec, actual.Spec)
}
