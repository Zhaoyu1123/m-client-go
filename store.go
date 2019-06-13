package robot

import (
	"k8s.io/client-go/tools/cache"
)

type store interface {
	List() []interface{}

	ListKeys() []string

	GetByKey(key string) (item interface{}, exists bool, err error)
}

type IndexerSet []cache.Indexer

func (set IndexerSet) List() (l []interface{}) {
	for _, indexer := range set {
		l = append(l, indexer.List()...)
	}
	return
}

func (set IndexerSet) ListKeys() (keys []string) {
	for _, indexer := range set {
		keys = append(keys, indexer.ListKeys()...)
	}
	return
}

func (set IndexerSet) GetByKey(key string) (item interface{}, exists bool, err error) {
	for _, indexer := range set {
		item, exists, err = indexer.GetByKey(key)
		if err != nil {
			return
		}
		if exists {
			return
		}
	}
	return
}
