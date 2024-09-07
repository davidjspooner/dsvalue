package path

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/davidjspooner/dsvalue/pkg/reflected"
	"github.com/davidjspooner/dsvalue/pkg/value"
	"gopkg.in/yaml.v3"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		input          string
		expectedError  error
		expectedString string
	}{
		{
			input:          ".foo.bar",
			expectedString: ".foo.bar",
			expectedError:  nil,
		},
		{
			input:          ".baz[0].qux",
			expectedString: ".baz[0].qux",
			expectedError:  nil,
		},
		{
			input:          ".foo[-1].bar[:].baz",
			expectedString: ".foo[-1].bar[:].baz",
			expectedError:  nil,
		},
		{
			input:          ".foo[-2:-1].bar[:].baz",
			expectedString: ".foo[-2:-1].bar[:].baz",
			expectedError:  nil,
		},
		{
			input:          "[].field[*]",
			expectedString: "[:].field[:]",
			expectedError:  nil,
		},
		{
			input:         "0",
			expectedError: &ErrInvalidPath{Path: "0", Inner: errors.New("expected '. or [', but got '0'")},
		},
		{
			input:         ".foo[0",
			expectedError: &ErrInvalidPath{Path: ".foo[0", Inner: errors.New("expected ': or ]', but got <EOF>")},
		},
		{
			input:         ".foo[0].",
			expectedError: &ErrInvalidPath{Path: ".foo[0].", Inner: errors.New("expected 'identifier', but got <EOF>")},
		},
		{
			input:          ".",
			expectedString: ".",
		},
		{
			input:         ".foo bar",
			expectedError: &ErrInvalidPath{Path: ".foo bar", Inner: errors.New("expected '. or [', but got ' '")},
		},
	}

	for _, test := range tests {
		path, err := CompilePath(test.input)
		if test.expectedError == nil && err != nil {
			t.Errorf("Parsing path: %q, expected no error, but got %q", test.input, err)
		} else if test.expectedError != nil && err == nil {
			t.Errorf("Parsing path: %q, expected %q, but got no error", test.input, test.expectedError)
		} else if test.expectedError != nil && err != nil && test.expectedError.Error() != err.Error() {
			t.Errorf("Parsing path: %q, expected %q, but got %q", test.input, test.expectedError, err)
		} else if err == nil {
			output := path.String()
			if output != test.expectedString {
				t.Errorf("Parsing path: %q, expected %q, but got %q", test.input, test.expectedString, path)
			}
		}
	}
}

var sampleYaml = `
apiVersion: v1
kind: Service
metadata:
  annotations:
    meta.helm.sh/release-name: traefik
    meta.helm.sh/release-namespace: ingress
    metallb.universe.tf/ip-allocated-from-pool: default-pool
  creationTimestamp: "2024-08-11T03:50:04Z"
  labels:
    app.kubernetes.io/instance: traefik-ingress
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: traefik
    helm.sh/chart: traefik-30.0.2
  name: traefik
  namespace: ingress
  resourceVersion: "4083"
  uid: d54d8410-018a-42fb-9dd1-b4b4a3df3f11
spec:
  allocateLoadBalancerNodePorts: true
  clusterIP: 10.152.183.108
  clusterIPs:
  - 10.152.183.108
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: web
    nodePort: 30494
    port: 80
    protocol: TCP
    targetPort: web
  - name: websecure
    nodePort: 30657
    port: 443
    protocol: TCP
    targetPort: websecure
  selector:
    app.kubernetes.io/instance: traefik-ingress
    app.kubernetes.io/name: traefik
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer:
    ingress:
    - ip: 192.168.201.128
      ipMode: VIP
`

func TestPathEvaluateFor(t *testing.T) {

	var obj any
	d := yaml.NewDecoder(strings.NewReader(sampleYaml))
	err := d.Decode(&obj)
	if err != nil {
		t.Errorf("Error decoding yaml: %v", err)
		return
	}

	path, err := CompilePath(".status.loadBalancer.ingress[:].ip")
	if err != nil {
		t.Errorf("Error parsing path: %v", err)
		return
	}

	rv := reflect.ValueOf(obj)
	object, err := reflected.NewReflectedObject(rv, value.UnknownSource)
	if err != nil {
		t.Errorf("Error creating reflected object: %v", err)
		return
	}
	result, err := path.EvaluateFor(object)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	s := result.WithoutSource()
	t.Logf("Value: %s", s)
	_ = result
}
