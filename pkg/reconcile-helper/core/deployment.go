package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deployment reconciles a k8s deployment object.
func Deployment(ctx context.Context, r client.Client, deployment *appsv1.Deployment, log logr.Logger) error {
	foundDeployment := &appsv1.Deployment{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Deployment", "namespace", deployment.Namespace, "name", deployment.Name)
			if err := r.Create(ctx, deployment); err != nil {
				log.Error(err, "Unable to create Deployment")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting Deployment")
			return err
		}
	}
	if !justCreated && CopyDeploymentFields(deployment, foundDeployment, log) {
		log.Info("Updating Deployment", "namespace", deployment.Namespace, "name", deployment.Name)
		if err := r.Update(ctx, foundDeployment); err != nil {
			log.Error(err, "Unable to update Deployment")
			return err
		}
	}

	return nil
}

// CopyDeploymentFields copies fields from one deployment to another.
// Returns true if the fields copied from don't match to.
func CopyDeploymentFields(from, to *appsv1.Deployment, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling Deployment due to label change")
			log.V(2).Info("difference in Deployment labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling Deployment due to label change")
		log.V(2).Info("difference in Deployment labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	// for k, v := range to.Annotations {
	// 	if from.Annotations[k] != v {
	// 		log.V(1).Info("reconciling Deployment due to annotation change")
	// 		log.V(2).Info("difference in Deployment annotations", "wanted", from.Annotations, "existing", to.Annotations)
	// 		requireUpdate = true
	// 	}
	// }
	// if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
	// 	log.V(1).Info("reconciling Deployment due to annotation change")
	// 	log.V(2).Info("difference in Deployment annotations", "wanted", from.Annotations, "existing", to.Annotations)
	// 	requireUpdate = true
	// }
	// to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec.Replicas, from.Spec.Replicas) {
		log.V(1).Info("reconciling Deployment due to replica change")
		log.V(2).Info("difference in Deployment replicas", "wanted", from.Spec.Replicas, "existing", to.Spec.Replicas)
		requireUpdate = true
	}
	to.Spec.Replicas = from.Spec.Replicas

	for k, v := range to.Spec.Template.Labels {
		if from.Spec.Template.Labels[k] != v {
			log.V(1).Info("reconciling Deployment due to template label change")
			log.V(2).Info("difference in Deployment template labels", "wanted", from.Spec.Template.Labels, "existing", to.Spec.Template.Labels)
			requireUpdate = true
		}
	}
	if len(to.Spec.Template.Labels) == 0 && len(from.Spec.Template.Labels) != 0 {
		log.V(1).Info("reconciling Deployment due to template label change")
		log.V(2).Info("difference in Deployment template labels", "wanted", from.Spec.Template.Labels, "existing", to.Spec.Template.Labels)
		requireUpdate = true
	}
	to.Spec.Template.Labels = from.Spec.Template.Labels

	for k, v := range to.Spec.Template.Annotations {
		if from.Spec.Template.Annotations[k] != v {
			log.V(1).Info("reconciling Deployment due to template annotation change")
			log.V(2).Info("difference in Deployment template annotations", "wanted", from.Spec.Template.Annotations, "existing", to.Spec.Template.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Spec.Template.Annotations) == 0 && len(from.Spec.Template.Annotations) != 0 {
		log.V(1).Info("reconciling Deployment due to template annotation change")
		log.V(2).Info("difference in Deployment template annotations", "wanted", from.Spec.Template.Annotations, "existing", to.Spec.Template.Annotations)
		requireUpdate = true
	}
	to.Spec.Template.Annotations = from.Spec.Template.Annotations

	if !reflect.DeepEqual(to.Spec.Template.Spec.Volumes, from.Spec.Template.Spec.Volumes) {
		log.V(1).Info("reconciling Deployment due to volumes change")
		log.V(2).Info("difference in Deployment volumes", "wanted", from.Spec.Template.Spec.Volumes, "existing", to.Spec.Template.Spec.Volumes)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Volumes = from.Spec.Template.Spec.Volumes

	if !reflect.DeepEqual(to.Spec.Template.Spec.ServiceAccountName, from.Spec.Template.Spec.ServiceAccountName) {
		log.V(1).Info("reconciling Deployment due to service account name change")
		log.V(2).Info("difference in Deployment service account name", "wanted", from.Spec.Template.Spec.ServiceAccountName, "existing", to.Spec.Template.Spec.ServiceAccountName)
		requireUpdate = true
	}
	to.Spec.Template.Spec.ServiceAccountName = from.Spec.Template.Spec.ServiceAccountName

	if !reflect.DeepEqual(to.Spec.Template.Spec.SecurityContext, from.Spec.Template.Spec.SecurityContext) {
		log.V(1).Info("reconciling Deployment due to security context change")
		log.V(2).Info("difference in Deployment security context", "wanted", from.Spec.Template.Spec.SecurityContext, "existing", to.Spec.Template.Spec.SecurityContext)
		requireUpdate = true
	}
	to.Spec.Template.Spec.SecurityContext = from.Spec.Template.Spec.SecurityContext

	if !reflect.DeepEqual(to.Spec.Template.Spec.Affinity, from.Spec.Template.Spec.Affinity) {
		log.V(1).Info("reconciling Deployment due to affinity change")
		log.V(2).Info("difference in Deployment affinity", "wanted", from.Spec.Template.Spec.Affinity, "existing", to.Spec.Template.Spec.Affinity)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Affinity = from.Spec.Template.Spec.Affinity

	if !reflect.DeepEqual(to.Spec.Template.Spec.Tolerations, from.Spec.Template.Spec.Tolerations) {
		log.V(1).Info("reconciling Deployment due to toleration change")
		log.V(2).Info("difference in Deployment tolerations", "wanted", from.Spec.Template.Spec.Tolerations, "existing", to.Spec.Template.Spec.Tolerations)
		requireUpdate = true
	}
	to.Spec.Template.Spec.Tolerations = from.Spec.Template.Spec.Tolerations

	if !reflect.DeepEqual(to.Spec.Template.Spec.TopologySpreadConstraints, from.Spec.Template.Spec.TopologySpreadConstraints) {
		log.V(1).Info("reconciling Deployment due to topology spread constraints change")
		log.V(2).Info("difference in Deployment topology spread constraints", "wanted", from.Spec.Template.Spec.TopologySpreadConstraints, "existing", to.Spec.Template.Spec.TopologySpreadConstraints)
		requireUpdate = true
	}
	to.Spec.Template.Spec.TopologySpreadConstraints = from.Spec.Template.Spec.TopologySpreadConstraints

	for i := range to.Spec.Template.Spec.Containers {

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].Name, from.Spec.Template.Spec.Containers[i].Name) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] name change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container [%v] name", i), "wanted", from.Spec.Template.Spec.Containers[i].Name, "existing", to.Spec.Template.Spec.Containers[i].Name)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].Name = from.Spec.Template.Spec.Containers[i].Name

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].Image, from.Spec.Template.Spec.Containers[i].Image) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] image change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] image", i), "wanted", from.Spec.Template.Spec.Containers[i].Image, "existing", to.Spec.Template.Spec.Containers[i].Image)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].Image = from.Spec.Template.Spec.Containers[i].Image

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].WorkingDir, from.Spec.Template.Spec.Containers[i].WorkingDir) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] working dir change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] working dir", i), "wanted", from.Spec.Template.Spec.Containers[i].WorkingDir, "existing", to.Spec.Template.Spec.Containers[i].WorkingDir)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].WorkingDir = from.Spec.Template.Spec.Containers[i].WorkingDir

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].Ports, from.Spec.Template.Spec.Containers[i].Ports) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] port change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] ports", i), "wanted", from.Spec.Template.Spec.Containers[i].Ports, "existing", to.Spec.Template.Spec.Containers[i].Ports)

			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].Ports = from.Spec.Template.Spec.Containers[i].Ports

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].Env, from.Spec.Template.Spec.Containers[i].Env) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] env change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] env", i), "wanted", from.Spec.Template.Spec.Containers[i].Env, "existing", to.Spec.Template.Spec.Containers[i].Env)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].Env = from.Spec.Template.Spec.Containers[i].Env

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].EnvFrom, from.Spec.Template.Spec.Containers[i].EnvFrom) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] EnvFrom change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] EnvFrom", i), "wanted", from.Spec.Template.Spec.Containers[i].EnvFrom, "existing", to.Spec.Template.Spec.Containers[i].EnvFrom)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].EnvFrom = from.Spec.Template.Spec.Containers[i].EnvFrom

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].Resources, from.Spec.Template.Spec.Containers[i].Resources) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] resource change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] resources", i), "wanted", from.Spec.Template.Spec.Containers[i].Resources, "existing", to.Spec.Template.Spec.Containers[i].Resources)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].Resources = from.Spec.Template.Spec.Containers[i].Resources

		if !reflect.DeepEqual(to.Spec.Template.Spec.Containers[i].VolumeMounts, from.Spec.Template.Spec.Containers[i].VolumeMounts) {
			log.V(1).Info(fmt.Sprintf("reconciling Deployment due to container[%v] VolumeMounts change", i))
			log.V(2).Info(fmt.Sprintf("difference in Deployment container[%v] VolumeMounts", i), "wanted", from.Spec.Template.Spec.Containers[i].VolumeMounts, "existing", to.Spec.Template.Spec.Containers[i].VolumeMounts)
			requireUpdate = true
		}
		to.Spec.Template.Spec.Containers[i].VolumeMounts = from.Spec.Template.Spec.Containers[i].VolumeMounts
	}

	return requireUpdate
}
