// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// //go:generate protoc -I ./model/pod --go_out=plugins=grpc:./model/pod ./model/pod/pod.proto

package exampleplugincrd

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	crdutils "github.com/ant31/crd-validation/pkg"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Masterminds/semver"
	"github.com/ligato/cn-infra/config"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/crd-example/pkg/apis/crdexample.io/v1"
	client "github.com/ligato/crd-example/pkg/client/clientset/versioned"
	factory "github.com/ligato/crd-example/pkg/client/informers/externalversions"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Plugin watches K8s resources and causes all changes to be reflected in the ETCD
// data store.
type Plugin struct {
	Deps

	pluginStopCh    chan struct{}
	k8sClientConfig *rest.Config
	k8sClientset    *kubernetes.Clientset
	apiclientset    *apiextcs.Clientset
	crdClient       client.Interface

	// These can be used to stop all the informers, as well as control loops
	// within the application.
	stopChExample      chan struct{}
	stopChExampleEmbed chan struct{}
	// sharedFactory is a shared informer factorys used as a cache for
	// items in the API server. It saves each informer listing and watches the
	// same resources independently of each other, thus providing more up to
	// date results with less 'effort'
	sharedFactory factory.SharedInformerFactory

	// Informer factories per CRD object
	informerExample      cache.SharedIndexInformer
	informerExampleEmbed cache.SharedIndexInformer
}

// Deps defines dependencies of netmesh plugin.
type Deps struct {
	local.PluginInfraDeps
	// Kubeconfig with k8s cluster address and access credentials to use.
	KubeConfig config.PluginConfig
}

const (
	exampleCRDVersion              = "0.0.1"
	exampleCRDVersionAnnotationKey = "crdexample.io/example-crd-version"
	exampleObjectUpdateRetries     = 5
)

var (
	cfg crdutils.Config
)

// Init builds K8s client-set based on the supplied kubeconfig and initializes
// all reflectors.
func (plugin *Plugin) Init() error {
	var err error
	plugin.Log.SetLevel(logging.DebugLevel)
	plugin.pluginStopCh = make(chan struct{})

	kubeconfig := plugin.KubeConfig.GetConfigName()
	plugin.Log.WithField("kubeconfig", kubeconfig).Info("Loading kubernetes client config")
	plugin.k8sClientConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("Failed to build kubernetes client config: %s", err)
	}

	plugin.k8sClientset, err = kubernetes.NewForConfig(plugin.k8sClientConfig)
	if err != nil {
		return fmt.Errorf("Failed to build kubernetes client: %s", err)
	}

	plugin.stopChExample = make(chan struct{})
	plugin.stopChExampleEmbed = make(chan struct{})

	return nil
}

// exampleCrdValidation generates OpenAPIV3 validator for CrdExample CRD
func exampleCrdValidation() *apiextv1beta1.CustomResourceValidation {
	maxLength := int64(64)
	validation := &apiextv1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &apiextv1beta1.JSONSchemaProps{
			Properties: map[string]apiextv1beta1.JSONSchemaProps{
				"spec": apiextv1beta1.JSONSchemaProps{
					Required: []string{"name"},
					Properties: map[string]apiextv1beta1.JSONSchemaProps{
						"name": apiextv1beta1.JSONSchemaProps{
							Type:        "string",
							MaxLength:   &maxLength,
							Description: "ExampleCrd Name",
							Pattern:     `^[a-zA-Z0-9]+[\-a-zA-Z0-9]*$`,
						},
						"uuid": apiextv1beta1.JSONSchemaProps{
							Type:        "string",
							MaxLength:   &maxLength,
							Description: "ExampleCrd UUID",
							Pattern:     `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`,
						},
					},
				},
			},
		},
	}
	return validation
}

// exampleCrdEmbed generates OpenAPIV3 validator for CrdExampleEmbed CRD
func exampleCrdEmbedValidation() *apiextv1beta1.CustomResourceValidation {
	maxLength := int64(64)
	validation := &apiextv1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &apiextv1beta1.JSONSchemaProps{
			Properties: map[string]apiextv1beta1.JSONSchemaProps{
				"spec": apiextv1beta1.JSONSchemaProps{
					Required: []string{"name"},
					Properties: map[string]apiextv1beta1.JSONSchemaProps{
						"name": apiextv1beta1.JSONSchemaProps{
							Type:        "string",
							MaxLength:   &maxLength,
							Description: "CrdExampleEmbed Name",
							Pattern:     `^[a-zA-Z0-9]+[\-a-zA-Z0-9]*$`,
						},
						"uuid": apiextv1beta1.JSONSchemaProps{
							Type:        "string",
							MaxLength:   &maxLength,
							Description: "CrdExampleEmbed UUID",
							Pattern:     `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`,
						},
					},
				},
			},
		},
	}
	return validation
}

