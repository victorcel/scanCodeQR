package main

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"

	// import gif, jpeg, png
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"

	// import bmp, tiff, webp
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)

func main() {
	addr := "localhost:8888"
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", post)
	server := &http.Server{Addr: addr, Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}

}

// handle post request
func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		msg := fmt.Sprintf("Failed to parse form: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("code")
	if err != nil {
		msg := fmt.Sprintf("Failed to get form file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			msg := fmt.Sprintf("Failed to close file: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}(file)
	b := new(bytes.Buffer)
	if _, err := io.Copy(b, file); err != nil {
		msg := fmt.Sprintf("Failed to read file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	res, err := scan(b.Bytes())
	if err != nil {
		msg := fmt.Sprintf("Internal server error: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(res))
	if err != nil {
		msg := fmt.Sprintf("Failed to write response: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
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

	res := strings.Join(strRes, "\n")
	return res, nil
}
