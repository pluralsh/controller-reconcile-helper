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

// Service reconciles a k8s service object.
func Service(ctx context.Context, r client.Client, service *corev1.Service, log logr.Logger) error {
	foundService := &corev1.Service{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Service", "namespace", service.Namespace, "name", service.Name)
			if err = r.Create(ctx, service); err != nil {
				log.Error(err, "Unable to create Service")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting Service")
			return err
		}
	}
	if !justCreated && CopyServiceFields(service, foundService, log) {
		log.Info("Updating Service", "namespace", service.Namespace, "name", service.Name)
		if err := r.Update(ctx, foundService); err != nil {
			log.Error(err, "Unable to update Service")
			return err
		}
	}

	return nil
}

// CopyServiceFields copies the owned fields from one Service to another
// Returns true if the fields copied from don't match to.
func CopyServiceFields(from, to *corev1.Service, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling service due to label change")
			log.V(2).Info("difference in service labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling service due to label change")
		log.V(2).Info("difference in service labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling service due to annotation change")
			log.V(2).Info("difference in service annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling service due to annotation change")
		log.V(2).Info("difference in service annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we can't overwrite the clusterIp field
	if !reflect.DeepEqual(to.Spec.Selector, from.Spec.Selector) {
		log.V(1).Info("reconciling service due to selector change")
		log.V(2).Info("difference in service selector", "wanted", from.Spec.Selector, "existing", to.Spec.Selector)
		requireUpdate = true
	}
	to.Spec.Selector = from.Spec.Selector

	if !reflect.DeepEqual(to.Spec.Type, from.Spec.Type) {
		log.V(1).Info("reconciling service due to ports change")
		log.V(2).Info("difference in service type", "wanted", from.Spec.Type, "existing", to.Spec.Type)
		requireUpdate = true
	}
	to.Spec.Type = from.Spec.Type

	if to.Spec.Type == corev1.ServiceTypeLoadBalancer {

		if !reflect.DeepEqual(to.Spec.Ports[0].AppProtocol, from.Spec.Ports[0].AppProtocol) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports[0].AppProtocol, "existing", to.Spec.Ports[0].AppProtocol)
			requireUpdate = true
		}
		to.Spec.Ports[0].AppProtocol = from.Spec.Ports[0].AppProtocol

		if !reflect.DeepEqual(to.Spec.Ports[0].Name, from.Spec.Ports[0].Name) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports[0].Name, "existing", to.Spec.Ports[0].Name)
			requireUpdate = true
		}
		to.Spec.Ports[0].Name = from.Spec.Ports[0].Name

		if !reflect.DeepEqual(to.Spec.Ports[0].Port, from.Spec.Ports[0].Port) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports[0].Port, "existing", to.Spec.Ports[0].Port)
			requireUpdate = true
		}
		to.Spec.Ports[0].Port = from.Spec.Ports[0].Port

		if !reflect.DeepEqual(to.Spec.Ports[0].Protocol, from.Spec.Ports[0].Protocol) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports[0].Protocol, "existing", to.Spec.Ports[0].Protocol)
			requireUpdate = true
		}
		if !reflect.DeepEqual(to.Spec.Ports[0].TargetPort, from.Spec.Ports[0].TargetPort) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports[0].TargetPort, "existing", to.Spec.Ports[0].TargetPort)
			requireUpdate = true
		}
		to.Spec.Ports[0].TargetPort = from.Spec.Ports[0].TargetPort
	} else {

		if !reflect.DeepEqual(to.Spec.Ports, from.Spec.Ports) {
			log.V(1).Info("reconciling service due to ports change")
			log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports, "existing", to.Spec.Ports)
			requireUpdate = true
		}
		to.Spec.Ports = from.Spec.Ports
	}

	return requireUpdate
}
