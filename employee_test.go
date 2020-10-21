package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangePassword(t *testing.T) {
	assert := assert.New(t)
	cpassword := ChangePassword{}

	response := postRequest(t, "/change_password", cpassword)
	expectedOutput := `{"status":"Fail","message":"Password changed failed","error":{"type":"Invalid Value","detail":"Required OldPassword"}}`
	assert.Equal(response.Body.String(), expectedOutput, "OK")

	cpassword = ChangePassword{OldPassword: "nopassword"}
	response = postRequest(t, "/change_password", cpassword)
	expectedOutput = `{"status":"Fail","message":"Password changed failed","error":{"type":"Invalid Value","detail":"Required NewPassword"}}`
	assert.Equal(response.Body.String(), expectedOutput, "Ok")

	cpassword = ChangePassword{OldPassword: "saven", NewPassword: "sago"}
	response = postRequest(t, "/change_password", cpassword)
	expectedOutput = `{"status":"Fail","message":"Password changed failed","error":{"type":"Invalid Data","detail":"Invalid OldPassword"}}`
	assert.Equal(response.Body.String(), expectedOutput, "Ok")

	cpassword = ChangePassword{OldPassword: "nopassword", NewPassword: "sago"}
	response = postRequest(t, "/change_password", cpassword)
	expectedOutput = `{"status":"Success","message":"Password changed successfully","error":{"type":"","detail":""}}`
	assert.Equal(response.Body.String(), expectedOutput, "Ok")

}

func TestUploadImage(t *testing.T) {
	assert := assert.New(t)
	uimage := UploadImage{}

	response := postRequest(t, "/upload_image", uimage)
	expectedOutput := `{"status":"Fail","message":"Image upload failed","error":{"type":"Invalid Value","detail":"Required Image"}}`
	assert.Equal(response.Body.String(), expectedOutput, "OK")

	uimage = UploadImage{Image: "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="}
	response = postRequest(t, "/upload_image", uimage)
	expectedOutput = `{"status":"Success","message":"Image uploaded successfully","error":{"type":"","detail":""}}`
	assert.Equal(response.Body.String(), expectedOutput, "OK")

}

func postRequest(t *testing.T, url string, input interface{}) *httptest.ResponseRecorder {
	data, _ := json.Marshal(input)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		t.Log(err)
	}
	request.Header.Set("Content-type", "application/json")

	response := httptest.NewRecorder()
	if url == "/upload_image" {
		handler := http.HandlerFunc(uploadImage)
		handler.ServeHTTP(response, request)
	} else if url == "/change_password" {
		handler := http.HandlerFunc(changePassword)
		handler.ServeHTTP(response, request)
	}
	return response
}
