package main

import (
    "log"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
    "path"
	"path/filepath"
	"regexp"
	"time"
)

type Myhandler struct{}

const (
	templateDir = "./view/"
	uploadDir   = "./upload/"
)

var (
    mux map[string]func(http.ResponseWriter, *http.Request)
    workDir, _= filepath.Abs(filepath.Dir(os.Args[0]))
)

func main() {
	server := http.Server{
		Addr:        ":9090",
		Handler:     &Myhandler{},
		ReadTimeout: 10 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/upload"] = upload
	server.ListenAndServe()
}



func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    addr := r.Header.Get("X-Real-IP")
    if addr == "" {
        addr = r.Header.Get("X-Forwarded-For")
        if addr == "" {
            addr = r.RemoteAddr
        }
    }
    log.Printf("Started %s %s for %s %v", r.Method, r.URL.Path, addr, r.Host)

	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	if ok, _ := regexp.MatchString("^/file", r.URL.String()); ok {
        dir := path.Dir(r.URL.String())
        realDir, _ := filepath.Rel("/file", dir)
		http.StripPrefix(dir, http.FileServer(http.Dir("/" + realDir))).ServeHTTP(w, r)
    } else if ok, _ := regexp.MatchString("^/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
	} else {
        index(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
    redirIndex(w, r, nil)
}

func redirIndex(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
    data := map[string]interface{}{"Title":"首页", "WorkDir":workDir, "UploadSuccess":false, "SuccessInfo":""}
    for k, v := range m {
        data[k] = v
    }

	t, _ := template.ParseFiles(templateDir + "index.html")
    t.Execute(w, data)
}


func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(templateDir + "file.html")
        t.Execute(w, map[string]interface{}{"Title":"上传文件", "Host":r.Host})
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
            fmt.Fprintf(w, "%v:%v", "上传错误", err)
			return
		}

        os.Mkdir(uploadDir, os.ModeDir)

		filename := handler.Filename
        log.Printf("upload %s", filename)
		f, _ := os.Create(uploadDir + handler.Filename)
		if err != nil {
            fmt.Fprintf(w, "%v:%v", "上传失败", err)
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
            fmt.Fprintf(w, "%v:%v", "写入失败", err)
			return
		}
		filedir, _ := filepath.Abs(uploadDir + filename)

        redirIndex(w, r, map[string]interface{}{"UploadSuccess":true, "SuccessInfo":filename + "上传完成,服务器地址:"+filedir})
	}
}

