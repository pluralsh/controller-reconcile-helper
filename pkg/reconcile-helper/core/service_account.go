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

// ServiceAccount reconciles a Service Account object.
func ServiceAccount(ctx context.Context, r client.Client, serviceAccount *corev1.ServiceAccount, log logr.Logger) error {
	foundServiceAccount := &corev1.ServiceAccount{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: serviceAccount.Name, Namespace: serviceAccount.Namespace}, foundServiceAccount); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating ServiceAccount", "namespace", serviceAccount.Namespace, "name", serviceAccount.Name)
			if err = r.Create(ctx, serviceAccount); err != nil {
				log.Error(err, "Unable to create ServiceAccount")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting ServiceAccount")
			return err
		}
	}
	if !justCreated && CopyServiceAccount(serviceAccount, foundServiceAccount, log) {
		log.Info("Updating ServiceAccount", "namespace", serviceAccount.Namespace, "name", serviceAccount.Name)
		if err := r.Update(ctx, foundServiceAccount); err != nil {
			log.Error(err, "Unable to update ServiceAccount")
			return err
		}
	}

	return nil
}

// CopyServiceAccount copies the owned fields from one Service Account to another
// Returns true if the fields copied from don't match to.
func CopyServiceAccount(from, to *corev1.ServiceAccount, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling ServiceAccount due to label change")
			log.V(2).Info("difference in ServiceAccount labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling ServiceAccount due to label change")
		log.V(2).Info("difference in ServiceAccount labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling ServiceAccount due to annotation change")
			log.V(2).Info("difference in ServiceAccount annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling ServiceAccount due to label change")
		log.V(2).Info("difference in ServiceAccount labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we this will lead to unnecessary reconciles
	if !reflect.DeepEqual(to.ImagePullSecrets, from.ImagePullSecrets) {
		log.V(1).Info("reconciling ServiceAccount due to ImagePullSecrets change")
		log.V(2).Info("difference in ServiceAccount ImagePullSecrets", "wanted", from.ImagePullSecrets, "existing", to.ImagePullSecrets)
		requireUpdate = true
	}
	to.ImagePullSecrets = from.ImagePullSecrets

	if !reflect.DeepEqual(to.AutomountServiceAccountToken, from.AutomountServiceAccountToken) {
		log.V(1).Info("reconciling ServiceAccount due to AutomountServiceAccountToken change")
		log.V(2).Info("difference in ServiceAccount AutomountServiceAccountToken", "wanted", from.AutomountServiceAccountToken, "existing", to.AutomountServiceAccountToken)
		requireUpdate = true
	}
	to.AutomountServiceAccountToken = from.AutomountServiceAccountToken

	return requireUpdate
}
