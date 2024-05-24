package e2e

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/support/kind"
)

func TestE2E(t *testing.T) {
	// 横着した
	kindClusterName := envconf.RandomName("my-cluster", 16)
	p := kind.NewProvider().WithName(kindClusterName)
	_, err := p.Create(context.Background())
	if err != nil {
		t.Fatalf("failed to create cluster: %v", err)
	}
	t.Cleanup(func() {
		err := p.Destroy(context.Background())
		if err != nil {
			t.Fatalf("failed to destroy cluster: %v", err)
		}
	})
	c, err := klient.NewWithKubeConfigFile(p.GetKubeconfig())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	err = p.WaitForControlPlane(context.Background(), c)
	if err != nil {
		t.Fatalf("failed to wait for control plane: %v", err)
	}
	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	err = c.Resources().Create(context.Background(), &ns)
	if err != nil {
		t.Fatalf("failed to create namespace: %v", err)
	}
}
