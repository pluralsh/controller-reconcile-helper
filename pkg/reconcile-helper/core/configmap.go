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

// ConfigMap reconciles a ConfigMap object.
func ConfigMap(ctx context.Context, r client.Client, configMap *corev1.ConfigMap, log logr.Logger) error {
	foundConfigMap := &corev1.ConfigMap{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, foundConfigMap); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating ConfigMap", "namespace", configMap.Namespace, "name", configMap.Name)
			if err = r.Create(ctx, configMap); err != nil {
				log.Error(err, "Unable to create ConfigMap")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting ConfigMap")
			return err
		}
	}
	if !justCreated && CopyConfigMap(configMap, foundConfigMap, log) {
		log.Info("Updating ConfigMap", "namespace", configMap.Namespace, "name", configMap.Name)
		if err := r.Update(ctx, foundConfigMap); err != nil {
			log.Error(err, "Unable to update ConfigMap")
			return err
		}
	}

	return nil
}

// CopyConfigMap copies the owned fields from one Service Account to another
// Returns true if the fields copied from don't match to.
func CopyConfigMap(from, to *corev1.ConfigMap, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling ConfigMap due to label change")
			log.V(2).Info("difference in ConfigMap labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling ConfigMap due to label change")
		log.V(2).Info("difference in ConfigMap labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling ConfigMap due to annotation change")
			log.V(2).Info("difference in ConfigMap annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling ConfigMap due to annotation change")
		log.V(2).Info("difference in ConfigMap annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we this will lead to unnecessary reconciles
	if !reflect.DeepEqual(to.Data, from.Data) {
		log.V(1).Info("reconciling ConfigMap due to Data change")
		log.V(2).Info("difference in ConfigMap Data", "wanted", from.Data, "existing", to.Data)
		requireUpdate = true
	}
	to.Data = from.Data

	if !reflect.DeepEqual(to.BinaryData, from.BinaryData) {
		log.V(1).Info("reconciling ConfigMap due to BinaryData change")
		log.V(2).Info("difference in ConfigMap BinaryData", "wanted", from.BinaryData, "existing", to.BinaryData)
		requireUpdate = true
	}
	to.BinaryData = from.BinaryData

	return requireUpdate
}
