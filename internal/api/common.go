package api

import "net/http"

func setHeader(header http.Header) {
	header.Set("Access-Control-Allow-Headers", "Content-Type")
	header.Set("Access-Control-Allow-Methods", "GET")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "application/json")
}
