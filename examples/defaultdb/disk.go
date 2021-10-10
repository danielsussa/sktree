package defaultdb

import (
	"github.com/peterbourgon/diskv"
)

type DiskPersistence struct {
	d *diskv.Diskv
}

func (dmp DiskPersistence) Find(key string) (string, bool) {
	b, err := dmp.d.Read(key)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func (dmp DiskPersistence) Add(key, value string) error {
	err := dmp.d.Write(key, []byte(value))
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
