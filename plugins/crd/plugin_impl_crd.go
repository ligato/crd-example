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

	crdClient := plugin.apiclientset.ApiextensionsV1beta1().CustomResourceDefinitions()

	// First check if the CRD already exists
	oldCRD, err := crdClient.Get(crd.Name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		plugin.Log.Errorf("error getting CRD %s, type %s", crd.Name, crd.Spec.Names.Kind)
		return err
	}
	if apierrors.IsNotFound(err) {
		// If the CRD does not exist, try to create it
		if _, err := crdClient.Create(crd); err != nil {
			plugin.Log.Errorf("error creating CRD %s, type %s", crd.Name, crd.Spec.Names.Kind)
			return err
		}
		plugin.Log.Infof("created CRD %s, type %s", crd.Name, crd.Spec.Names.Kind)
	}
	if err == nil {
		// Now we try to update the CRD
		crd.ResourceVersion = oldCRD.ResourceVersion
		if _, err := crdClient.Update(crd); err != nil {
			plugin.Log.Errorf("error updating CRD %s, type %s", crd.Name, crd.Spec.Names.Kind)
			return err
		}
		plugin.Log.Infof("updated CRD %s, type %s", crd.Name, crd.Spec.Names.Kind)
	}

	return nil
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
