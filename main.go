package main

import (
	"encoding/json"
	"fmt"
	"github.com/application-research/barge/core"
	"github.com/application-research/estuary/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
	app.Name = "barge"
	app.Commands = []*cli.Command{
		core.LoginCmd,
		core.InitCmd,
		core.ConfigCmd,
		core.PlumbCmd,
		core.CollectionsCmd,
		core.BargeAddCmd,
		core.BargeStatusCmd,
		core.BargeSyncCmd,
		core.BargeCheckCmd,
		core.BargeShareCmd,
		UiWebCmd,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logging",
		},
	}
	app.Before = func(cctx *cli.Context) error {
		if err := loadConfig(); err != nil {
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var UiWebCmd = &cli.Command{
	Name: "web-ui",
	Action: func(context *cli.Context) error {

		// create the dir first.
		os.Mkdir("upload", 0775)

		//	host
		fs := http.FileServer(http.Dir("./web"))
		http.Handle("/", fs)
		http.HandleFunc("/api/v0/plumb/file", func(w http.ResponseWriter, r *http.Request) {
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
				contentResponse, err = core.PlumbAddFile(context, "./upload/"+handler.Filename, handler.Filename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					jsonResponse, _ = json.Marshal(map[string]string{
						"status": string(http.StatusBadRequest),
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
				contentResponse, err = core.PlumbAddCar(context, "./upload/"+handler.Filename, handler.Filename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)

					jsonResponse, _ = json.Marshal(map[string]string{
						"status": string(http.StatusBadRequest),
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

func loadConfig() error {
	bargeDir, err := homedir.Expand("~/.barge")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(bargeDir, 0775); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.barge")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.WriteConfigAs(filepath.Join(bargeDir, "config"))
		} else {
			fmt.Printf("read err: %#v\n", err)
			return err
		}
	}
	return nil
}
