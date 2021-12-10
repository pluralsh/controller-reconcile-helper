package reconcile

import (
	"context"
	"reflect"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-logr/logr"

	// istioNetworking "istio.io/api/networking/v1beta1"
	istioNetworkingClient "istio.io/client-go/pkg/apis/networking/v1beta1"
	// istioSecurity "istio.io/api/security/v1beta1"
	platformv1alpha1 "github.com/pluralsh/kubeflow-controller/apis/platform/v1alpha1"
	postgresv1 "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	istioSecurityClient "istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	hncv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
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
	if !justCreated && CopyDeploymentSetFields(deployment, foundDeployment, log) {
		log.Info("Updating Deployment", "namespace", deployment.Namespace, "name", deployment.Name)
		if err := r.Update(ctx, foundDeployment); err != nil {
			log.Error(err, "Unable to update Deployment")
			return err
		}
	}

	return nil
}

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

// VirtualService reconciles an Istio virtual service object.
func VirtualService(ctx context.Context, r client.Client, virtualservice *istioNetworkingClient.VirtualService, log logr.Logger) error {
	foundVirtualService := &istioNetworkingClient.VirtualService{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: virtualservice.Name, Namespace: virtualservice.Namespace}, foundVirtualService); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating VirtualService", "namespace", virtualservice.Namespace, "name", virtualservice.Name)
			if err := r.Create(ctx, virtualservice); err != nil {
				log.Error(err, "Unable to create VirtualService")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting VirtualService")
			return err
		}
	}
	if !justCreated && CopyVirtualService(virtualservice, foundVirtualService, log) {
		log.Info("Updating VirtualService", "namespace", virtualservice.Namespace, "name", virtualservice.Name)
		if err := r.Update(ctx, foundVirtualService); err != nil {
			log.Error(err, "Unable to update VirtualService")
			return err
		}
	}
	return nil
}

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

// RequestAuthentication reconciles an Istio RequestAuthentication object.
func RequestAuthentication(ctx context.Context, r client.Client, requestAuth *istioSecurityClient.RequestAuthentication, log logr.Logger) error {
	foundRequestAuth := &istioSecurityClient.RequestAuthentication{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: requestAuth.Name, Namespace: requestAuth.Namespace}, foundRequestAuth); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Istio RequestAuthentication", "namespace", requestAuth.Namespace, "name", requestAuth.Name)
			if err = r.Create(ctx, requestAuth); err != nil {
				log.Error(err, "Unable to create Istio RequestAuthentication")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting Istio RequestAuthentication")
			return err
		}
	}
	if !justCreated && CopyRequestAuthentication(requestAuth, foundRequestAuth, log) {
		log.Info("Updating Istio RequestAuthentication", "namespace", requestAuth.Namespace, "name", requestAuth.Name)
		if err := r.Update(ctx, foundRequestAuth); err != nil {
			log.Error(err, "Unable to update Istio RequestAuthentication")
			return err
		}
	}

	return nil
}

// AuthorizationPolicy reconciles an Istio AuthorizationPolicy object.
func AuthorizationPolicy(ctx context.Context, r client.Client, authPolicy *istioSecurityClient.AuthorizationPolicy, log logr.Logger) error {
	foundAuthPolicy := &istioSecurityClient.AuthorizationPolicy{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: authPolicy.Name, Namespace: authPolicy.Namespace}, foundAuthPolicy); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Istio AuthorizationPolicy", "namespace", authPolicy.Namespace, "name", authPolicy.Name)
			if err = r.Create(ctx, authPolicy); err != nil {
				log.Error(err, "Unable to create Istio AuthorizationPolicy")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting Istio AuthorizationPolicy")
			return err
		}
	}
	if !justCreated && CopyAuthorizationPolicy(authPolicy, foundAuthPolicy, log) {
		log.Info("Updating Istio AuthorizationPolicy", "namespace", authPolicy.Namespace, "name", authPolicy.Name)
		if err := r.Update(ctx, foundAuthPolicy); err != nil {
			log.Error(err, "Unable to update Istio AuthorizationPolicy")
			return err
		}
	}

	return nil
}

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

