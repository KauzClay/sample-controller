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

	contourclient "knative.dev/net-contour/pkg/client/injection/client"
	proxyinformer "knative.dev/net-contour/pkg/client/injection/informers/projectcontour/v1/httpproxy"
	kingressclient "knative.dev/networking/pkg/client/injection/client"
	kingressinformer "knative.dev/networking/pkg/client/injection/informers/networking/v1alpha1/ingress"
	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
	inginformer "knative.dev/pkg/client/injection/kube/informers/networking/v1/ingress"

	//hostruleinformer "github.com/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1beta1/informers/externalversions/ako/v1beta1"

	"knative.dev/net-contour/pkg/reconciler/contour/config"
	"knative.dev/networking/pkg/apis/networking"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	networkcfg "knative.dev/networking/pkg/config"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"

	"k8s.io/client-go/tools/cache"
)

// NewController returns a new Ingress controller for Project Contour.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)

	kingressInformer := kingressinformer.Get(ctx)
	proxyInformer := proxyinformer.Get(ctx)
	ingInformer := inginformer.Get(ctx)
	//hostRuleInformer := hostruleinformer.Informer

	c := &Reconciler{
		kingressClient: kingressclient.Get(ctx),
		contourClient:  contourclient.Get(ctx),
		contourLister:  proxyInformer.Lister(),
		kingressLister: kingressInformer.Lister(),
		ingLister:      ingInformer.Lister(),
	}
	myFilterFunc := reconciler.AnnotationFilterFunc(networking.IngressClassAnnotationKey, ContourIngressClassName, false)
	impl := kingressreconciler.NewImpl(ctx, c, ContourIngressClassName,
		func(impl *controller.Impl) controller.Options {
			configsToResync := []interface{}{
				&config.Contour{},
				&networkcfg.Config{},
			}

			resyncIngressesOnConfigChange := configmap.TypeFilter(configsToResync...)(func(string, interface{}) {
				impl.FilteredGlobalResync(myFilterFunc, kingressInformer.Informer())
			})
			configStore := config.NewStore(logger.Named("config-store"), resyncIngressesOnConfigChange)
			configStore.WatchConfigs(cmw)
			return controller.Options{
				ConfigStore:       configStore,
				PromoteFilterFunc: myFilterFunc,
			}
		})

	kingressInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: myFilterFunc,
		Handler:    controller.HandleAll(impl.Enqueue),
	})

	// Enqueue us if any of our children kingress resources change.
	kingressInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&v1alpha1.Ingress{}),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	proxyInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&v1alpha1.Ingress{}),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
