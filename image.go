package ical

import (
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path"
	"strings"
)

type Image struct {
	URL      string
	MIMEType string
	// BASE64 encoded string
	Value string
	// Optional: Defaults to DisplayBadge
	Display DisplayType
}

func (img *Image) generate(builder *strings.Builder) error {
	err := img.valid()
	if err != nil {
		return err
	}

	if img.Display == "" {
		img.Display = DisplayBadge
	}

	if img.URL != "" {
		if img.MIMEType == "" {
			ext := path.Ext(img.URL)
			img.MIMEType = mime.TypeByExtension(ext)
			if img.MIMEType == "" {
				img.MIMEType = "image/png" //Fallback
			}
		}

		builder.WriteString(fmt.Sprintf("IMAGE;VALUE=URI;DISPLAY=%s;FMTTYPE=%s:%s", img.Display, img.MIMEType, img.URL) + lineBreak)
	} else if img.Value != "" {
		if img.MIMEType == "" {
			img.MIMEType = "image/png" //Fallback
		}

		base64 := foldText(img.Value)
		builder.WriteString(fmt.Sprintf("IMAGE;FMTTYPE=%s;DISPLAY=%s;ENCODING=BASE64;VALUE=BINARY:%s", img.MIMEType, img.Display, base64) + lineBreak)
	} else {
		return ErrImageNotUrlOrBase64
	}
	return nil
}

func (img *Image) valid() error {
	if img.URL == "" && img.Value == "" {
		return ErrImageNotUrlOrBase64
	}

	if img.URL != "" && path.Ext(img.URL) == "" {
		return ErrImageNotUrlOrBase64
	}

	return nil
}

func ImageFromURL(url string) *Image {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return &Image{
			URL: url,
		}
	} else {
		if strings.Contains(url, "://") {
			return nil // Invalid URL scheme
		}

		url = fmt.Sprintf("https://%s", url)
		return &Image{
			URL: url,
		}
	}
}

func ImageFromBase64(mimeType string, base64 string) (*Image, error) {
	if mimeType == "" {
		return nil, ErrInvalidMimeType
	}

	if base64 == "" {
		return nil, ErrImageNotUrlOrBase64
	}

	return &Image{
		MIMEType: mimeType,
		Value:    base64,
	}, nil
}

func ImageFromFile(filePath string) (*Image, error) {
	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return nil, ErrInvalidMimeType
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)

	return &Image{
		MIMEType: mimeType,
		Value:    string(encoded),
	}, nil
}

type DisplayType string

const (
	DisplayBadge     DisplayType = "BADGE"
	DisplayThumbnail DisplayType = "THUMBNAIL"
	DisplayFull      DisplayType = "FULL"
)
