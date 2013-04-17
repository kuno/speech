package main

import (
        "os"
        "fmt"
        "path"
        "io/ioutil"
        "net/http"
        "html/template"
        "github.com/russross/blackfriday"
)

const (
        Target = "slices"
        Template = "./index.html.temp"
        listenAddr = "0.0.0.0:5000"
)

type bytes []byte

type HTML template.HTML

type HTMLS []template.HTML

type Slice struct {
        Title string
        Articles HTMLS
}

type Handler func (w http.ResponseWriter, r *http.Request)

func findDir(target string, files []os.FileInfo) bool {
        found := false

        for _, f := range(files) {
                if target == f.Name() {
                        found = true
                }
        }

        return found
}

func genSlice(t, dir string) Slice {
        p := ""
        contents := make(HTMLS, 0)

        files, _ := ioutil.ReadDir(dir)

        fmt.Println(files)

        // if err != nil {
        //         return
        // }

        for _, f := range(files) {
                if f.Name() != "." {
                        p = path.Join(dir, f.Name())

                        m, _ := ioutil.ReadFile(p)

                        // if err != nil {
                        //         return
                        // }
                        h := blackfriday.MarkdownCommon(m)

                        contents = append(contents, template.HTML(string(h)))
                }
        }

        fmt.Println("contents are")
        fmt.Println(contents)

        s := Slice{
                Title: t,
                Articles: contents,
        }

        return s
}

func genHandler(s Slice) (h Handler) {
        return func(w http.ResponseWriter, r *http.Request) {
                html, err := ioutil.ReadFile(Template)

                if err != nil {
                        fmt.Println("error")
                }

                t := template.New("template")
                temp, err := t.Parse(string(html))

                if err != nil {
                }

                err = temp.Execute(w, s)

                if err != nil {
                }
        }

}

func main() {
        var dir string

        if len(os.Args) > 1 {
                dir = os.Args[1]
        } else {
                dir = "."
        }

        infos, err := ioutil.ReadDir(dir)

        if err != nil {
                return
        }

        if findDir(Target, infos) {
                fmt.Println("Found slices directory")
        } else {
                fmt.Println("Not found directory")
        }

        slicesDir := path.Join(dir, Target)

        fmt.Println(slicesDir)
        slice := genSlice(dir, slicesDir)
        handler := genHandler(slice)

        http.Handle("/", http.FileServer(http.Dir("./")))
        http.HandleFunc("/go", handler)
        http.ListenAndServe(listenAddr, nil)
}