// Create the CRD resource, ignore error if it already exists
func createCRD(plugin *Plugin, FullName, Group, Version, Plural, Name string) error {
	flagset := flag.NewFlagSet(Name, flag.ExitOnError)
	flagset.Var(&cfg.Labels, "labels", "Labels")

	crd := crdutils.NewCustomResourceDefinition(crdutils.Config{
		SpecDefinitionName:    FullName,
		EnableValidation:      true,
		Labels:                crdutils.Labels{LabelsMap: cfg.Labels.LabelsMap},
		ResourceScope:         string(apiextv1beta1.NamespaceScoped),
		Group:                 Group,
		Kind:                  Name,
		Version:               Version,
		Plural:                Plural,
		GetOpenAPIDefinitions: v1.GetOpenAPIDefinitions,
	})

	crdClient := plugin.apiclientset
	// Add Example CRD version
	crd.ObjectMeta.Annotations = map[string]string{
		exampleCRDVersionAnnotationKey: exampleCRDVersion,
	}
	// Starting with 1.11.X SpecReplicasPath and StatusReplicasPath have become mandatory fields,
	// but since NewCustomResourceDefinition function does not set them, CRD creation fails in +1.11.0 k8s clusters.
	// As a workaround, setting these two fields manually here.
	crd.Spec.Subresources.Scale.SpecReplicasPath = ".spec.replicas"
	crd.Spec.Subresources.Scale.StatusReplicasPath = ".status.replicas"
	if err := createCRDObject(crd, crdClient); err != nil {
		plugin.Log.Errorf("fail to create CRD with error: %#v", err)
		return err
	}

	return nil
}

func createCRDObject(newCRD *apiextv1beta1.CustomResourceDefinition, crdClient *apiextcs.Clientset) error {
	// First check if the CRD already exists
	oldCRD, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(newCRD.ObjectMeta.Name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("error getting CRD %s, type %s with error: %+v", newCRD.ObjectMeta.Name, newCRD.Spec.Names.Kind, err)
	}
	if apierrors.IsNotFound(err) {
		// If the CRD does not exist, try to create it. There is a check for possible race condition
		// when another example daemon manages to create CRDs between CustomResourceDefinitions().Get()
		// and CustomResourceDefinitions().Create(), if Create returns error AlreadyExists, then createCRDObject
		// does not fail.
		_, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(newCRD)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("fail creating CRD %s, type %s with error: %#v", newCRD.ObjectMeta.Name, newCRD.Spec.Names.Kind, err)
		}
		return nil
	}
	// Check if CRD has the version annotation, if not, then it means it is old CRD
	// and we update it uncoditionally.
	version, ok := oldCRD.ObjectMeta.Annotations[exampleCRDVersionAnnotationKey]
	if !ok {
		// Exisiting CRD does not have the version annotation, updating it to new definition
		// uncoditionally
		return updateCRD(newCRD, crdClient)
	}
	// Existing CRD has version info, next check is to see if existing CRD version is "<" or "==" or ">"
	// if "<" than new CRD, it will be updated, if "==", then no action , if ">" than new CRD, then we will fail as
	// it possible, that older version of NSM controller started on a cluster with newer CRD definitions.
	existingVersion, err := semver.NewVersion(version)
	if err != nil {
		// Since we failed to process existing CRD version, then we update CRD attempting to bring it to the right level
		return updateCRD(newCRD, crdClient)
	}
	newVersion, _ := semver.NewVersion(newCRD.ObjectMeta.Annotations[exampleCRDVersionAnnotationKey])
	if existingVersion.LessThan(newVersion) {
		// It is upgrade case, updating CRD to the new CRD version
		return updateCRD(newCRD, crdClient)
	}
	if existingVersion.GreaterThan(newVersion) {
		// Downgrade scenario, we have to fail and let the user to resolve this inconsistency
		return fmt.Errorf("fail creating CRD %s, as desired version %s is lower than already exisiting CRD object version %s",
			newCRD.Name, newVersion.String(), existingVersion.String())
	}
	// Old CRD version "==" to new CRD version, do nothing
	return nil
}

