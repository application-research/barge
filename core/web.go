package core

import (
	"encoding/json"
	"fmt"
	"github.com/application-research/estuary/util"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var UiWebCmd = &cli.Command{
	Name:        "web",
	Description: "barge web is a command to start the web UI",
	Usage:       "barge web",
	Action: func(context *cli.Context) error {

		// create the dir first.
		os.Mkdir("upload", 0775)

		//	website host.
		fs := http.FileServer(http.Dir("./web"))

		//	rest endpoints.
		http.Handle("/", fs)
		http.HandleFunc("/api/v0/plumb/file", func(w http.ResponseWriter, r *http.Request) {
			enableCors(&w)
			var contentResponse *util.ContentAddResponse
			var jsonResponse []byte
			var err error

			if r.Method == "POST" { // post only

				//	get the file
				file, handler, err := r.FormFile("file")
				if err != nil {
					return
				}

				defer file.Close()
				defer func() {
					// remove the temp file
					os.Remove("upload/" + handler.Filename)
				}()

				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)

				// Create a temporary file within our temp-images directory that follows
				// a particular naming pattern
				fileBytes, err := ioutil.ReadAll(file)
				ioutil.WriteFile("./upload/"+handler.Filename, fileBytes, 0644)

				//	grab the local file and pass it here.
				fmt.Println(r.FormValue("fpath"))
				contentResponse, err = PlumbAddFile(context, "./upload/"+handler.Filename, handler.Filename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					jsonResponse, _ = json.Marshal(map[string]string{
						"status": fmt.Sprint(http.StatusBadRequest),
						"error":  err.Error(),
					})
					_, err = io.WriteString(w, string(jsonResponse))
				}
			}

			contentResponseJson, err := json.Marshal(contentResponse)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = io.WriteString(w, string(contentResponseJson))
			if err != nil {
				return
			}
		})
		http.HandleFunc("/api/v0/plumb/files", func(w http.ResponseWriter, r *http.Request) {
			//	upload a list of files.

		})
		http.HandleFunc("/api/v0/plumb/car", func(w http.ResponseWriter, r *http.Request) {
			var contentResponse *util.ContentAddResponse
			var jsonResponse []byte
			var err error

			if r.Method == "POST" { // post only

				//	get the file
				file, handler, err := r.FormFile("file")
				if err != nil {
					return
				}

				defer file.Close()
				defer func() {
					// remove the temp file
					os.Remove("upload/" + handler.Filename)
				}()

				fmt.Printf("Uploaded File: %+v\n", handler.Filename)
				fmt.Printf("File Size: %+v\n", handler.Size)
				fmt.Printf("MIME Header: %+v\n", handler.Header)

				// Create a temporary file within our temp-images directory that follows
				// a particular naming pattern
				fileBytes, err := ioutil.ReadAll(file)
				ioutil.WriteFile("./upload/"+handler.Filename, fileBytes, 0644)

				//	grab the local file and pass it here.
				fmt.Println(r.FormValue("fpath"))
				contentResponse, err = PlumbAddCar(context, "./upload/"+handler.Filename, handler.Filename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					jsonResponse, _ = json.Marshal(map[string]string{
						"status": fmt.Sprint(http.StatusBadRequest),
						"error":  err.Error(),
					})
					_, err = io.WriteString(w, string(jsonResponse))
				}
			}

			contentResponseJson, err := json.Marshal(contentResponse)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = io.WriteString(w, string(contentResponseJson))
			if err != nil {
				return
			}
		})
		http.HandleFunc("/api/v0/get-files", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {

			}
			fmt.Println("get files")

		})
		log.Print("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
		return err
	},
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
