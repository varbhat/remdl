package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
)

type Resp struct {
	Start    time.Time
	URL      string
	Filename string
	Size     int64
	Err      error
}

func authHelper(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	c, err := r.Cookie("session_token")
	if err != nil {
		if val := r.URL.Query().Get("token"); val != "" {
			if val != Cred.Token {
				return bauthHelper(w, r)
			}
		} else {
			return bauthHelper(w, r)
		}
	} else {
		if c.Value != Cred.Token {
			return bauthHelper(w, r)
		}
	}
	return nil
}

func bauthHelper(w http.ResponseWriter, r *http.Request) (err error) {
	var password string
	var ok bool
	w.Header().Set("WWW-Authenticate", `Basic realm="Wirelit Auth Protected Space"`)
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return fmt.Errorf("not authorized")
	}
	if username != Cred.Username || password != Cred.Password {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return fmt.Errorf("not authorized")
	} else {
		session := uuid.New().String()
		Cred.Token = session
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   session,
			Expires: time.Now().Add(48 * time.Hour),
		})
	}
	return nil
}

func main() {
	http.HandleFunc("/dl/", DirServe)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/cancel", cancelreq)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/toggle", toggleLock)
	http.HandleFunc("/creds", credHandler)
	http.HandleFunc("/changecreds", changecreds)
	http.HandleFunc("/", mainHandler)

	Info.Println("Starting server on", Cred.ListenAddress)
	if Cred.TLSCertPath != "" && Cred.TLSKeyPath != "" {
		Info.Println("Serving the HTTPS with TLS Cert ", Cred.TLSCertPath, " and TLS Key", Cred.TLSKeyPath)
		Err.Fatal(http.ListenAndServeTLS(Cred.ListenAddress, Cred.TLSCertPath, Cred.TLSKeyPath, nil))
	} else {
		Err.Fatal(http.ListenAndServe(Cred.ListenAddress, nil))
	}
}

func changecreds(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	if r.Method != http.MethodPost {
		return
	}
	usn := r.FormValue("username")
	pw := r.FormValue("password")
	if len(usn) > 0 {
		Cred.Username = usn
		Info.Println("Username changed")
	}
	if len(pw) > 0 {
		Cred.Password = pw
		Info.Println("Password changed")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func credHandler(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	fmt.Fprint(w, `<!DOCTYPE html>
		<html lang="en">
		  <head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Remdl</title>
		  </head>
		  <body>
			<main>
			  <a href="/">Home</a>
			  <a href="/dl">Downloads</a>
			  <a href="/list">List</a>`, "\n")
	if Cred.Locked {
		fmt.Fprint(w, `<a href="/toggle">Unlock</a>`, "\n")
	} else {
		fmt.Fprint(w, `<a href="/toggle">Unlock</a>`, "\n")
	}
	fmt.Fprintf(w, `<article>
		<form method="POST" action="/changecreds">
		  <fieldset>
		  		  <legend>Change Credentials</legend>
				  <label>Username:</label><br/>
				  <input type="text" name="username"><br/>
				  <label>Password:</label><br />
				  <input type="password" name="password"><br />
				  <input type="submit">
		  </fieldset>
		</form>
	 </article>
	</main>
	</body>
	</html>`)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="en">
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Remdl</title>
	  </head>
	  <body>
		<main>
		  <a href="/dl">Downloads</a>
		  <a href="/list">List</a>
		  <a href="/creds">Creds</a>`, "\n")
	if Cred.Locked {
		fmt.Fprint(w, `<a href="/toggle">Unlock</a>`, "\n")
	} else {
		fmt.Fprint(w, `<a href="/toggle">Unlock</a>`, "\n")
	}
	fmt.Fprintf(w, `<article>
	<form method="POST" action="/add">
	  <fieldset>
			  <label>URL:</label><br />
			  <input type="text" name="url"><br />
			  <label>Filename:</label><br />
			  <input type="text" name="filename"><br />
			  <input type="submit">
	  </fieldset>
	</form>
 </article>
</main>
</body>
</html>`)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="en">
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Remdl</title>
	  </head>
	  <body>
		<main>
		  <a href="/">Home</a>
		  <a href="/dl">Downloads</a>
		  <a href="/list">Refresh</a>
		  <article>`, "\n")

	for url, val := range Engine.Responses {
		fmt.Fprintf(w, `<fieldset>
		<h3>%s</h3>
		<p>Size:      %s </p>
		<p>Complete:  %s </p>
		<p><progress value="%f" max="100"></progress> %f%%</p>
		<p>Duration:  %s </p>
		<p>Completed: %v </p>
		<p>Started %s</p>
		<p><a href="/cancel?url=%s"><button>Cancel</button></a><p>`, url, humanize.Bytes(uint64(val.Size())), humanize.Bytes(uint64(val.BytesComplete())), val.Progress(), val.Progress()*100, val.Duration(), val.IsComplete(), humanize.Time(val.Start), url)
		if val.IsComplete() {
			if val.Err() != nil {
				fmt.Fprintf(w, "\n<p>Error: %s\n", val.Err().Error())
			}
		}
		fmt.Fprint(w, "\n</fieldset>")
	}
	fmt.Fprint(w, `</article>
    </main>
  </body>
</html>`)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	if r.Method != http.MethodPost {
		return
	}
	url := r.FormValue("url")
	fn := r.FormValue("filename")
	go addreq(url, fn)
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func toggleLock(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}
	Cred.Locked = !Cred.Locked
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addreq(addurl string, filename string) {
	if len(filename) > 0 {
		filename = filepath.Join(Cred.DirPath, path.Clean(filename))
		if !strings.HasPrefix(filename, Cred.DirPath) {
			Info.Println("Invalid Path")
			return
		}
	}

	u, err := url.Parse(addurl)
	if err != nil {
		Warn.Println(err)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return
	}
	h := u.Hostname()
	if h == "0.0.0.0" || h == "127.0.0.1" {
		return
	}

	req, err := NewRequest(filename, addurl)
	if err != nil {
		Warn.Println(err)
		return
	}
	Engine.Responses[addurl] = DefaultClient.Do(req)
	Info.Println("Download added:", addurl, filename)
}

func cancelreq(w http.ResponseWriter, r *http.Request) {
	err := authHelper(w, r)
	if err != nil {
		Err.Println(err)
		return
	}

	if val := r.URL.Query().Get("url"); val != "" {
		Info.Println(val, " will be removed if present")
		resp, ok := Engine.Responses[val]
		if !ok {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		err := resp.Cancel()
		if err == context.Canceled {
			http.Error(w, "Cancelled", http.StatusInternalServerError)
		}
		if err != context.Canceled {
			http.Error(w, "Not stopped", http.StatusInternalServerError)
			return
		}
		Info.Println(resp.Filename, "Stopped")
		delete(Engine.Responses, val)
		Info.Println("Removed. Currently Downloading: ", len(Engine.Responses))
	}
}
