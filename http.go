package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"syscall"

	"github.com/nmaggioni/goat"
)

type (
	bucketKeys map[string][]KeyValue
	stats      struct {
		BucketsNumber int        `json:"bucketsNumber"`
		KeysNumber    int        `json:"keysNumber"`
		DiskUsed      int64      `json:"diskUsedBytes"`
		DiskFree      uint64     `json:"diskFreeBytes"`
		Keys          bucketKeys `json:"keys"`
	}
)

func getAllBucketsAndKeys() (bucketKeys, error) {
	buckets, err := ListBuckets()
	if err != nil {
		return nil, err
	}
	bk := make(bucketKeys)
	for _, bucketName := range buckets {
		keys, _ := ListBucketKeys(bucketName)
		bk[bucketName] = keys
	}
	return bk, nil
}

func setPoweredByHeader(res http.ResponseWriter) {
	res.Header().Set("X-Powered-By", "gerph")
}

func listAllBucketsKeys(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	bk, err := getAllBucketsAndKeys()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
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

func serveWebInterface(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	if params["asset"] != "" {
		http.ServeFile(res, req, "./public/assets/"+params["asset"])
	} else {
		http.ServeFile(res, req, "./public"+req.URL.Path[1:])
	}
}

func serveWebStats(res http.ResponseWriter, req *http.Request, params goat.Params) {
	bucketsNumber, _ := CountBuckets()
	keysNumber, _ := CountKeys()
	dbFileStats, _ := os.Stat(DBPath)
	dbSizeBytes := dbFileStats.Size()
	var diskAvailableBytes uint64
	if runtime.GOOS == "windows" {
		diskAvailableBytes = 0  // Undefined syscall?  http://stackoverflow.com/a/20110856
	} else {
		var stat syscall.Statfs_t
		cwd, _ := os.Getwd()
		syscall.Statfs(cwd, &stat)
		diskAvailableBytes = stat.Bavail * uint64(stat.Bsize)
	}
	keys, _ := getAllBucketsAndKeys()

	webStats := stats{bucketsNumber, keysNumber, dbSizeBytes, diskAvailableBytes, keys}

	setPoweredByHeader(res)
	goat.WriteJSON(res, webStats)
}

func serveBackup(res http.ResponseWriter, req *http.Request, params goat.Params) {
	setPoweredByHeader(res)
	res.Header().Set("Content-Disposition", "attachment; filename=\""+path.Base(DBPath)+"\"")
	res.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(res, req, DBPath)
}

// Listen configures and starts a web server and its router.
func Listen(port string) {
	router := goat.New()
	router.Get("/", "web_ui", serveWebInterface)
	router.Get("/asset/:asset", "serve_asset", serveWebInterface)
	router.Get("/stats", "web_ui_stats", serveWebStats)
	router.Get("/backup", "download_backup", serveBackup)

	api := router.Subrouter("/api")
	api.Options("/", "help", func(res http.ResponseWriter, req *http.Request, _ goat.Params) {
		setPoweredByHeader(res)
		goat.WriteJSON(res, api.Index())
	})
	api.Get("/", "list_all_buckets_and_keys", listAllBucketsKeys)
	api.Get("/:bucket", "list_bucket_keys", listBucketKeys)
	api.Delete("/:bucket", "delete_bucket", deleteBucket)
	api.Get("/:bucket/:key", "get_bucket_key", getBucketKey)
	api.Put("/:bucket/:key", "set_bucket_key", setBucketKey)
	api.Post("/:bucket/:key", "set_bucket_key", setBucketKey)
	api.Delete("/:bucket/:key", "delete_bucket_key", deleteBucketKey)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Unable to start web server: %s", err.Error())
	}
}
