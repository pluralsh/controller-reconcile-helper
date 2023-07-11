package core

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Secret reconciles a k8s secret object.
func Secret(ctx context.Context, r client.Client, secret *corev1.Secret, log logr.Logger) error {
	foundSecret := &corev1.Secret{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, foundSecret); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Secret", "namespace", secret.Namespace, "name", secret.Name)
			if err := r.Create(ctx, secret); err != nil {
				log.Error(err, "Unable to create Secret")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting Secret")
			return err
		}
	}
	if !justCreated && CopySecretFields(secret, foundSecret, log) {
		log.Info("Updating Secret", "namespace", secret.Namespace, "name", secret.Name)
		if err := r.Update(ctx, foundSecret); err != nil {
			log.Error(err, "Unable to update Secret")
			return err
		}
	}

	return nil
}

// CopySecretFields copies the owned fields from one Service to another
// Returns true if the fields copied from don't match to.
func CopySecretFields(from, to *corev1.Secret, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling secret due to label change")
			log.V(2).Info("difference in secret labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling secret due to label change")
		log.V(2).Info("difference in secret labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling secret due to annotation change")
			log.V(2).Info("difference in secret annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling secret due to annotation change")
		log.V(2).Info("difference in secret annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if to.Type != corev1.SecretTypeServiceAccountToken {
		// Don't copy the entire Spec, because we can't overwrite the clusterIp field
		if !reflect.DeepEqual(to.Data, from.Data) {
			log.V(1).Info("reconciling secret due to data change")
			log.V(2).Info("difference in secret selector", "wanted", from.Data, "existing", to.Data)
			requireUpdate = true
		}
		to.Data = from.Data
	}

	return requireUpdate
}
