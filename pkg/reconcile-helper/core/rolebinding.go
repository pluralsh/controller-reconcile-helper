package core

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

// CopyRoleBinding copies the owned fields from one Role Binding to another
// Returns true if the fields copied from don't match to.
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
