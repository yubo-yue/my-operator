/*
Copyright 2026.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/yubo-yue/my-operator/first-operator/api/v1alpha1"
)

// FirstAppReconciler reconciles a FirstApp object
type FirstAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.example.com,resources=firstapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=firstapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.example.com,resources=firstapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FirstAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the FirstApp instance
	firstApp := &appsv1alpha1.FirstApp{}
	err := r.Get(ctx, req.NamespacedName, firstApp)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("FirstApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get FirstApp")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: firstApp.Name, Namespace: firstApp.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForFirstApp(firstApp)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := firstApp.Spec.Size
	if size == 0 {
		size = 1 // default size
	}
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the FirstApp status with the pod names
	// List the pods for this firstApp's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(firstApp.Namespace),
		client.MatchingLabels(labelsForFirstApp(firstApp.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "FirstApp.Namespace", firstApp.Namespace, "FirstApp.Name", firstApp.Name)
		return ctrl.Result{}, err
	}

	// Update status.AvailableReplicas if needed
	availableReplicas := int32(0)
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			availableReplicas++
		}
	}

	if firstApp.Status.AvailableReplicas != availableReplicas {
		firstApp.Status.AvailableReplicas = availableReplicas
		err := r.Status().Update(ctx, firstApp)
		if err != nil {
			log.Error(err, "Failed to update FirstApp status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForFirstApp returns a FirstApp Deployment object
func (r *FirstAppReconciler) deploymentForFirstApp(m *appsv1alpha1.FirstApp) *appsv1.Deployment {
	ls := labelsForFirstApp(m.Name)
	replicas := m.Spec.Size
	if replicas == 0 {
		replicas = 1
	}

	image := m.Spec.Image
	if image == "" {
		image = "nginx:latest"
	}

	port := m.Spec.Port
	if port == 0 {
		port = 80 // default port
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "firstapp",
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
	// Set FirstApp instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForFirstApp returns the labels for selecting the resources
// belonging to the given firstApp CR name.
func labelsForFirstApp(name string) map[string]string {
	return map[string]string{"app": "firstapp", "firstapp_cr": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *FirstAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.FirstApp{}).
		Owns(&appsv1.Deployment{}).
		Named("firstapp").
		Complete(r)
}
