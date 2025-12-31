package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Expect multipart form
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save to temp dir
	tmpDir, err := os.MkdirTemp("", "convert-")
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	inPath := filepath.Join(tmpDir, fh.Filename)
	outPath := inPath + ".pdf" // fallback output path
	f, err := os.Create(inPath)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(f, file); err != nil {
		f.Close()
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}
	f.Close()

	// Run libreoffice to convert to PDF
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "libreoffice", "--headless", "--convert-to", "pdf", "--outdir", tmpDir, inPath)
	// capture combined output for diagnostics
	if out, err := cmd.CombinedOutput(); err != nil {
		http.Error(w, "conversion failed: "+err.Error()+" - "+string(out), http.StatusInternalServerError)
		return
	}

	// Find resulting PDF (LibreOffice names file with same base and .pdf extension)
	base := fh.Filename
	ext := filepath.Ext(base)
	name := base[0 : len(base)-len(ext)]
	outPath = filepath.Join(tmpDir, name+".pdf")

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		// sometimes output may be in different name; try to find any pdf in tmpDir
		matches, _ := filepath.Glob(filepath.Join(tmpDir, "*.pdf"))
		if len(matches) > 0 {
			outPath = matches[0]
		} else {
			http.Error(w, "conversion produced no PDF", http.StatusInternalServerError)
			return
		}
	}

	// Stream PDF back
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+".pdf\"")
	pdfFile, err := os.Open(outPath)
	if err != nil {
		http.Error(w, "failed to open converted file", http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()
	if _, err := io.Copy(w, pdfFile); err != nil {
		// nothing more we can do
		return
	}
}
