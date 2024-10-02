package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	timestamp := time.Now().UnixMilli()

	secretStr := "sjsjsjsjsjs"

	secretWithTimestamp := []byte(strconv.FormatInt(timestamp, 10) + "\n" + secretStr)
	secret := []byte(secretStr)

	h := hmac.New(sha256.New, secret)
	h.Write(secretWithTimestamp)
	hmacSign := h.Sum(nil)

	b64Sign := base64.StdEncoding.EncodeToString(hmacSign)
	parseSign := url.QueryEscape(b64Sign)

	data := url.Values{}
	data.Set("timestamp", strconv.FormatInt(timestamp, 10))
	data.Set("sign", parseSign)
	data.Set("from", "1234567890")
	data.Set("content", "Your sms verify code is 182734")

	req, err := http.NewRequest("POST", "http://localhost:18080/post", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status OK; got %v and %s", resp.Status, string(body))
	}

}
