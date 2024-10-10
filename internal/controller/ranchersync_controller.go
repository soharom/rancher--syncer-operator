/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	ranchersynciov1alpha1 "github.com/soharom/rancher-image-sync/api/v1alpha1"
	"github.com/soharom/rancher-image-sync/internal"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	syncerLog = ctrl.Log.WithName("controller")
)

// RancherSyncReconciler reconciles a RancherSync object
type RancherSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type RancherSyncNamespaced struct {
	ranchersyncNamespace string
	ranchersyncName      string
}

// +kubebuilder:rbac:groups=rancher.sync.io,resources=ranchersyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rancher.sync.io,resources=ranchersyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rancher.sync.io,resources=ranchersyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RancherSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *RancherSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	instance := &ranchersynciov1alpha1.RancherSync{}

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			syncerLog.Info("Ranchersync  resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil

		}
		syncerLog.Error(err, "Failed to get Ranchersync")
		return ctrl.Result{}, err
	}
	if err := r.GenerateSeceretResources(instance.Spec.Token, instance.Spec.Api, instance, ctx); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *RancherSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ranchersynciov1alpha1.RancherSync{}).
		Complete(r)
}

// Function define secret and the data

func (r *RancherSyncReconciler) secretForMemcahed(
	ranchersync *ranchersynciov1alpha1.RancherSync, clusterSecretData *internal.ClusterSecretData) (*corev1.Secret, error) {
	ranchersyncNamespaced := &RancherSyncNamespaced{
		ranchersyncNamespace: ranchersync.ObjectMeta.Namespace,
		ranchersyncName:      ranchersync.ObjectMeta.Name,
	}

	secr := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ranchersyncNamespaced.ranchersyncName + "-" + clusterSecretData.Name,
			Namespace: ranchersyncNamespaced.ranchersyncNamespace,
			Labels: map[string]string{
				"creator": "ranchersync-controller",
			},
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"name":       clusterSecretData.Name,
			"kubeconfig": clusterSecretData.Kubeconfig,
		},
	}
	if err := ctrl.SetControllerReference(ranchersync, secr, r.Scheme); err != nil {
		return nil, err
	}
	return secr, nil
}

//Function

func (r *RancherSyncReconciler) GenerateSeceretResources(token string, url string, instance *ranchersynciov1alpha1.RancherSync, ctx context.Context) error {
	client := internal.NewClient(token, true)
	clusters, err := client.GetClusters(url)

	if err != nil {
		syncerLog.Error(err, "Failed to get clusters on  Rancher :", "Endpoint", url)
		return err
	}
	if clusters != nil {

		for _, cluster := range clusters.ClusterDatas {
			generatedConfig, err := client.GenerateClusterConfig(cluster.ClusterActions.GenerateKubeconfigEndpoint)
			if err != nil {
				syncerLog.Error(err, "Failed to generate cluster kubeconfig :", "Endpoint", url, "Clusterr Name", clusters.ClusterDatas[0].Id)
			}
			clusterSecretData := &internal.ClusterSecretData{
				Name:       cluster.Name,
				Kubeconfig: generatedConfig.Config,
			}
			found := &corev1.Secret{}
			err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
			if err != nil && apierrors.IsNotFound(err) {
				// Define a new Secret

				sec, err := r.secretForMemcahed(instance, clusterSecretData)
				if err != nil {
					syncerLog.Error(err, "Failed to define new secret resource for Memcached")

					// The following implementation will update the status
					r.Status().Update(ctx, instance)

					return err
				}

				syncerLog.Info("Creating a new Secret", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name)
				if err = r.Create(ctx, sec); err != nil {
					syncerLog.Error(err, "Failed to create new Secret", "Secret.Namespace", sec.Namespace, "Secret.Name", sec.Name)
					return err
				}

				return nil
			} else if err != nil {
				syncerLog.Error(err, "Failed to get Secret")
				// Let's return the error for the reconciliation be re-trigged again
				return err
			}

		}
	} else {
		syncerLog.Info("The cluster is nil:")
		return nil
	}
	return nil
}
