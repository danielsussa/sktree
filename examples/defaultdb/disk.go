package defaultdb

import (
	"encoding/json"
	tree "github.com/danielsussa/tmp_tree"
	"github.com/peterbourgon/diskv"
)

type DiskPersistence struct {
	d *diskv.Diskv
}

func (dmp DiskPersistence) Find(key string) (*tree.Node, bool) {
	b, err := dmp.d.Read(key)
	if err != nil {
		return nil, false
	}

	var node *tree.Node
	err = json.Unmarshal(b, &node)
	if err != nil {
		panic(err)
	}
	return node, true
}

func (dmp DiskPersistence) Add(key string, node *tree.Node) error {
	b, _ := json.Marshal(node)
	err := dmp.d.Write(key, b)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewDefaultDiskDB(filename string) DiskPersistence {
	flatTransform := func(s string) []string {
		return []string{}
	}

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1GB cache.
	d := diskv.New(diskv.Options{
		BasePath:     filename,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024 * 1024,
	})
	return DiskPersistence{d: d}
}
