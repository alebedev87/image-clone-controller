package deployment

import (
	"context"

	"image-clone-controller/pkg/controller/utils"
	"image-clone-controller/pkg/registry"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("deployment-controller")

// Add creates a new deployment controller and adds it to the manager
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new deployment reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployment{
		client:    mgr.GetClient(),
		regClient: registry.NewClientFromConfig(),
	}
}

// add adds a new controller to the given manager
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("deployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	pred := utils.NewBlacklistNamespacePredicateFromConfig()
	if err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}, pred); err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDeployment{}

// ReconcileDeployment reconciles a Deployment object
type ReconcileDeployment struct {
	// split client (reads from the cache, writes to API)
	client    client.Client
	regClient *registry.Client
}

// Reconcile migrates Deployments to backed up images
func (r *ReconcileDeployment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("deployment", request.NamespacedName)
	logger.Info("Reconciling deployment")

	// fetch deployment instance
	instance := &appsv1.Deployment{}
	err := r.client.Get(context.Background(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// object was deleted - nothing to do
			return reconcile.Result{}, nil
		}
		// error getting the deployment - requeue the request
		return reconcile.Result{}, err
	}

	// checking the images
	numChangedImg, numErrorImg := 0, 0
	for i, c := range instance.Spec.Template.Spec.Containers {
		if r.regClient.Belongs(c.Image) {
			continue
		}
		logger.Info("Cloning the image", "Image", c.Image)
		if newImg, err := r.regClient.Backup(c.Image); err == nil {
			instance.Spec.Template.Spec.Containers[i].Image = newImg
			numChangedImg++
		} else {
			logger.Error(err, "Failed to clone the image", "Image", c.Image)
			numErrorImg++
			// best effort: backup as many as possible, no requeue
		}
	}

	// migrating to the new images
	if numChangedImg > 0 {
		logger.Info("Updating the deployment to backed up images", "Changed images", numChangedImg)
		r.client.Update(context.Background(), instance)
	} else if numErrorImg == 0 {
		logger.Info("Deployment is fully backed up!")
	}

	return reconcile.Result{}, nil
}
