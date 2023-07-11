package core

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Namespace reconciles a Namespace object.
func Namespace(ctx context.Context, r client.Client, namespace *corev1.Namespace, log logr.Logger) error {
	foundNamespace := &corev1.Namespace{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: namespace.Name, Namespace: namespace.Namespace}, foundNamespace); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Namespace", "namespace", namespace.Name)
			if err = r.Create(ctx, namespace); err != nil {
				// IncRequestErrorCounter("error creating namespace", SEVERITY_MAJOR)
				log.Error(err, "Unable to create Namespace")
				return err
			}
			err = backoff.Retry(
				func() error {
					return r.Get(ctx, types.NamespacedName{Name: namespace.Name}, foundNamespace)
				},
				backoff.WithMaxRetries(backoff.NewConstantBackOff(3*time.Second), 5))
			if err != nil {
				// IncRequestErrorCounter("error namespace create completion", SEVERITY_MAJOR)
				log.Error(err, "Error Namespace create completion")
				return err
				// return r.appendErrorConditionAndReturn(ctx, namespace,
				// "Owning namespace failed to create within 15 seconds")
			}
			log.Info("Created Namespace: "+foundNamespace.Name, "status", foundNamespace.Status.Phase)
			justCreated = true
		} else {
			// IncRequestErrorCounter("error getting Namespace", SEVERITY_MAJOR)
			log.Error(err, "Error getting Namespace")
			return err
		}
	}
	if !justCreated && CopyNamespace(namespace, foundNamespace, log) {
		log.Info("Updating Namespace", "namespace", namespace.Name)
		if err := r.Update(ctx, foundNamespace); err != nil {
			log.Error(err, "Unable to update Namespace")
			return err
		}
	}

	return nil
}

// CopyNamespace copies the owned fields from one Namespace to another
// Returns true if the fields copied from don't match to.
func CopyNamespace(from, to *corev1.Namespace, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling Namespace due to label change")
			log.V(2).Info("difference in Namespace labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling Namespace due to label change")
		log.V(2).Info("difference in Namespace labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling Namespace due to annotation change")
			log.V(2).Info("difference in Namespace annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling Namespace due to annotation change")
		log.V(2).Info("difference in Namespace annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	return requireUpdate
}
