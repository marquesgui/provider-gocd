package pipelineconfig

import (
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/marquesgui/provider-gocd/apis/config/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errGetSecret = "cannot get secret"
)

// GetSecretValue retrieves the value of a secret key from a Kubernetes secret.
func GetSecretValue(ctx context.Context, kube client.Client, selector *xpv1.SecretKeySelector) (string, error) {
	nn := types.NamespacedName{
		Name:      selector.Name,
		Namespace: selector.Namespace,
	}
	var secret corev1.Secret
	if err := kube.Get(ctx, nn, &secret); err != nil {
		return "", errors.Wrap(err, errGetSecret)
	}

	val, ok := secret.Data[selector.Key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s/%s", selector.Key, selector.Namespace, selector.Name)
	}

	return string(val), nil
}

// GetValueFrom retrieves the value from a given environment variable source.
func GetValueFrom(ctx context.Context, kube client.Client, from *v1alpha1.EnvVarSource) (string, bool, error) {
	if from.ConfigMapKeyRef != nil {
		nn := types.NamespacedName{
			Name:      from.ConfigMapKeyRef.Name,
			Namespace: from.ConfigMapKeyRef.Namespace,
		}
		var cm corev1.ConfigMap
		if err := kube.Get(ctx, nn, &cm); err != nil {
			return "", false, errors.Wrap(err, "cannot get configmap")
		}

		val, ok := cm.Data[from.ConfigMapKeyRef.Key]
		if !ok {
			return "", false, fmt.Errorf("key %s not found in configmap %s/%s", from.ConfigMapKeyRef.Key, from.ConfigMapKeyRef.Namespace, from.ConfigMapKeyRef.Name)
		}
		return val, false, nil
	}
	if from.SecretKeyRef != nil {
		val, err := GetSecretValue(ctx, kube, from.SecretKeyRef)
		return val, true, err
	}
	return "", false, nil
}
