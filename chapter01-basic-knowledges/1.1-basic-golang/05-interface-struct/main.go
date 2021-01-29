package main

import (
	"fmt"
)

// ICitizenService is the interface for citizen service
type ICitizenService interface {
	Validate(c *Citizen) bool
	CreateCitizenCard(c *Citizen) error
}

// ThaiCitizenService is the service to work with citizen data
type ThaiCitizenService struct {
}

// NewThaiCitizenService is constructor function for ThaiCitizenService
func NewThaiCitizenService() *ThaiCitizenService {
	return &ThaiCitizenService{}
}

// Validate validate citizen information
func (svc *ThaiCitizenService) Validate(c *Citizen) bool {
	return len(c.Firstname) > 0 && len(c.Lastname) > 0 && len(c.CitizenID) > 0
}

// CreateCitizenCard will request API to create citizen card if citizen information is valid
func (svc *ThaiCitizenService) CreateCitizenCard(c *Citizen) error {
	// TODO: Request API to create citizen card
	println(fmt.Sprintf("Successfully create Thai citizen card for ID=%s", c.CitizenID))
	return nil
}

// JapanCitizenService is the service to work with citizen data
type JapanCitizenService struct {
}

// NewJapanCitizenService is constructor function for JapanCitizenService
func NewJapanCitizenService() *JapanCitizenService {
	return &JapanCitizenService{}
}

// Validate validate citizen information
func (svc *JapanCitizenService) Validate(c *Citizen) bool {
	return len(c.Firstname) > 0 && len(c.Lastname) > 0 && len(c.CitizenID) > 0
}

// CreateCitizenCard will request API to create citizen card if citizen information is valid
func (svc *JapanCitizenService) CreateCitizenCard(c *Citizen) error {
	// TODO: Request API to create citizen card
	println(fmt.Sprintf("Successfully create Japan citizen card for ID=%s", c.CitizenID))
	return nil
}

func main() {

	// 1. Constructor function will return struct pointer
	citizen := NewCitizen("Chaiyapong", "Lapliengtrakul", "1122334455")
	citizenSvc := NewThaiCitizenService()
	// citizenSvc := NewJapanCitizenService()

	// 2. Function that accept interface, the argument must implement function declare in interface
	//    In this case CitizenService struct must implement ICitizenService
	err := createCitizenCard(citizenSvc, citizen)
	if err != nil {
		printError(err)
	}
	printUnderline()

	// 3. Struct that doesn't implement interface cannot send as argument to the function with interface
	// anotherSvc := NewAnotherService()
	// createCitizenCard(anotherSvc, citizen) // This will error
}

func createCitizenCard(svc ICitizenService, c *Citizen) error {
	if !svc.Validate(c) {
		return fmt.Errorf("Citizen data is invalid")
	}

	err := svc.CreateCitizenCard(c)
	if err != nil {
		return err
	}

	return nil
}

// Citizen is type represent person in country
type Citizen struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	CitizenID string `json:"citizen_id"`
}

// NewCitizen is constructor function for Citizen
func NewCitizen(firstname string, lastname string, citizenID string) *Citizen {
	return &Citizen{
		Firstname: firstname,
		Lastname:  lastname,
		CitizenID: citizenID,
	}
}

// AnotherService is another service
type AnotherService struct{}

// NewAnotherService return new instance of another service
func NewAnotherService() *AnotherService {
	return &AnotherService{}
}

func printError(err error) {
	fmt.Println("error = ", err.Error())
}

func printUnderline() {
	fmt.Println("---")
}
