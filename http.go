package main

import (
	"log"
	"net/http"

	"github.com/bahlo/goat"
)

type bucketKeys map[string][]KeyValue

func listAllKeys(res http.ResponseWriter, req *http.Request, params goat.Params) {
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
	router.Get("/", "listAllKeys", listAllKeys)
	router.Get("/:bucket", "listBucketKeys", listBucketKeys)
	router.Delete("/:bucket", "deleteBucket", deleteBucket)
	router.Get("/:bucket/:key", "getBucketKey", getBucketKey)
	router.Put("/:bucket/:key", "setBucketKey", setBucketKey)
	router.Post("/:bucket/:key", "setBucketKey", setBucketKey)
	router.Delete("/:bucket/:key", "deleteBucketKey", deleteBucketKey)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Unable to start web server: %s", err.Error())
	}
}
