package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	addr := ":8888"
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", post)
	server := &http.Server{Addr: addr, Handler: mux}
	log.Printf("Server started on: %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Failed to start server: ", err)
	}

}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse form: %v", err))
		return
	}

	file, _, err := r.FormFile("code")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get form file: %v", err))
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	b := new(bytes.Buffer)
	if _, err := io.Copy(b, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to read file: %v", err))
		return
	}

	res, err := scan(b.Bytes())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
		return
	}
	log.Println(res)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": res})
}

func scan(b []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	source := gozxing.NewLuminanceSourceFromImage(img)
	bin := gozxing.NewHybridBinarizer(source)
	bbm, err := gozxing.NewBinaryBitmap(bin)
	if err != nil {
		return "", fmt.Errorf("error during processing: %v", err)
	}

	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultiple(bbm, nil)
	if err != nil {
		return "", fmt.Errorf("unable to decode QRCode: %v", err)
	}

	var strRes []string
	for _, element := range result {
		strRes = append(strRes, element.String())
	}

	return strings.Join(strRes, "\n"), nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}
}
