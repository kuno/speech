package main

import (
        "os"
        "fmt"
        "path"
        "io/ioutil"
        "net/http"
        "html/template"
)

const (
        Target = "slices"
        Template = "./index.html.temp"
        listenAddr = "localhost:5000"
)

type bytes []byte

type Value struct {
        Articles []template.HTML
}

func findDir(target string, files []os.FileInfo) bool {
        found := false

        for _, f := range(files) {
                if target == f.Name() {
                        found = true
                }
        }

        return found
}

func readFiles(dir string) Value {
        p := ""
        contents := make([]template.HTML, 100)

        files, _ := ioutil.ReadDir(dir)

        fmt.Println(files)

        // if err != nil {
        //         return
        // }

        for _, f := range(files) {
                if f.Name() != "." {
                        p = path.Join(dir, f.Name())

                        c, _ := ioutil.ReadFile(p)

                        // if err != nil {
                        //         return
                        // }

                        contents = append(contents, template.HTML(string(c)))
                }
        }

        fmt.Println(contents)

        v := Value{
                Articles: contents,
        }

        return v
}

func handler(w http.ResponseWriter, r *http.Request) {
        html, err := ioutil.ReadFile(Template)

        if err != nil {
                fmt.Println("error")
                html = bytes("err")
        }

        t := template.New("template")
        temp, err := t.Parse(string(html))

        if err != nil {
        }

        a := make([]template.HTML, 1)
        a[0] = template.HTML("<h2>Hello</h2>")

        v := Value{
                Articles: a,
        }
        err = temp.Execute(w, v)

        if err != nil {
        }

        return
}

func render() {
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
        readFiles(slicesDir)

        http.Handle("/", http.FileServer(http.Dir("./")))
        http.HandleFunc("/go", handler)
        http.ListenAndServe(listenAddr, nil)
}

