/*
Copyright 2024 The Knative Authors

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

package kingress

import (
	"context"

	corev1listers "k8s.io/client-go/listers/core/v1"

	networkingv1listers "k8s.io/client-go/listers/networking/v1"
	contourclientset "knative.dev/net-contour/pkg/client/clientset/versioned"
	contourlisters "knative.dev/net-contour/pkg/client/listers/projectcontour/v1"
	kingressclientset "knative.dev/networking/pkg/client/clientset/versioned"
	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
	networkingv1alpha1 "knative.dev/networking/pkg/client/listers/networking/v1alpha1"

	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/networking/pkg/status"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
)

// Reconciler implements controller.Reconciler for Ingress resources.
type Reconciler struct {
	kingressClient kingressclientset.Interface
	contourClient  contourclientset.Interface

	// Listers index properties about resources
	contourLister  contourlisters.HTTPProxyLister
	kingressLister networkingv1alpha1.IngressLister
	serviceLister  corev1listers.ServiceLister
	ingLister      networkingv1listers.IngressLister
	//httpRouteLister avi.HTTPRouteLister

	statusManager status.Manager
	tracker       tracker.Interface
}

var _ kingressreconciler.Interface = (*Reconciler)(nil)

// ReconcileKind reconciles ingress resource.
func (r *Reconciler) ReconcileKind(ctx context.Context, ing *v1alpha1.Ingress) reconciler.Event {
	return nil
}
