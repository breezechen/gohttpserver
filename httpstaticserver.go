package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type HTTPStaticServer struct {
	Root  string
	Theme string

	m *mux.Router
}

func NewHTTPStaticServer(root string, theme string) *HTTPStaticServer {
	if root == "" {
		root = "."
	}
	m := mux.NewRouter()
	s := &HTTPStaticServer{
		Root:  root,
		Theme: theme, // TODO: need to parse from command line
		m:     m,
	}
	m.HandleFunc("/-/raw/{path:.*}", s.hFileOrDirectory)
	m.HandleFunc("/-/zip/{path:.*}", s.hZip)
	m.HandleFunc("/-/json/{path:.*}", s.hJSONList)
	m.HandleFunc("/{path:.*}", s.hIndex).Methods("GET")
	return s
}

func (s *HTTPStaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func (s *HTTPStaticServer) hIndex(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	relPath := filepath.Join(s.Root, path)
	finfo, err := os.Stat(relPath)
	if err == nil && finfo.IsDir() {
		tmpl.Execute(w, s)
	} else {
		http.ServeFile(w, r, relPath)
	}
}

func (s *HTTPStaticServer) hZip(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	CompressToZip(w, path)
}

func (s *HTTPStaticServer) hFileOrDirectory(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Path:", s.Root, path)
	http.ServeFile(w, r, filepath.Join(s.Root, path))
}

type ListResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size string `json:"size"`
}

func (s *HTTPStaticServer) hJSONList(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	lrs := make([]ListResponse, 0)
	fd, err := os.Open(filepath.Join(s.Root, path))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer fd.Close()

	files, err := fd.Readdir(-1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, file := range files {
		lr := ListResponse{
			Name: file.Name(),
			Path: filepath.Join(path, file.Name()), // lstrip "/"
		}
		if file.IsDir() {
			fileName := deepPath(filepath.Join(s.Root, path), file.Name())
			lr.Name = fileName
			lr.Path = filepath.Join(path, fileName)
			lr.Type = "dir"
			lr.Size = "-"
		} else {
			lr.Type = "file"
			lr.Size = formatSize(file)
		}
		lrs = append(lrs, lr)
	}

	data, _ := json.Marshal(lrs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func deepPath(basedir, name string) string {
	isDir := true
	// loop max 5, incase of for loop not finished
	maxDepth := 5
	for depth := 0; depth <= maxDepth && isDir; depth += 1 {
		finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			break
		}
		if finfos[0].IsDir() {
			name = filepath.Join(name, finfos[0].Name())
		} else {
			break
		}
	}
	return name
}
