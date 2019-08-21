package platform

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// GetInfraStatus queries the k8s api for the infrastructure CR and retrieves the current status.
func GetInfraStatus(m manager.Manager) (*configv1.InfrastructureStatus, error) {
	c, err := getClient()
	if err != nil {
		return nil, err
	}
	infraName := types.NamespacedName{Name: "cluster"}
	infra := &configv1.Infrastructure{}
	err = c.Get(context.Background(), infraName, infra)
	if err != nil {
		return nil, err
	}
	return &infra.Status, nil
}

// GetType returns the platform type given an infrastructure status. If PlatformStatus is set,
// it will get the platform type from it, otherwise it will get it from InfraStatus.Platform which
// is deprecated in 4.2
func GetType(infraStatus *configv1.InfrastructureStatus) configv1.PlatformType {
	if infraStatus.PlatformStatus != nil && len(infraStatus.PlatformStatus.Type) > 0 {
		return infraStatus.PlatformStatus.Type
	}
	return infraStatus.Platform
}

func getClient() (client.Client, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	cfg, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	//apis.AddToScheme(scheme.Scheme)
	dynamicClient, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, err
	}

	return dynamicClient, nil
}