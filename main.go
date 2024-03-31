package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5" // router
	"github.com/go-chi/chi/v5/middleware"
	"os"

	//"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	getter "github.com/hashicorp/go-getter"
	"github.com/robfig/cron/v3"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: "./",
		Dir: true,
		//the repository with a subdirectory I would like to clone only
		Src:  "github.com/KTruntaev/auto_updating_testing",
		Mode: getter.ClientModeDir,
		//define the type of detectors go getter should use, in this case only github is needed
		Detectors: []getter.Detector{
			&getter.GitHubDetector{},
		},
		//provide the getter needed to download the files
		Getters: map[string]getter.Getter{
			"git": &getter.GitGetter{},
		},
	}

	//download the files
	if err := client.Get(); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting path %s: %v", client.Src, err)
		os.Exit(1)
	}

	// scheduler ! //////////////////////////////////////////////

	c := cron.New()
	_, err := c.AddFunc("* * * * *", func() {
		fmt.Println("Hello world!")

		//download the files
		if err := client.Get(); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting path %s: %v", client.Src, err)
			os.Exit(1)
		}
	})

	if err != nil {
		panic(err)
	}
	c.Start()

	// Wait for the Cron job to run
	//time.Sleep(5 * time.Minute)
	//
	//c.Stop()

	/////////////////////////////////////////////////////////////

	////download the files
	//if err := client.Get(); err != nil {
	//	fmt.Fprintf(os.Stderr, "Error getting path %s: %v", client.Src, err)
	//	os.Exit(1)
	//}
	//now you should check your temp directory for the files to see if they exist

	//MAIN CODE
	//Initialize router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	////Initialize the template engine
	////Files are provided as a slice of strings
	//paths := []string{
	//	"./static/templates/encycl_article.html",
	//}

	//var err error

	//articleTempl, err = template.ParseFiles(paths...)
	//if err != nil {
	//	panic(err)
	//}

	//Route the webpaths
	//router.Route("/articles/", func(r chi.Router) {
	//	r.Get("/*", serveArticle) // GET /articles/{article_name.html}
	//	//r.Get()
	//})

	var fs = http.FileServer(http.Dir("articles"))

	router.Handle("/*", http.StripPrefix("", fs))

	fmt.Println("Server running!")
	http.ListenAndServe(":3333", fs)
}
