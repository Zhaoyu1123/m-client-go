package robot

import (
	"k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Robot is an interface for monitor k8s multi-cluster resources.
type Robot interface {
	// Discover define which resources will be discovered under the fixed namespace of k8s
	// If the namespace is empty, it will discover all k8s namespaces
	Discover(resources []Resource, namespace string)

	// Run start up the robot.
	// Start monitoring resources and sending events to the queue.
	Run()

	// Stop stop monitoring resources
	// Empty queue and recycle
	Stop()

	queue

	store
}

type controller struct {
	clients 	[]*kubernetes.Clientset
	informers 	informerSet

	stop        chan struct{}

	queue

	store
}

var _ Robot = &controller{}

func NewRobot(masterUrl, kubeconfigPath []string) (Robot, error) {
	cs, err := newClientSet(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return &controller{
		clients:cs,
		queue: newWorkQueue(),
		stop: make(chan struct{}, 1),
	}, nil
}

func (c *controller) Discover(resources []Resource, namespace string) {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.push(QueueObject{EventAdd, key})
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				c.push(QueueObject{EventUpdate, key})
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.push(QueueObject{EventDelete, key})
			}
		},
	}

	mis := make(mapIndexerSet)
	fs := make(informerSet, 0)
	for _, r := range resources {
		for _, client := range c.clients {
			lw := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), r.String(), namespace, fields.Everything())
			var indexer cache.Indexer
			var informer cache.Controller
			switch r {
			case Services:
				indexer, informer = cache.NewIndexerInformer(lw, &v1.Service{}, 0, handler, cache.Indexers{})
				mis[Services] = append(mis[Services], indexer)
			case Pods:
				indexer, informer = cache.NewIndexerInformer(lw, &v1.Pod{}, 0, handler, cache.Indexers{})
				mis[Pods] = append(mis[Pods], indexer)
			case Endpoints:
				indexer, informer = cache.NewIndexerInformer(lw, &v1.Endpoints{}, 0, handler, cache.Indexers{})
				mis[Endpoints] = append(mis[Endpoints], indexer)
			case ConfigMaps:
				indexer, informer = cache.NewIndexerInformer(lw, &v1.ConfigMap{}, 0, handler, cache.Indexers{})
				mis[ConfigMaps] = append(mis[ConfigMaps], indexer)
			}
			fs = append(fs, informer)
		}
	}

	c.informers = fs
	c.store = mis
}

func (c *controller) Run()  {
	defer c.queue.close()

	c.informers.run(c.stop)

	<-c.stop
}

func (c *controller) Stop()  {
	c.stop <- struct{}{}
}

type informerSet []cache.Controller

func (s informerSet) run(done chan struct{}) {
	for _, one := range s {

		go one.Run(done)

		if !cache.WaitForCacheSync(done, one.HasSynced) {
			panic("Timed out waiting for caches to sync")
		}
	}
}

func newClientSet (masterUrl, kubeconfigPath []string) ([]*kubernetes.Clientset, error) {
	cs := make([]*kubernetes.Clientset, 0)
	if len(masterUrl) != 0 {
		for _, uri := range masterUrl {
			config, err := clientcmd.BuildConfigFromFlags(uri, "")
			if err != nil {
				return nil, err
			}
			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return nil, err
			}
			cs = append(cs, client)
		}
	}

	if len(kubeconfigPath) != 0 {
		for _, path := range kubeconfigPath {
			config, err := clientcmd.BuildConfigFromFlags("", path)
			if err != nil {
				return nil, err
			}
			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return nil, err
			}
			cs = append(cs, client)
		}
	}
	return cs, nil
}