// RoleBinding reconciles a Role Binding object.
func RoleBinding(ctx context.Context, r client.Client, roleBinding *rbacv1.RoleBinding, log logr.Logger) error {
	foundRoleBinding := &rbacv1.RoleBinding{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, foundRoleBinding); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating RoleBinding", "namespace", roleBinding.Namespace, "name", roleBinding.Name)
			if err = r.Create(ctx, roleBinding); err != nil {
				log.Error(err, "Unable to create RoleBinding")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting RoleBinding")
			return err
		}
	}
	if !justCreated && CopyRoleBinding(roleBinding, foundRoleBinding, log) {
		log.Info("Updating RoleBinding", "namespace", roleBinding.Namespace, "name", roleBinding.Name)
		if err := r.Update(ctx, foundRoleBinding); err != nil {
			log.Error(err, "Unable to update RoleBinding")
			return err
		}
	}

	return nil
}

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

// kubeflowEnvironment reconciles a Kubeflow Environment object.
func KubeflowEnvironment(ctx context.Context, r client.Client, environment *platformv1alpha1.Environment, log logr.Logger) error {
	foundEnvironment := &platformv1alpha1.Environment{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: environment.Name, Namespace: environment.Namespace}, foundEnvironment); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating KubeflowEnvironment", "namespace", environment.Namespace, "name", environment.Name)
			if err = r.Create(ctx, environment); err != nil {
				log.Error(err, "Unable to create KubeflowEnvironment")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting KubeflowEnvironment")
			return err
		}
	}
	if !justCreated && CopyKubeflowEnvironment(environment, foundEnvironment, log) {
		log.Info("Updating KubeflowEnvironmentg", "namespace", environment.Namespace, "name", environment.Name)
		if err := r.Update(ctx, foundEnvironment); err != nil {
			log.Error(err, "Unable to update KubeflowEnvironment")
			return err
		}
	}

	return nil
}

// SubnamespaceAnchor reconciles a HNC Subnamespace object.
func SubnamespaceAnchor(ctx context.Context, r client.Client, userEnv *hncv1alpha2.SubnamespaceAnchor, log logr.Logger) error {
	foundUserEnv := &hncv1alpha2.SubnamespaceAnchor{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: userEnv.Name, Namespace: userEnv.Namespace}, foundUserEnv); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating SubnamespaceAnchor", "namespace", userEnv.Namespace, "name", userEnv.Name)
			if err = r.Create(ctx, userEnv); err != nil {
				log.Error(err, "Unable to create SubnamespaceAnchor")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting SubnamespaceAnchor")
			return err
		}
	}
	if !justCreated && CopySubnamespaceAnchor(userEnv, foundUserEnv, log) {
		log.Info("Updating SubnamespaceAnchor", "namespace", userEnv.Namespace, "name", userEnv.Name)
		if err := r.Update(ctx, foundUserEnv); err != nil {
			log.Error(err, "Unable to update SubnamespaceAnchor")
			return err
		}
	}

	return nil
}

// Postgresql reconciles a Postgresql Database object.
func Postgresql(ctx context.Context, r client.Client, postgres *postgresv1.Postgresql, log logr.Logger) error {
	foundPostgresql := &postgresv1.Postgresql{}
	justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: postgres.Name, Namespace: postgres.Namespace}, foundPostgresql); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating PostgreSQL Database", "namespace", postgres.Namespace, "name", postgres.Name)
			if err = r.Create(ctx, postgres); err != nil {
				log.Error(err, "Unable to create PostgreSQL Database")
				return err
			}
			justCreated = true
		} else {
			log.Error(err, "Error getting PostgreSQL Database")
			return err
		}
	}
	if !justCreated && CopyPostgresql(postgres, foundPostgresql, log) {
		log.Info("Updating PostgreSQL Database", "namespace", postgres.Namespace, "name", postgres.Name)
		if err := r.Update(ctx, foundPostgresql); err != nil {
			log.Error(err, "Unable to update PostgreSQL Database")
			return err
		}
	}

	return nil
}

