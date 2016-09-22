package main

import (
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SetupDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB() {
	DB.Close()
}

func GetKey(bucket string, key string) (string, error) {
	var destVal string
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			destVal = string(b.Get([]byte(key)))
			return nil
		}
		return errors.New("no such bucket")
	})
	if err != nil {
		return "", err
	}
	return destVal, nil
}

func SetKey(bucket string, key string, value string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteKey(bucket string, key string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		if b != nil {
			err = b.Delete([]byte(key))
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("no such bucket")
	})
}

func DeleteBucket(bucket string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			// Inner keys must be deleted to prevent them to reappear if the bucket is recreated.
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				err := b.Delete(k)
				if err != nil {
					return err
				}
			}
			return tx.DeleteBucket([]byte(bucket))
		}
		return errors.New("no such bucket")
	})
}

func ListBuckets() ([]string, error) {
    var buckets []string
	err := DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
            buckets = append(buckets, string(name))
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
    return buckets, nil
}

func ListBucketKeys(bucket string) ([]KeyValue, error) {
    var kv []KeyValue
    err := DB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				kv = append(kv, KeyValue{string(k), string(v)})
			}
			return nil
		}
		return errors.New("no such bucket")
    })
    if err != nil {
        return nil, err
    }
    return kv, nil
}