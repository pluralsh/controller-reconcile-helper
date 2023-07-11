package core

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Statefulset reconciles a k8s statefulset object.
func StatefulSet(ctx context.Context, r client.Client, statefulset *appsv1.StatefulSet, log logr.Logger) error {
	foundStatefulset := &appsv1.StatefulSet{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: statefulset.Name, Namespace: statefulset.Namespace}, foundStatefulset); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating StatefulSet", "namespace", statefulset.Namespace, "name", statefulset.Name)
			if err := r.Create(ctx, statefulset); err != nil {
				log.Error(err, "Unable to create StatefulSet")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting StatefulSet")
			return err
		}
	}
	if !justCreated && CopyStatefulSetFields(statefulset, foundStatefulset, log) {
		log.Info("Updating StatefulSet", "namespace", statefulset.Namespace, "name", statefulset.Name)
		if err := r.Update(ctx, foundStatefulset); err != nil {
			log.Error(err, "Unable to update StatefulSet")
			return err
		}
	}

	return nil
}

// CopyStatefulSetFields copies the owned fields from one StatefulSet to another
// Returns true if the fields copied from don't match to.
func CopyStatefulSetFields(from, to *appsv1.StatefulSet, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling StatefulSet due to label change")
			log.V(2).Info("difference in StatefulSet labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling StatefulSet due to label change")
		log.V(2).Info("difference in StatefulSet labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling StatefulSet due to annotation change")
			log.V(2).Info("difference in StatefulSet annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling StatefulSet due to annotation change")
		log.V(2).Info("difference in StatefulSet annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec.Replicas, from.Spec.Replicas) {
		log.V(1).Info("reconciling StatefulSet due to replica change")
		log.V(2).Info("difference in StatefulSet replicas", "wanted", from.Spec.Replicas, "existing", to.Spec.Replicas)
		requireUpdate = true
	}
	to.Spec.Replicas = from.Spec.Replicas

	for k, v := range to.Spec.Template.Labels {
		if from.Spec.Template.Labels[k] != v {
			log.V(1).Info("reconciling StatefulSet due to template label change")
			log.V(2).Info("difference in StatefulSet template labels", "wanted", from.Spec.Template.Labels, "existing", to.Spec.Template.Labels)
			requireUpdate = true
		}
	}
	if len(to.Spec.Template.Labels) == 0 && len(from.Spec.Template.Labels) != 0 {
		log.V(1).Info("reconciling StatefulSet due to template label change")
		log.V(2).Info("difference in StatefulSet template labels", "wanted", from.Spec.Template.Labels, "existing", to.Spec.Template.Labels)
		requireUpdate = true
	}
	to.Spec.Template.Labels = from.Spec.Template.Labels

	for k, v := range to.Spec.Template.Annotations {
		if from.Spec.Template.Annotations[k] != v {
			log.V(1).Info("reconciling StatefulSet due to template annotation change")
			log.V(2).Info("difference in StatefulSet template annotations", "wanted", from.Spec.Template.Annotations, "existing", to.Spec.Template.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Spec.Template.Annotations) == 0 && len(from.Spec.Template.Annotations) != 0 {
		log.V(1).Info("reconciling StatefulSet due to template annotation change")
		log.V(2).Info("difference in StatefulSet template annotations", "wanted", from.Spec.Template.Annotations, "existing", to.Spec.Template.Annotations)
		requireUpdate = true
	}
	to.Spec.Template.Annotations = from.Spec.Template.Annotations

	if !reflect.DeepEqual(to.Spec.Template.Spec.Volumes, from.Spec.Template.Spec.Volumes) {
		log.V(1).Info("reconciling StatefulSet due to volumes change")
		log.V(2).Info("difference in StatefulSet volumes", "wanted", from.Spec.Template.Spec.Volumes, "existing", to.Spec.Template.Spec.Volumes)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Volumes = from.Spec.Template.Spec.Volumes

	if !reflect.DeepEqual(to.Spec.Template.Spec.ServiceAccountName, from.Spec.Template.Spec.ServiceAccountName) {
		log.V(1).Info("reconciling StatefulSet due to service account name change")
		log.V(2).Info("difference in StatefulSet service account name", "wanted", from.Spec.Template.Spec.ServiceAccountName, "existing", to.Spec.Template.Spec.ServiceAccountName)
		requireUpdate = true
	}
	to.Spec.Template.Spec.ServiceAccountName = from.Spec.Template.Spec.ServiceAccountName

	if !reflect.DeepEqual(to.Spec.Template.Spec.SecurityContext, from.Spec.Template.Spec.SecurityContext) {
		log.V(1).Info("reconciling StatefulSet due to security context change")
		log.V(2).Info("difference in StatefulSet security context", "wanted", from.Spec.Template.Spec.SecurityContext, "existing", to.Spec.Template.Spec.SecurityContext)
		requireUpdate = true
	}
	to.Spec.Template.Spec.SecurityContext = from.Spec.Template.Spec.SecurityContext

	if !reflect.DeepEqual(to.Spec.Template.Spec.Affinity, from.Spec.Template.Spec.Affinity) {
		log.V(1).Info("reconciling StatefulSet due to affinity change")
		log.V(2).Info("difference in StatefulSet affinity", "wanted", from.Spec.Template.Spec.Affinity, "existing", to.Spec.Template.Spec.Affinity)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Affinity = from.Spec.Template.Spec.Affinity

	if !reflect.DeepEqual(to.Spec.Template.Spec.Tolerations, from.Spec.Template.Spec.Tolerations) {
		log.V(1).Info("reconciling StatefulSet due to toleration change")
		log.V(2).Info("difference in StatefulSet tolerations", "wanted", from.Spec.Template.Spec.Tolerations, "existing", to.Spec.Template.Spec.Tolerations)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Tolerations = from.Spec.Template.Spec.Tolerations

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].Name, from.Spec.Template.Spec.Containers[0].Name) {
		log.V(1).Info("reconciling StatefulSet due to container[0] name change")
		log.V(2).Info("difference in StatefulSet container[0] name", "wanted", from.Spec.Template.Spec.Containers[0].Name, "existing", to.Spec.Template.Spec.Containers[0].Name)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].Name = from.Spec.Template.Spec.Containers[0].Name

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].Image, from.Spec.Template.Spec.Containers[0].Image) {
		log.V(1).Info("reconciling StatefulSet due to container[0] image change")
		log.V(2).Info("difference in StatefulSet container[0] image", "wanted", from.Spec.Template.Spec.Containers[0].Image, "existing", to.Spec.Template.Spec.Containers[0].Image)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].Image = from.Spec.Template.Spec.Containers[0].Image

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].WorkingDir, from.Spec.Template.Spec.Containers[0].WorkingDir) {
		log.V(1).Info("reconciling StatefulSet due to container[0] working dir change")
		log.V(2).Info("difference in StatefulSet container[0] working dir", "wanted", from.Spec.Template.Spec.Containers[0].WorkingDir, "existing", to.Spec.Template.Spec.Containers[0].WorkingDir)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].WorkingDir = from.Spec.Template.Spec.Containers[0].WorkingDir

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].Ports, from.Spec.Template.Spec.Containers[0].Ports) {
		log.V(1).Info("reconciling StatefulSet due to container[0] port change")
		log.V(2).Info("difference in StatefulSet container[0] ports", "wanted", from.Spec.Template.Spec.Containers[0].Ports, "existing", to.Spec.Template.Spec.Containers[0].Ports)

		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].Ports = from.Spec.Template.Spec.Containers[0].Ports

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].Env, from.Spec.Template.Spec.Containers[0].Env) {
		log.V(1).Info("reconciling StatefulSet due to container[0] env change")
		log.V(2).Info("difference in StatefulSet container[0] env", "wanted", from.Spec.Template.Spec.Containers[0].Env, "existing", to.Spec.Template.Spec.Containers[0].Env)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].Env = from.Spec.Template.Spec.Containers[0].Env

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].EnvFrom, from.Spec.Template.Spec.Containers[0].EnvFrom) {
		log.V(1).Info("reconciling StatefulSet due to container[0] EnvFrom change")
		log.V(2).Info("difference in StatefulSet container[0] EnvFrom", "wanted", from.Spec.Template.Spec.Containers[0].EnvFrom, "existing", to.Spec.Template.Spec.Containers[0].EnvFrom)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].EnvFrom = from.Spec.Template.Spec.Containers[0].EnvFrom

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].Resources, from.Spec.Template.Spec.Containers[0].Resources) {
		log.V(1).Info("reconciling StatefulSet due to container[0] resource change")
		log.V(2).Info("difference in StatefulSet container[0] resources", "wanted", from.Spec.Template.Spec.Containers[0].Resources, "existing", to.Spec.Template.Spec.Containers[0].Resources)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].Resources = from.Spec.Template.Spec.Containers[0].Resources

	if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[0].VolumeMounts, from.Spec.Template.Spec.Containers[0].VolumeMounts) {
		log.V(1).Info("reconciling StatefulSet due to container[0] VolumeMounts change")
		log.V(2).Info("difference in StatefulSet container[0] VolumeMounts", "wanted", from.Spec.Template.Spec.Containers[0].VolumeMounts, "existing", to.Spec.Template.Spec.Containers[0].VolumeMounts)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Containers[0].VolumeMounts = from.Spec.Template.Spec.Containers[0].VolumeMounts

	return requireUpdate
}