// Reference: https://github.com/pwittrock/kubebuilder-workshop/blob/master/pkg/util/util.go

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

	if !reflect.DeepEqual(to.Spec.Replicas, to.Spec.Replicas) {
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

func CopyDeploymentSetFields(from, to *appsv1.Deployment, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling deployment due to label change")
			log.V(2).Info("difference in deployment labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling deployment due to label change")
		log.V(2).Info("difference in deployment labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling deployment due to annotation change")
			log.V(2).Info("difference in deployment annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling deployment due to annotation change")
		log.V(2).Info("difference in deployment annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec.Replicas, from.Spec.Replicas) {
		log.V(1).Info("reconciling deployment due to replica change")
		log.V(2).Info("difference in deployment replicas", "wanted", from.Spec.Replicas, "existing", to.Spec.Replicas)
		requireUpdate = true
	}
	to.Spec.Replicas = from.Spec.Replicas

	if !reflect.DeepEqual(to.Spec.Template.Spec, from.Spec.Template.Spec) {
		log.V(1).Info("reconciling deployment due to template spec change")
		log.V(2).Info("difference in deployment template spec", "wanted", from.Spec.Template.Spec, "existing", to.Spec.Template.Spec)
		requireUpdate = true
	}
	to.Spec.Template.Spec = from.Spec.Template.Spec

	return requireUpdate
}

// CopyServiceFields copies the owned fields from one Service to another
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

	if !reflect.DeepEqual(to.Spec.Ports, from.Spec.Ports) {
		log.V(1).Info("reconciling service due to ports change")
		log.V(2).Info("difference in service ports", "wanted", from.Spec.Ports, "existing", to.Spec.Ports)
		requireUpdate = true
	}
	to.Spec.Ports = from.Spec.Ports

	return requireUpdate
}

// Copy configuration related fields to another instance and returns true if there
// is a diff and thus needs to update.
func CopyVirtualService(from, to *istioNetworkingClient.VirtualService, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling virtualservice due to label change")
			log.V(2).Info("difference in virtualservice labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling virtualservice due to label change")
		log.V(2).Info("difference in virtualservice labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling virtualservice due to annotations change")
			log.V(2).Info("difference in virtualservice annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling virtualservice due to annotations change")
		log.V(2).Info("difference in virtualservice annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec, from.Spec) {
		log.V(1).Info("reconciling virtualservice due to spec change")
		log.V(2).Info("difference in virtualservice spec", "wanted", from.Spec, "exising", to.Spec)
		requireUpdate = true
	}
	to.Spec = from.Spec

	return requireUpdate
}

