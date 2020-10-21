package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

// Response model used to send service response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   Error  `json:"error,omitempty"`
}

// Error model used to prepare error information
type Error struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

// Employee model used to maintain employee details
type Employee struct {
	ProfileImage string `json:"profile_image"`
	Password     string `json:"password"`
}

// ChangePassword model used to read input from request
type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UploadImage model used to read input from request
type UploadImage struct {
	Image string `json:"image"`
}

var (
	employee = Employee{ProfileImage: "", Password: "nopassword"}
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/upload_image", uploadImage)
	router.HandleFunc("/change_password", changePassword)

	handler := cors.AllowAll().Handler(router)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Printf("failed to start server")
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	response := Response{}

	if r.Method == "POST" {
		uimage := UploadImage{}
		err := json.NewDecoder(r.Body).Decode(&uimage)
		if err != nil {
			response = Response{Status: "Fail", Message: "Image upload failed", Error: Error{Type: "Invalid Request", Detail: "Invalid json"}}
		}

		if response.Status != "Fail" {
			if uimage.Image == "" {
				response = Response{Status: "Fail", Message: "Image upload failed", Error: Error{Type: "Invalid Value", Detail: "Required Image"}}
			}

			if response.Status != "Fail" {
				employee.ProfileImage = uimage.Image
				response = Response{Status: "Success", Message: "Image uploaded successfully", Error: Error{}}
			}
		}

	} else {
		response = Response{Status: "Fail", Message: "Image upload failed", Error: Error{Type: "Invalid Request", Detail: "Invalid http method"}}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Fprintf(w, string(resp))
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	response := Response{}

	if r.Method == "POST" {
		cpassowrd := ChangePassword{}
		err := json.NewDecoder(r.Body).Decode(&cpassowrd)
		if err != nil {
			response = Response{Status: "Fail", Message: "Password changed failed", Error: Error{Type: "Invalid Request", Detail: "Invalid json"}}
		}

		if response.Status != "Fail" {

			if cpassowrd.OldPassword == "" {
				response = Response{Status: "Fail", Message: "Password changed failed", Error: Error{Type: "Invalid Value", Detail: "Required OldPassword"}}
			} else if cpassowrd.NewPassword == "" {
				response = Response{Status: "Fail", Message: "Password changed failed", Error: Error{Type: "Invalid Value", Detail: "Required NewPassword"}}
			}

			if response.Status != "Fail" {
				if cpassowrd.OldPassword != employee.Password {
					response = Response{Status: "Fail", Message: "Password changed failed", Error: Error{Type: "Invalid Data", Detail: "Invalid OldPassword"}}
				}

				if response.Status != "Fail" {
					employee.Password = cpassowrd.NewPassword
					response = Response{Status: "Success", Message: "Password changed successfully", Error: Error{}}
				}
			}
		}

	} else {
		response = Response{Status: "Fail", Message: "Password changed failed", Error: Error{Type: "Invalid Request", Detail: "Invalid http method"}}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Fprintf(w, string(resp))
}
