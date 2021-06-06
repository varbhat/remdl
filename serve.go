package main

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func DirServe(w http.ResponseWriter, r *http.Request) {
	var lcheck bool
	if Cred.Locked {
		err := authHelper(w, r)
		if err != nil {
			Err.Println(err)
			return
		}
		lcheck = true
	}
	if val := r.URL.Query().Get("dl"); val != "" {
		file := filepath.Join(Cred.DirPath, strings.TrimPrefix(r.URL.Path, "/dl/"))
		info, err := os.Stat(file)
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		} else if info.IsDir() {
			if val == "tar" {
				TarDir(file, w, path.Base(r.URL.Path))
				return
			} else if val == "zip" {
				ZipDir(file, w, path.Base(r.URL.Path))
				return
			}
		}
	}
	if val := r.URL.Query().Get("act"); val == "rm" {
		if !lcheck {
			err := authHelper(w, r)
			if err != nil {
				Err.Println(err)
				return
			}
		}
		file := filepath.Join(Cred.DirPath, strings.TrimPrefix(r.URL.Path, "/dl/"))
		_, err := os.Stat(file)
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		err = os.RemoveAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.Write([]byte("Successfully Removed"))
			return
		}
	}

	if strings.HasSuffix(r.URL.Path, "/") {
		fentries, err := os.ReadDir(filepath.Join(Cred.DirPath, strings.TrimPrefix(r.URL.Path, "/dl/")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
		fmt.Fprintf(w, `<a href="/">Home</a>
<a href="/dl">Downloads</a>
<a href="/list">List</a>`)
		fmt.Fprintf(w, "<pre>\n")
		for _, fentry := range fentries {
			name := fentry.Name()
			in, er := fentry.Info()
			if er != nil {
				continue
			}

			url := url.URL{Path: name}
			if in.IsDir() {
				fmt.Fprintf(w, "<fieldset><p><a href=\"%s\">%s</a>  <a href=\"%s?act=rm\"><button>Delete</button></a> <a href=\"%s?dl=tar\"><button>TAR</button></a> <a href=\"%s?dl=zip\"><button>ZIP</button></a></p></fieldset>\n", url.String(), template.HTMLEscapeString(name), url.String(), url.String(), url.String())
			} else {
				fmt.Fprintf(w, "<fieldset><p><a href=\"%s\">%s</a>  <a href=\"%s?act=rm\"><button>Delete</button></a></p></fieldset>\n", url.String(), template.HTMLEscapeString(name+"/"), url.String())
			}
		}
		fmt.Fprintf(w, "</pre>\n")
		return
	}

	http.StripPrefix("/dl/", http.FileServer(http.Dir(Cred.DirPath))).ServeHTTP(w, r)
}

func TarDir(dirpath string, w http.ResponseWriter, name string) {
	w.Header().Set("Content-Type", "application/x-tar")
	w.Header().Set("Content-disposition", `attachment; filename="`+name+`.tar"`)
	w.WriteHeader(http.StatusOK)
	tw := tar.NewWriter(w)
	defer tw.Close()

	_ = filepath.WalkDir(dirpath, func(p string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, ierr := de.Info()
		if ierr != nil {
			return ierr
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(dirpath, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		h.Name = rel
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		n, err := io.Copy(tw, f)
		if info.Size() != n {
			return errors.New("size mismatch: " + rel)
		}
		return err
	})
}

func ZipDir(dirpath string, w http.ResponseWriter, name string) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-disposition", `attachment; filename="`+name+`.zip"`)
	w.WriteHeader(http.StatusOK)
	zw := zip.NewWriter(w)
	defer zw.Close()

	_ = filepath.WalkDir(dirpath, func(p string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, ierr := de.Info()
		if ierr != nil {
			return ierr
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(dirpath, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		h.Name = rel
		//h.Method = zip.Deflate

		zf, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}

		n, err := io.Copy(zf, f)
		if info.Size() != n {
			return errors.New("size mismatch: " + rel)
		}

		return err

	})

}
