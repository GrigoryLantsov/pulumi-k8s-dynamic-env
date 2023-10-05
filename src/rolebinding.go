package main

import (
	rbac "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/rbac/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateNewRoleBindings should be invoked reflective of your implementation's specifics.
func CreateNewRoleBindings(ctx *pulumi.Context, namespace pulumi.StringPtrOutput) error {
	// Create a RoleBinding for the view role
	_, err := rbac.NewRoleBinding(ctx, "RoleBinding_view", &rbac.RoleBindingArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(DeployKubeNamespace + "-ro-view"),
			Namespace: namespace,
		},
		RoleRef: &rbac.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("view"),
		},
		Subjects: rbac.SubjectArray{
			&rbac.SubjectArgs{
				ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
				Kind:     pulumi.String("Group"),
				Name:     pulumi.String(DeployKubeNamespace + "-ro"),
			},
		},
	})
	if err != nil {
		return err
	}

	// Create a RoleBinding for the edit role
	_, err = rbac.NewRoleBinding(ctx, "RoleBinding_edit", &rbac.RoleBindingArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(DeployKubeNamespace + "-rw-edit"),
			Namespace: namespace,
		},
		RoleRef: &rbac.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("edit"),
		},
		Subjects: rbac.SubjectArray{
			&rbac.SubjectArgs{
				ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
				Kind:     pulumi.String("Group"),
				Name:     pulumi.String(DeployKubeNamespace + "-rw"),
			},
		},
	})
	if err != nil {
		return err
	}

	// Create a RoleBinding for the admin role
	_, err = rbac.NewRoleBinding(ctx, "RoleBinding_admin", &rbac.RoleBindingArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(DeployKubeNamespace + "-admins-admin"),
			Namespace: namespace,
		},
		RoleRef: &rbac.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("admin"),
		},
		Subjects: rbac.SubjectArray{
			&rbac.SubjectArgs{
				ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
				Kind:     pulumi.String("Group"),
				Name:     pulumi.String(DeployKubeNamespace + "-admins"),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
