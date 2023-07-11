package core

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"

	networkv1 "k8s.io/api/networking/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NetworkPolicy reconciles a NetworkPolicy object.
func NetworkPolicy(ctx context.Context, r client.Client, networkPolicy *networkv1.NetworkPolicy, log logr.Logger) error {
	foundNetworkPolicy := &networkv1.NetworkPolicy{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: networkPolicy.Name, Namespace: networkPolicy.Namespace}, foundNetworkPolicy); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating NetworkPolicy", "namespace", networkPolicy.Namespace, "name", networkPolicy.Name)
			if err = r.Create(ctx, networkPolicy); err != nil {
				log.Error(err, "Unable to create NetworkPolicy")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting NetworkPolicy")
			return err
		}
	}
	if !justCreated && CopyNetworkPolicy(networkPolicy, foundNetworkPolicy, log) {
		log.Info("Updating NetworkPolicy", "namespace", networkPolicy.Namespace, "name", networkPolicy.Name)
		if err := r.Update(ctx, foundNetworkPolicy); err != nil {
			log.Error(err, "Unable to update NetworkPolicy")
			return err
		}
	}

	return nil
}

// Reference: https://github.com/pwittrock/kubebuilder-workshop/blob/master/pkg/util/util.go

// CopyNetworkPolicy copies the owned fields from one NetworkPolicy to another
// Returns true if the fields copied from don't match to.
func CopyNetworkPolicy(from, to *networkv1.NetworkPolicy, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling NetworkPolicy due to label change")
			log.V(2).Info("difference in NetworkPolicy labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling NetworkPolicy due to label change")
		log.V(2).Info("difference in NetworkPolicy labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling NetworkPolicy due to annotation change")
			log.V(2).Info("difference in NetworkPolicy annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling NetworkPolicy due to annotation change")
		log.V(2).Info("difference in NetworkPolicy annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec, from.Spec) {
		log.V(1).Info("reconciling NetworkPolicy due to spec change")
		log.V(2).Info("difference in NetworkPolicy spec", "wanted", from.Spec, "existing", to.Spec)
		requireUpdate = true
	}
	to.Spec = from.Spec

	return requireUpdate
}
