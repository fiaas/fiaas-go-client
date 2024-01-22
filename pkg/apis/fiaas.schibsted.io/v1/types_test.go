package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


var ExampleConfig = Config{
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


func TestConfigDeepCopy(t *testing.T) {
	cp := &ExampleConfig

	var out Config
	cp.DeepCopyInto(&out)

	// verify data was actually copied
	assert.Equal(t, *cp, out)

	// modify some of the data in various levels and types in out
	out["admin_access"] = false
	if replicas, ok := out["replicas"].(map[string]interface{}); ok {
		replicas["minimum"] = 1
		replicas["maximum"] = 5
	} else {
		t.Fatalf("replicas was not map[string]interface{}: %v", replicas)
	}

	if metrics, ok := out["metrics"].(map[string]interface{}); ok {
		if prometheus, ok := metrics["prometheus"].(map[string]interface{}); ok {
			prometheus["port"] = "other-port"
		} else {
			t.Fatalf("metrics.prometheus was not map[string]interface{}: %v", prometheus)
		}
	} else {
		t.Fatalf("metrics was not map[string]interface{}: %v", metrics)
	}

	if ingress, ok := out["ingress"].([]interface{}); ok {
		if item := ingress[0].(map[string]interface{}); ok {
			item["host"] = "other.example.com"
			if pathmappings, ok := item["paths"].([]interface{}); ok {
				if pathmapping, ok := pathmappings[0].(map[string]interface{}); ok {
					pathmapping["path"] = "/other"
					pathmapping["port"] = "other-port"
				} else {
					t.Fatalf("ingress.paths[0] was not map[string]interface{}: %v", pathmapping)
				}

			} else {
				t.Fatalf("ingress.paths was not []interface{}: %v", pathmappings)
			}
		} else {
			t.Fatalf("ingress[0] was not map[string]interface{}: %v", item)
		}
	} else {
		t.Fatalf("ingress was not []interface{}")
	}

	// should no longer be equal since out is a copy
	assert.NotEqual(t, *cp, out)

	expected := &Config{
		"admin_access": false,
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
				"host": "other.example.com",
				"paths": []interface{}{
					map[string]interface{}{
						"path": "/other",
						"port": "other-port",
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
				"port":    "other-port",
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
			"maximum":                  5,
			"minimum":                  1,
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

	assert.Equal(t, *expected, out)
}

func TestConfigDeepCopyMultiple(t *testing.T) {
	// verify calling gob.Register via DeepCopy/DeepCopyInto multiple times is OK
	out1 := ExampleConfig.DeepCopy()
	out2 := ExampleConfig.DeepCopy()

	var outInto1 Config
	ExampleConfig.DeepCopyInto(&outInto1)
	var outInto2 Config
	ExampleConfig.DeepCopyInto(&outInto2)

	assert.Equal(t, ExampleConfig, out1)
	assert.Equal(t, ExampleConfig, out2)
	assert.Equal(t, ExampleConfig, outInto1)
	assert.Equal(t, ExampleConfig, outInto2)
}
