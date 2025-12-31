package main

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"

    "cups-web/internal/ipp"
    "cups-web/internal/auth"
)

type printResp struct {
    JobID string `json:"jobId,omitempty"`
    OK    bool   `json:"ok"`
}

func printHandler(w http.ResponseWriter, r *http.Request) {
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

    printer := r.FormValue("printer")
    if printer == "" {
        http.Error(w, "missing printer field", http.StatusBadRequest)
        return
    }

    // Read content (for now fully in-memory)
    var buf bytes.Buffer
    if _, err := io.Copy(&buf, file); err != nil {
        http.Error(w, "read file", http.StatusInternalServerError)
        return
    }

    // Determine mime from header, fallback
    mime := fh.Header.Get("Content-Type")
    if mime == "" {
        mime = http.DetectContentType(buf.Bytes())
    }

    sess, _ := auth.GetSession(r)

    job, err := ipp.SendPrintJob(printer, bytes.NewReader(buf.Bytes()), mime, sess.Username, fh.Filename)
    if err != nil {
        http.Error(w, "print error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(printResp{JobID: job, OK: true})
}
