package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}



func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	//r.ParseMultipartForm(32 << 20) // 32Mb
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file

	//f, err := os.OpenFile("uploads/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)


	dst, err := os.Create("uploads/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>Successfully Uploaded File</h1>\n")

	

	fmt.Fprintf(w, `<form   enctype="multipart/form-data"   action="http://localhost:8080/upload"   method="get" >`)
	fmt.Fprintf(w, `<h3><input type="submit" value="upload again" /></h3></form>`)
}



func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}