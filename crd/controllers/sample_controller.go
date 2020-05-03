/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	samplev1alpha1 "crd/api/v1alpha1"
)

// SampleReconciler reconciles a Sample object
type SampleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wreulicke.github.io,resources=samples,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wreulicke.github.io,resources=samples/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete

func (r *SampleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("sample", req.NamespacedName)

	var sample samplev1alpha1.Sample
	log.Info("fetching Foo Resource")
	if err := r.Get(ctx, req.NamespacedName, &sample); err != nil {
		log.Error(err, "unable to fetch Foo")
		return ctrl.Result{}, err
	}
	deployment := apps.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      sample.Spec.DeploymentName,
	}, &deployment); err != nil {
		log.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, err
	}
	log.Info(sample.Namespace + "/" + sample.Name + " is found")
	log.Info(deployment.Namespace + "/" + deployment.Name + " is found")
	copy := deployment.DeepCopy()
	copy.ObjectMeta.Annotations["test"] = "xxx"
	if err := r.Update(ctx, copy); err != nil {
		log.Error(err, "unable to update Deployment")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SampleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplev1alpha1.Sample{}).
		Complete(r)
}
