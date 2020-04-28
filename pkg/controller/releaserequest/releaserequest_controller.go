package releaserequest

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	deployv1alpha1 "github.com/kube-sailmaker/request-operator/pkg/apis/deploy/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	_ "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_releaserequest")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ReleaseRequest Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileReleaseRequest{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("releaserequest-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ReleaseRequest
	err = c.Watch(&source.Kind{Type: &deployv1alpha1.ReleaseRequest{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ReleaseRequest
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &deployv1alpha1.ReleaseRequest{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileReleaseRequest implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileReleaseRequest{}

// ReconcileReleaseRequest reconciles a ReleaseRequest object
type ReconcileReleaseRequest struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ReleaseRequest object and makes changes based on the state read
// and what is in the ReleaseRequest.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileReleaseRequest) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ReleaseRequest")

	// Fetch the ReleaseRequest instance
	instance := &deployv1alpha1.ReleaseRequest{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Status.Code != "" {
		reqLogger.Info("it is not new, so skipping")
		return reconcile.Result{Requeue: false}, nil
	}

	reqLogger.Info("Processing this new CRD ", "Instance Name", instance.Name)

	for _, app := range instance.Spec.Apps {
		reqLogger.Info("Queue Job to handle these all", "App Name", app.Name, "App Version", app.Version)
	}

	// Define a new Pod object
	job := newJobForCR(instance)

	// Set ReleaseRequest instance as the owner and controller
	//if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}

	// Check if this Job already exists
	found := &batchv1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.client.Create(context.TODO(), job)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Job created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Job already exists - don't requeue
	reqLogger.Info("Skip reconcile: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
	return reconcile.Result{}, nil
}

// newJobForCR returns a busybox pod with the same name/namespace as the cr
func newJobForCR(cr *deployv1alpha1.ReleaseRequest) *batchv1.Job {
	labels := map[string]string{
		"app": cr.Name,
	}
	backoffLimit := int32(1)
	deadlineSeconds:= int64(3600)
	cleanupAfterSeconds := int32(300)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-release-job",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			ActiveDeadlineSeconds: &deadlineSeconds,
			Selector: &metav1.LabelSelector{

			},
			TTLSecondsAfterFinished: &cleanupAfterSeconds,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name: "busybox",
							Image: "busybox",
							Command: []string{"sleep", "60"},
						},
					},
				},
			},
		},
	}
}