// updateCRD attempts to update existing CRD with new definitions. number of attempts is defined in
// nsmObjectUpdateRetries. In case of a conflict error code is returned, the update is re-attempted
// as per optimistic concurrency approach.
func updateCRD(newCRD *apiextv1beta1.CustomResourceDefinition, crdClient *apiextcs.Clientset) error {
	for i := 0; i < exampleObjectUpdateRetries; i++ {
		oldCRD, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(newCRD.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("update: error getting CRD %s, type %s with error: %+v", newCRD.ObjectMeta.Name, newCRD.Spec.Names.Kind, err)
		}
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		_, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Update(newCRD)
		if err == nil {
			return nil
		} else if !apierrors.IsConflict(err) {
			return fmt.Errorf("update: fail updating CRD %s, type %s with error: %#v", newCRD.ObjectMeta.Name, newCRD.Spec.Names.Kind, err)
		}
	}
	return fmt.Errorf("update: fail to update CRD after %d retries", exampleObjectUpdateRetries)
}

func informerCrdExample(plugin *Plugin) {
	plugin.informerExample = plugin.sharedFactory.Crdexample().V1().CrdExamples().Informer()
	// We add a new event handler, watching for changes to API resources.
	plugin.informerExample.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: exampleCrdEnqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					exampleCrdEnqueue(cur)
				}
			},
			DeleteFunc: exampleCrdEnqueue,
		},
	)
}

func informerCrdExampleEmbed(plugin *Plugin) {
	plugin.informerExampleEmbed = plugin.sharedFactory.Crdexample().V1().CrdExampleEmbeds().Informer()
	// We add a new event handler, watching for changes to API resources.
	plugin.informerExample.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: exampleCrdEmbedEnqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					exampleCrdEmbedEnqueue(cur)
				}
			},
			DeleteFunc: exampleCrdEmbedEnqueue,
		},
	)
}

// AfterInit This will create all of the CRDs for NetworkServiceMesh.
func (plugin *Plugin) AfterInit() error {
	var err error

	// Create clientset and create our CRD, this only needs to run once
	plugin.apiclientset, err = apiextcs.NewForConfig(plugin.k8sClientConfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an instance of our own API client
	plugin.crdClient, err = client.NewForConfig(plugin.k8sClientConfig)
	if err != nil {
		plugin.Log.Errorf("Error creating CRD client: %s", err.Error())
		panic(err.Error())
	}

	err = createCRD(plugin, v1.FullCRDExampleName,
		v1.Group,
		v1.GroupVersion,
		v1.CRDExamplePlural,
		v1.CRDExampleTypeName)

	if err != nil {
		plugin.Log.Error("Error initializing CrdExample CRD")
		return err
	}

	err = createCRD(plugin, v1.FullCRDExampleEmbedName,
		v1.Group,
		v1.GroupVersion,
		v1.CRDExampleEmbedPlural,
		v1.CRDExampleEmbedTypeName)

	if err != nil {
		plugin.Log.Error("Error initializing CrdExampleEmbed CRD")
		return err
	}

	// We use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	plugin.sharedFactory = factory.NewSharedInformerFactory(plugin.crdClient, time.Second*30)

	informerCrdExample(plugin)
	informerCrdExampleEmbed(plugin)

	// Start the informer. This will cause it to begin receiving updates from
	// the configured API server and firing event handlers in response.
	plugin.sharedFactory.Start(plugin.stopChExample)
	plugin.Log.Info("Started CrdExample informer factory.")

	// Wait for the informer cache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(plugin.stopChExample, plugin.informerExample.HasSynced) {
		plugin.Log.Error("Error waiting for informer cache to sync")
	}

	plugin.Log.Info("CrdExample Informer is ready")

	// Read forever from the work queues
	workforever(plugin, queueExample, plugin.informerExample, plugin.stopChExample)
	workforever(plugin, queueExampleEmbed, plugin.informerExampleEmbed, plugin.stopChExampleEmbed)

	return nil
}

// Close stops all reflectors.
func (plugin *Plugin) Close() error {
	close(plugin.pluginStopCh)
	return nil
}