// CopyAuthorizationPolicy copies the owned fields from one AuthorizationPolicy to another
func CopyAuthorizationPolicy(from, to *istioSecurityClient.AuthorizationPolicy, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling AuthorizationPolicy due to label change")
			log.V(2).Info("difference in AuthorizationPolicy labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling AuthorizationPolicy due to label change")
		log.V(2).Info("difference in AuthorizationPolicy labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling AuthorizationPolicy due to annotation change")
			log.V(2).Info("difference in AuthorizationPolicy annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling AuthorizationPolicy due to annotation change")
		log.V(2).Info("difference in AuthorizationPolicy annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we this can lead to unnecessary reconciles

	if !reflect.DeepEqual(to.Spec.Selector, from.Spec.Selector) {
		log.V(1).Info("reconciling AuthorizationPolicy due to selector change")
		log.V(2).Info("difference in AuthorizationPolicy selector", "wanted", from.Spec.Selector, "existing", to.Spec.Selector)
		requireUpdate = true
	}
	to.Spec.Selector = from.Spec.Selector

	if !reflect.DeepEqual(to.Spec.Action, from.Spec.Action) {
		log.V(1).Info("reconciling AuthorizationPolicy due to action change")
		log.V(2).Info("difference in AuthorizationPolicy action", "wanted", from.Spec.Action, "existing", to.Spec.Action)
		requireUpdate = true
	}
	to.Spec.Action = from.Spec.Action

	if !reflect.DeepEqual(to.Spec.Rules, from.Spec.Rules) {
		log.V(1).Info("reconciling AuthorizationPolicy due to rule change")
		log.V(2).Info("difference in AuthorizationPolicy rules", "wanted", from.Spec.Rules, "existing", to.Spec.Rules)
		requireUpdate = true
	}
	to.Spec.Rules = from.Spec.Rules

	return requireUpdate
}

// CopyRequestAuthentication copies the owned fields from one RequestAuthentication to another
func CopyRequestAuthentication(from, to *istioSecurityClient.RequestAuthentication, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling RequestAuthentication due to label change")
			log.V(2).Info("difference in RequestAuthentication labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling RequestAuthentication due to label change")
		log.V(2).Info("difference in RequestAuthentication labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling RequestAuthentication due to annotation change")
			log.V(2).Info("difference in RequestAuthentication annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling RequestAuthentication due to annotation change")
		log.V(2).Info("difference in RequestAuthentication annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we this can lead to unnecessary reconciles
	if !reflect.DeepEqual(to.Spec.Selector, from.Spec.Selector) {
		log.V(1).Info("reconciling RequestAuthentication due to selector change")
		log.V(2).Info("difference in RequestAuthentication selector", "wanted", from.Spec.Selector, "existing", to.Spec.Selector)
		requireUpdate = true
	}
	to.Spec.Selector = from.Spec.Selector

	if !reflect.DeepEqual(to.Spec.JwtRules, from.Spec.JwtRules) {
		log.V(1).Info("reconciling RequestAuthentication due to JwtRule change")
		log.V(2).Info("difference in RequestAuthentication JwtRules", "wanted", from.Spec.JwtRules, "existing", to.Spec.JwtRules)
		requireUpdate = true
	}
	to.Spec.JwtRules = from.Spec.JwtRules

	return requireUpdate
}

// CopyServiceAccount copies the owned fields from one Service Account to another
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

// CopyRoleBinding copies the owned fields from one Role Binding to another
func CopyRoleBinding(from, to *rbacv1.RoleBinding, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling RoleBinding due to label change")
			log.V(2).Info("difference in RoleBinding labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling RoleBinding due to label change")
		log.V(2).Info("difference in RoleBinding labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling RoleBinding due to annotation change")
			log.V(2).Info("difference in RoleBinding annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling RoleBinding due to annotation change")
		log.V(2).Info("difference in RoleBinding annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we this will lead to unnecessary reconciles
	if !reflect.DeepEqual(to.RoleRef, from.RoleRef) {
		log.V(1).Info("reconciling RoleBinding due to RoleRef change")
		log.V(2).Info("difference in RoleBinding RoleRef", "wanted", from.RoleRef, "existing", to.RoleRef)
		requireUpdate = true
	}
	to.RoleRef = from.RoleRef

	if !reflect.DeepEqual(to.Subjects, from.Subjects) {
		log.V(1).Info("reconciling RoleBinding due to Subject change")
		log.V(2).Info("difference in RoleBinding Subjects", "wanted", from.Subjects, "existing", to.Subjects)
		requireUpdate = true
	}
	to.Subjects = from.Subjects

	return requireUpdate
}

// CopyNetworkPolicy copies the owned fields from one NetworkPolicy to another
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

// CopyNamespace copies the owned fields from one Namespace to another
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

// CopyPersistentVolumeClaim copies the owned fields from one PersistentVolumeClaim to another
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

	for k, v := range to.Annotations {
		if from.Annotations[k] != v && k != "pv.kubernetes.io/bind-completed" {
			log.V(1).Info("reconciling PersistentVolumeClaim due to annotation change")
			log.V(2).Info("difference in PersistentVolumeClaim annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling PersistentVolumeClaim due to annotation change")
		log.V(2).Info("difference in PersistentVolumeClaim annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec.Resources.Requests, from.Spec.Resources.Requests) {
		log.V(1).Info("reconciling PersistentVolumeClaim due to resource requests change")
		log.V(2).Info("difference in PersistentVolumeClaim resource requests", "wanted", from.Spec.Resources.Requests, "existing", to.Spec.Resources.Requests)
		requireUpdate = true
	}
	to.Spec.Resources.Requests = from.Spec.Resources.Requests

	return requireUpdate
}

// CopyPostgresql copies the owned fields from one Postgres instance to another
func CopyPostgresql(from, to *postgresv1.Postgresql, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling PostgreSQL Database due to label change")
			log.V(2).Info("difference in PostgreSQL Database labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling PostgreSQL Database due to label change")
		log.V(2).Info("difference in PostgreSQL Database labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling PostgreSQL Database due to annotation change")
			log.V(2).Info("difference in PostgreSQL Database annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling PostgreSQL Database due to annotation change")
		log.V(2).Info("difference in PostgreSQL Database annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec, from.Spec) {
		log.V(1).Info("reconciling PostgreSQL Database due to spec change")
		log.V(2).Info("difference in PostgreSQL Database spec", "wanted", from.Spec, "existing", to.Spec)
		requireUpdate = true
	}
	to.Spec = from.Spec

	return requireUpdate
}

// CopyKubeflowEnvironment copies the owned fields from one Kubeflow Environment to another
func CopyKubeflowEnvironment(from, to *platformv1alpha1.Environment, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling KubeflowEnvironment due to label change")
			log.V(2).Info("difference in KubeflowEnvironment labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling KubeflowEnvironment due to label change")
		log.V(2).Info("difference in KubeflowEnvironment labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling KubeflowEnvironment due to annotations change")
			log.V(2).Info("difference in KubeflowEnvironment annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling KubeflowEnvironment due to annotations change")
		log.V(2).Info("difference in KubeflowEnvironment annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	if !reflect.DeepEqual(to.Spec, from.Spec) {
		log.V(1).Info("reconciling KubeflowEnvironment due to spec change")
		log.V(2).Info("difference in KubeflowEnvironment spec", "wanted", from.Spec, "existing", to.Spec)
		requireUpdate = true
	}
	to.Spec = from.Spec

	return requireUpdate
}

// CopySubnamespaceAnchor copies the owned fields from one Subnamespace to another
func CopySubnamespaceAnchor(from, to *hncv1alpha2.SubnamespaceAnchor, log logr.Logger) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			log.V(1).Info("reconciling SubnamespaceAnchor due to label change")
			log.V(2).Info("difference in SubnamespaceAnchor labels", "wanted", from.Labels, "existing", to.Labels)
			requireUpdate = true
		}
	}
	if len(to.Labels) == 0 && len(from.Labels) != 0 {
		log.V(1).Info("reconciling SubnamespaceAnchor due to label change")
		log.V(2).Info("difference in SubnamespaceAnchor labels", "wanted", from.Labels, "existing", to.Labels)
		requireUpdate = true
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			log.V(1).Info("reconciling SubnamespaceAnchor due to annotation change")
			log.V(2).Info("difference in SubnamespaceAnchor annotations", "wanted", from.Annotations, "existing", to.Annotations)
			requireUpdate = true
		}
	}
	if len(to.Annotations) == 0 && len(from.Annotations) != 0 {
		log.V(1).Info("reconciling SubnamespaceAnchor due to annotation change")
		log.V(2).Info("difference in SubnamespaceAnchor annotations", "wanted", from.Annotations, "existing", to.Annotations)
		requireUpdate = true
	}
	to.Annotations = from.Annotations

	return requireUpdate
}
