package robot

import (
	"k8s.io/client-go/tools/cache"
)

type store interface {
	List(Resource) []interface{}

	ListKeys(Resource) []string

	GetByKey(r Resource, key string) (items []interface{}, exists bool)
}

var _ store = mapIndexerSet{}

type mapIndexerSet map[Resource][]cache.Indexer

func (mt mapIndexerSet) List(r Resource) (l []interface{}) {
	switch r {
	case All:
		for _, set := range mt {
			for _, indexer := range set {
				l = append(l, indexer.List()...)
			}
		}
	case Services:
		for _, indexer := range mt[Services] {
			l = append(l, indexer.List()...)
		}
	case Pods:
		for _, indexer := range mt[Pods] {
			l = append(l, indexer.List()...)
		}
	case Endpoints:
		for _, indexer := range mt[Endpoints] {
			l = append(l, indexer.List()...)
		}
	case ConfigMaps:
		for _, indexer := range mt[ConfigMaps] {
			l = append(l, indexer.List()...)
		}
	}
	return
}

func (mt mapIndexerSet) ListKeys(r Resource) (keys []string) {
	switch r {
	case All:
		for _, set := range mt {
			for _, indexer := range set {
				keys = append(keys, indexer.ListKeys()...)
			}
		}
	case Services:
		for _, indexer := range mt[Services] {
			keys = append(keys, indexer.ListKeys()...)

		}
	case Pods:
		for _, indexer := range mt[Pods] {
			keys = append(keys, indexer.ListKeys()...)
		}
	case Endpoints:
		for _, indexer := range mt[Endpoints] {
			keys = append(keys, indexer.ListKeys()...)
		}
	case ConfigMaps:
		for _, indexer := range mt[ConfigMaps] {
			keys = append(keys, indexer.ListKeys()...)
		}
	}
	return
}

func (mt mapIndexerSet) GetByKey(r Resource, key string) ([]interface{}, bool) {
	var iterms []interface{}
	ok := false

	switch r {
	case All:
		for _, set := range mt {
			for _, indexer := range set {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					iterms = append(iterms, item)
				}
			}
		}
	case Services:
		for _, indexer := range mt[Services] {
			item, exists, err := indexer.GetByKey(key)
			if err != nil {
				continue
			}
			if exists {
				ok = true
				iterms = append(iterms, item)
			}
		}
	case Pods:
		for _, indexer := range mt[Pods] {
			item, exists, err := indexer.GetByKey(key)
			if err != nil {
				continue
			}
			if exists {
				ok = true
				iterms = append(iterms, item)
			}
		}
	case Endpoints:
		for _, indexer := range mt[Endpoints] {
			item, exists, err := indexer.GetByKey(key)
			if err != nil {
				continue
			}
			if exists {
				ok = true
				iterms = append(iterms, item)
			}
		}
	case ConfigMaps:
		for _, indexer := range mt[ConfigMaps] {
			item, exists, err := indexer.GetByKey(key)
			if err != nil {
				continue
			}
			if exists {
				ok = true
				iterms = append(iterms, item)
			}
		}
	}

	return iterms, ok
}
