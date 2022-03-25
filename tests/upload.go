package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type ImgurAPIResponse struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}
type Data struct {
	ID          string      `json:"id"`
	Title       interface{} `json:"title"`
	Description interface{} `json:"description"`
	Datetime    int         `json:"datetime"`
	Type        string      `json:"type"`
	Animated    bool        `json:"animated"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	Size        int         `json:"size"`
	Views       int         `json:"views"`
	Bandwidth   int         `json:"bandwidth"`
	Deletehash  string      `json:"deletehash"`
	Name        string      `json:"name"`
	Link        string      `json:"link"`
}

func uploadImage(imageData []byte, imageName string) (imgURL, imageDeleteHash string, err error) {
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)
	imgFileWriter, err := multipartWriter.CreateFormFile("image", imageName)
	if err != nil {
		return
	}
	_, err = imgFileWriter.Write(imageData)
	if err != nil {
		return
	}
	err = multipartWriter.Close()
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &buffer)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Client-ID "+os.Getenv("IMGUR_CLIENT_ID"))
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var imgurAPIResp ImgurAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&imgurAPIResp)
	if err != nil {
		return
	}
	if !imgurAPIResp.Success {
		err = fmt.Errorf("imgur API error: %v", imgurAPIResp.Status)
	}
	return imgurAPIResp.Data.Link, imgurAPIResp.Data.Deletehash, nil
}
