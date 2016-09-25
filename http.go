package main

import (
	"log"
	"net/http"

	"github.com/nmaggioni/goat"
)

type bucketKeys map[string][]KeyValue

func setPoweredByHeader(res http.ResponseWriter) {
	res.Header().Set("X-Powered-By", "gerph")
}

func listAllBucketsKeys(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	buckets, err := ListBuckets()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	bk := make(bucketKeys)
	for _, bucketName := range buckets {
		keys, _ := ListBucketKeys(bucketName)
		bk[bucketName] = keys
	}

	goat.WriteJSON(res, bk)
}

func listBucketKeys(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	keys, err := ListBucketKeys(params["bucket"])
	if err != nil {
		if err.Error() == "no such bucket" {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	goat.WriteJSON(res, keys)
}

func deleteBucket(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	if params["bucket"] != "" {
		err := DeleteBucket(params["bucket"])
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func getBucketKey(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	if params["bucket"] != "" && params["key"] != "" {
		value, err := GetKey(params["bucket"], params["key"])
		if err != nil {
			if err.Error() == "no such bucket" {
				res.WriteHeader(http.StatusNoContent)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if value != "" {
			res.WriteHeader(http.StatusOK)
			goat.WriteJSON(res, map[string]string{
				"key":   params["key"],
				"value": value,
			})
			return
		}
		res.WriteHeader(http.StatusNoContent)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func setBucketKey(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(req.Form["value"]) > 0 {
		value := req.Form["value"][0]
		err = SetKey(params["bucket"], params["key"], value)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		goat.WriteJSON(res, map[string]string{
			"key":   params["key"],
			"value": value,
		})
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func deleteBucketKey(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	if params["bucket"] != "" && params["key"] != "" {
		err := DeleteKey(params["bucket"], params["key"])
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

// Listen configures and starts a web server, enclosing it in an asynchronous goroutine.
func Listen(port string) {
	router := goat.New()
	router.Options("/", "help", func(res http.ResponseWriter, req *http.Request, _ goat.Params) {
		setPoweredByHeader(res)
		goat.WriteJSON(res, router.Index())
	})

	router.Get("/", "list_all_buckets_and_keys", listAllBucketsKeys)
	router.Get("/:bucket", "list_bucket_keys", listBucketKeys)
	router.Delete("/:bucket", "delete_bucket", deleteBucket)
	router.Get("/:bucket/:key", "get_bucket_key", getBucketKey)
	router.Put("/:bucket/:key", "set_bucket_key", setBucketKey)
	router.Post("/:bucket/:key", "set_bucket_key", setBucketKey)
	router.Delete("/:bucket/:key", "delete_bucket_key", deleteBucketKey)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Unable to start web server: %s", err.Error())
	}
}
