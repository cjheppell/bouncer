package tcpservice

import (
	"context"

	bouncerv1alpha1 "github.com/cjheppell/bouncer/pkg/apis/bouncer/v1alpha1"
	"github.com/cjheppell/bouncer/pkg/cloudprovider"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_tcpservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new TcpService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTcpService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("tcpservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource TcpService
	err = c.Watch(&source.Kind{Type: &bouncerv1alpha1.TcpService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner TcpService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &bouncerv1alpha1.TcpService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTcpService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTcpService{}

// ReconcileTcpService reconciles a TcpService object
type ReconcileTcpService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

const tcpserviceFinalizer = "finalizer.tcpservice.bouncer.cjheppell.github.com"

// Reconcile reads that state of the cluster for a TcpService object and makes changes based on the state read
// and what is in the TcpService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTcpService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling TcpService")

	// Fetch the TcpService instance
	instance := &bouncerv1alpha1.TcpService{}
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

	isTcpServiceMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isTcpServiceMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), tcpserviceFinalizer) {
			reqLogger.Info("finalizing tcp service", "tcpservice_nodePort", instance.Spec.NodePort)
			if err := r.finalizeTcpService(reqLogger, int32(instance.Spec.NodePort)); err != nil {
				return reconcile.Result{}, err
			}
			instance.SetFinalizers(remove(instance.GetFinalizers(), tcpserviceFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	if !instance.Status.Exposed {
		firewallClient, err := cloudprovider.GetFirewallClient()
		if err != nil {
			reqLogger.Error(err, "Error constructing firewall client")
			return reconcile.Result{}, nil
		}

		err = firewallClient.ExposePort(int32(instance.Spec.NodePort))
		if err != nil {
			reqLogger.Error(err, "Error exposing port")
			return reconcile.Result{}, err
		}

		instance.Status.Exposed = true
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "failed to update the status of the tcp service")
			return reconcile.Result{Requeue: false}, nil
		}
	}

	if !contains(instance.GetFinalizers(), tcpserviceFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			reqLogger.Error(err, "Failed to add a finalizer to the tcp service")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileTcpService) addFinalizer(logger logr.Logger, tcpservice *bouncerv1alpha1.TcpService) error {
	logger.Info("adding Finalizer for tcp service", "tcpservice_nodePort", tcpservice.Spec.NodePort)
	tcpservice.SetFinalizers(append(tcpservice.GetFinalizers(), tcpserviceFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), tcpservice)
	if err != nil {
		logger.Error(err, "failed to update tcp service with finalizer", "tcpservice_nodePort", tcpservice.Spec.NodePort)
		return err
	}
	return nil
}

func (r *ReconcileTcpService) finalizeTcpService(logger logr.Logger, port int32) error {
	firewallClient, err := cloudprovider.GetFirewallClient()
	if err != nil {
		return err
	}

	return firewallClient.UnexposePort(port)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
