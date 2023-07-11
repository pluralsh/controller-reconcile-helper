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

// PersistentVolumeClaim reconciles a k8s pvc object.
func PersistentVolumeClaim(ctx context.Context, r client.Client, pvc *corev1.PersistentVolumeClaim, log logr.Logger) error {
	foundPVC := &corev1.PersistentVolumeClaim{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, foundPVC); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating PersistentVolumeClaim", "namespace", pvc.Namespace, "name", pvc.Name)
			if err := r.Create(ctx, pvc); err != nil {
				log.Error(err, "Unable to create PersistentVolumeClaim")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting PersistentVolumeClaim")
			return err
		}
	}
	if !justCreated && CopyPersistentVolumeClaim(pvc, foundPVC, log) {
		log.Info("Updating PersistentVolumeClaim", "namespace", pvc.Namespace, "name", pvc.Name)
		if err := r.Update(ctx, foundPVC); err != nil {
			log.Error(err, "Unable to update PersistentVolumeClaim")
			return err
		}
	}

	return nil
}

// CopyPersistentVolumeClaim copies the owned fields from one PersistentVolumeClaim to another
// Returns true if the fields copied from don't match to.
func CopyPersistentVolumeClaim(from, to *corev1.PersistentVolumeClaim, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling PersistentVolumeClaim due to label change")
			log.V(2).Info("difference in PersistentVolumeClaim labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling PersistentVolumeClaim due to label change")
		log.V(2).Info("difference in PersistentVolumeClaim labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	// for k, v := range to.Annotations {
	// 	if from.Annotations[k] != v {
	// 		log.V(1).Info("reconciling PersistentVolumeClaim due to annotation change")
	// 		log.V(2).Info("difference in PersistentVolumeClaim annotations", "wanted", from.Annotations, "existing", to.Annotations)
	// 		requireUpdate = true
	// 	}
	// }
	// if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
	// 	log.V(1).Info("reconciling PersistentVolumeClaim due to annotation change")
	// 	log.V(2).Info("difference in PersistentVolumeClaim annotations", "wanted", from.Annotations, "existing", to.Annotations)
	// 	requireUpdate = true
	// }
	// to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec.Resources.Requests, from.Spec.Resources.Requests) {
		log.V(1).Info("reconciling PersistentVolumeClaim due to resource requests change")
		log.V(2).Info("difference in PersistentVolumeClaim resource requests", "wanted", from.Spec.Resources.Requests, "existing", to.Spec.Resources.Requests)
		requireUpdate = true
	}
	to.Spec.Resources.Requests = from.Spec.Resources.Requests

	return requireUpdate
}
