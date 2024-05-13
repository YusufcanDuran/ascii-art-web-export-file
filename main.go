package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		http.Error(w, "Missing text or banner", http.StatusBadRequest)
		return
	}

	turkcechar := "ışğüöçİŞÖÇĞÜ"
	if strings.ContainsAny(text, turkcechar) {
		http.Error(w, "Turkish characters detected", http.StatusBadRequest)
		return
	}

	asciiArt, err := generateAsciiArt(text, banner)
	if err != nil {
		http.Error(w, "Error generating ASCII Art", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	fmt.Fprintf(w, "%s", asciiArt)
}

func generateAsciiArt(text, banner string) (string, error) {
	cmd := exec.Command("go", "run", "argument.go", text, banner)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running command: %v", err)
	}
	return string(output), nil
}

func main() {
	http.Handle("/templates/style.css", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate", asciiArtHandler)
	http.HandleFunc("/export", asciiArtExportHandler)
	http.ListenAndServe(":8080", nil)
}

func asciiArtExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		http.Error(w, "Missing text or banner", http.StatusBadRequest)
		return
	}

	turkcechar := "ışğüöçİŞÖÇĞÜ"
	if strings.ContainsAny(text, turkcechar) {
		http.Error(w, "Turkish characters detected", http.StatusBadRequest)
		return
	}

	asciiArt, err := generateAsciiArt(text, banner)
	if err != nil {
		http.Error(w, "Error generating ASCII Art", http.StatusInternalServerError)
		return
	}

	// Dosyaya yazma işlemi
	file, err := os.Create("ascii_art.txt")
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(asciiArt)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	// Dosyayı indirme işlemi
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(asciiArt)))
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_art.txt")
	http.ServeFile(w, r, "ascii_art.txt")
}
