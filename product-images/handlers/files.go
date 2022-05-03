package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Vergangenheit/go-grpc-ms/product-images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{log: l, store: s}
}

// Upload REST implements the http.Handler interface
func (f *Files) UploadRest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// check if the file is a valid name and file
	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}
	// check that the filepath is a valid name and file
	f.saveFile(id, fn, rw, r.Body)

}

//
func (f *Files) UploadMultiPart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		f.log.Error("Unable to parse form", "err", err)
		return
	}
	id, iderr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Form value", "id", id)
	if iderr != nil {
		http.Error(rw, "Expected id as integer", http.StatusBadRequest)
		f.log.Error("Bad Request", "err", err)
		return
	}

	fl, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		f.log.Error("Bad request", "err", err)
		return
	}
	// save the file
	f.saveFile(r.FormValue("id"), fh.Filename, rw, fl)

}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
