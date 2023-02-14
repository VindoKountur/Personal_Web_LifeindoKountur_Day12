package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get file with name uploadImage from html
		fileNames, fileHandler, err := r.FormFile("uploadImage")
		if err != nil {
			fmt.Println("message : " + err.Error())
			json.NewEncoder(w).Encode("Error Retrieving the file")
			return
		}

		defer fileNames.Close() // For memory leak

		// change image name
		tempFile, err := ioutil.TempFile("uploads", "image-*"+fileHandler.Filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}

		defer tempFile.Close() // For memory leak

		// Read file name by byte
		fileBytes, err := ioutil.ReadAll(fileNames)
		if err != nil {
			fmt.Println(err)
		}

		//create imager temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		filename := data[8:] // uploads/image-akukeren.png

		ctx := context.WithValue(r.Context(), "dataFile", filename)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
