package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type apiResp struct {
	Images map[string]string `json:"images"`
}

func main() {
	// ---------- flags ----------
	url := flag.String("url", "http://localhost:5000/generate-images", "API endpoint")
	outDir := flag.String("o", defaultDownloads(), "Output directory for decoded images")

	name := flag.String("name", "pierre", "form field: name")
	role := flag.String("role", "expert", "form field: role")
	ensDomain := flag.String("ensDomain", "pierre.agent.protocol6022.eth", "form field: ensDomain")
	creator := flag.String("creator", "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146", "form field: creator (wallet addr)")
	walletsStr := flag.String("wallets", "0xf5cfC4723d64C0B53b3aF49C3a184f16E7cACA40", "form field: wallets (comma-sep)")
	stars := flag.Int("stars", 7, "form field: stars (int)")
	filePath := flag.String("file", "./avogadro.png", "PNG to send as baseImageBytes")
	mimeType := flag.String("mimetype", "image/png", "form field: baseImageMimeType")

	flag.Parse()

	// ---------- build multipart body ----------
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	// text fields
	writeField(w, "name", *name)
	writeField(w, "role", *role)
	writeField(w, "creator", *creator)
	writeField(w, "ensDomain", *ensDomain)
	writeField(w, "stars", itoa(*stars))
	writeField(w, "baseImageMimeType", *mimeType)

	wallets := strings.Split(*walletsStr, ",")
	for _, wallet := range wallets {
		writeField(w, "wallets", strings.TrimSpace(wallet))
	}

	// file field
	if *filePath != "" {
		fw, err := w.CreateFormFile("baseImageBytes", filepath.Base(*filePath))
		check(err)
		f, err := os.Open(*filePath)
		check(err)
		defer f.Close()
		_, err = io.Copy(fw, f)
		check(err)
	}

	check(w.Close())

	// ---------- HTTP call ----------
	req, err := http.NewRequest("POST", *url, &body)
	check(err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	check(err)
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	check(err)

	// ---------- decode & save ----------
	var r apiResp
	check(json.Unmarshal(raw, &r))

	check(os.MkdirAll(*outDir, 0o755))
	for name, b64 := range r.Images {
		data, err := base64.StdEncoding.DecodeString(b64)
		check(err)

		// force extension by key
		var filename string
		switch name {
		case "default":
			filename = "default.svg"
		case "icon":
			filename = "icon.png"
		default:
			filename = name // fallback to raw key
		}

		dst := filepath.Join(*outDir, filename)
		check(os.WriteFile(dst, data, 0o644))
		println("✓ saved", dst)
	}
}

func writeField(w *multipart.Writer, key, val string) {
	if val == "" {
		return
	}
	fw, err := w.CreateFormField(key)
	check(err)
	_, err = fw.Write([]byte(val))
	check(err)
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }

func defaultDownloads() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Downloads")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
