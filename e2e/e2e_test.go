// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	v1 "github.com/fiaas/fiaas-go-client/pkg/apis/fiaas.schibsted.io/v1"
	fiaasclientset "github.com/fiaas/fiaas-go-client/pkg/client/clientset/versioned"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var applicationTests = []struct {
	expectedYamlFilePath string
	application          v1.Application
}{
	{"expected/application.yml", v1.Application{
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
			Config:      v1.Config{"version": 3},
			AdditionalAnnotations: &v1.AdditionalLabelsOrAnnotations{
				Status: map[string]string{
					"pipeline.finn.no/CallbackURL": "http://example.com/callback",
				},
			},
		},
	}},
	{"expected/minimal.yml", v1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minimal",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fiaas.schibsted.io/v1",
			Kind:       "Application",
		},
		Spec: v1.ApplicationSpec{
			Application: "minimal",
			Image:       "fiaas/fiaas-deploy-daemon:latest",
			Config:      v1.Config{"version": 3},
		},
	}},
	{"expected/fullconfig.yml", v1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "full",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fiaas.schibsted.io/v1",
			Kind:       "Application",
		},
		Spec: v1.ApplicationSpec{
			Application: "full",
			Image:       "fiaas/fiaas-deploy-daemon:latest",
			Config:      createFullConfig(),
		},
	}},
}

func TestCreate(t *testing.T) {
	clientset, err := createClient()
	if err != nil {
		t.Fatalf("failed to create client: %s", err)
	}
	for i, testcase := range applicationTests {
		t.Run(fmt.Sprintf("%d/%s", i, testcase.expectedYamlFilePath), func(t *testing.T) {
			expected, err := applicationFromYaml(testcase.expectedYamlFilePath)
			if err != nil {
				t.Fatal(err)
			}

			applicationsClient := clientset.FiaasV1().Applications(testcase.application.Namespace)
			_, err = applicationsClient.Create(&testcase.application)
			if err != nil {
				t.Fatalf("failed to create application: %s", err)
			}

			defer func() {
				err := applicationsClient.Delete(testcase.application.Name, &metav1.DeleteOptions{})
				if err != nil && !apimachineryerrors.IsNotFound(err) {
					t.Fatalf("failed to delete application: %s", err)
				}
			}()
			actual, err := applicationsClient.Get(testcase.application.Name, metav1.GetOptions{})
			if err != nil {
				t.Fatalf("failed to get application: %s", err)
			}

			assert.Equal(t, expected.ObjectMeta.Name, actual.ObjectMeta.Name)
			assert.Equal(t, expected.ObjectMeta.Labels, actual.ObjectMeta.Labels)
			assert.Equal(t, expected.ObjectMeta.Annotations, actual.ObjectMeta.Annotations)
			assert.Equal(t, expected.Spec, actual.Spec)

		})
	}
}

func createClient() (*fiaasclientset.Clientset, error) {
	kubeconfigPath, ok := os.LookupEnv("KIND_KUBECONFIG")
	if !ok {
		return nil, fmt.Errorf("$KIND_KUBECONFIG must be set for e2e test")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to build kubeconfig from kubeconfigPath %s: %s", kubeconfigPath, err)
	}
	clientset, err := fiaasclientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create new clientset: %s", err)
	}
	return clientset, nil
}

func applicationFromYaml(yamlFilePath string) (*v1.Application, error) {
	yamlBytes, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s", yamlFilePath)
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(yamlBytes, nil, &v1.Application{})
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %s", yamlFilePath, err)
	}
	application, ok := obj.(*v1.Application)
	if !ok {
		return nil, fmt.Errorf("decoded file %s is not v1.Application", yamlFilePath)
	}
	return application, nil
}

func createFullConfig() v1.Config {
	config := v1.Config{
		"admin_access": true,
		"annotations": map[string]interface{}{
			"deployment": map[string]interface{}{
				"m": "n",
				"o": "p",
			},
			"horizontal_pod_autoscaler": map[string]interface{}{
				"i": "j",
				"k": "l",
			},
			"ingress": map[string]interface{}{
				"e": "f",
				"g": "h",
			},
			"pod": map[string]interface{}{
				"x": true,
				"z": true,
			},
			"service": map[string]interface{}{
				"a": "b",
				"c": "d",
			}},
		"extensions": map[string]interface{}{
			"strongbox": map[string]interface{}{
				"aws_region": "eu-central-1",
				"groups":     []interface{}{"secretgroup1", "secretgroup2"},
				"iam_role":   "arn:aws:iam::12345678:role/the-role-name",
			},
		},
		"healthchecks": map[string]interface{}{
			"liveness": map[string]interface{}{
				"http": map[string]interface{}{
					"http_headers": map[string]interface{}{
						"X-Custom-Header": "liveness-stuff",
					},
					"path": "/health",
					"port": "a",
				},
			},
			"readiness": map[string]interface{}{
				"failure_threshold":     6,
				"initial_delay_seconds": 5,
				"period_seconds":        5,
				"success_threshold":     2,
				"tcp":                   map[string]interface{}{"port": "b"},
				"timeout_seconds":       2,
			}},
		"ingress": []interface{}{
			map[string]interface{}{
				"host": "www.example.com",
				"paths": []interface{}{
					map[string]interface{}{
						"path": "/a",
						"port": "a",
					},
				},
			},
		},
		"labels": map[string]interface{}{
			"deployment": map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			"horizontal_pod_autoscaler": map[string]interface{}{
				"e": "f",
				"g": "h",
			},
			"ingress": map[string]interface{}{
				"i": "j",
				"k": "l",
			},
			"pod": map[string]interface{}{
				"q": "r",
				"s": "u",
			},
			"service": map[string]interface{}{
				"m": "n",
				"o": "p",
			},
		},
		"metrics": map[string]interface{}{
			"datadog": map[string]interface{}{
				"enabled": true,
				"tags": map[string]interface{}{
					"tag1": "value1",
					"tag2": "value2",
				},
			},
			"prometheus": map[string]interface{}{
				"enabled": true,
				"path":    "/prometheus-metrics-here",
				"port":    "a",
			},
		},
		"ports": []interface{}{
			map[string]interface{}{
				"name":        "a",
				"port":        1337,
				"protocol":    "http",
				"target_port": 31337,
			},
			map[string]interface{}{
				"name":        "b",
				"port":        1338,
				"protocol":    "tcp",
				"target_port": 31338,
			},
		},
		"replicas": map[string]interface{}{
			"cpu_threshold_percentage": 60,
			"maximum":                  20,
			"minimum":                  10,
		},
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    2,
				"memory": "1024Mi",
			},
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "512Mi",
			},
		},
		"secrets_in_environment": true,
		"version":                3,
	}
	return config
}
