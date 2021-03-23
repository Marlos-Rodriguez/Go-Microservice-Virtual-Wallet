package handlers

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	//internalJWT "github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/internal/jwt"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/images/internal/environment"

	//"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	//AWS
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IimagesHandlerService interface {
	ChangeAvatar(c *fiber.Ctx) error
}

type ImagesHandlerService struct {
	AmazonSession *session.Session
}

var (
	AWS_S3_REGION = environment.AccessENV("AWS_S3_REGION")
	AWS_S3_BUCKET = environment.AccessENV("AWS_S3_BUCKET")
)

//NewImageshandlerService Create a new handler to work with form-data and AWS
func NewImageshandlerService() IimagesHandlerService {
	go client.StartClient()

	newSession, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		return nil
	}

	return &ImagesHandlerService{AmazonSession: newSession}
}

func (s *ImagesHandlerService) ChangeAvatar(c *fiber.Ctx) error {
	//Get the ID
	ID := c.Params("id")

	if len(ID) < 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Review your input"})
	}

	/*
		//Check the JWT ID
		tk := c.Locals("user").(*jwt.Token)
		if err := internalJWT.GetClaims(*tk); err != nil {
			return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
		}

		if match, err := internalJWT.CheckID(ID); !match || err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "Error in process JWT", "data": err.Error()})
		}*/

	//Get the Image from form-Data
	fileHeader, err := c.FormFile("avatar")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in get the image", "data": err.Error()})
	}

	var extension string = fileHeader.Header.Get("Content-Type")
	var size float64 = float64((fileHeader.Size / 1024) / 1024)

	if size > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in get the image"})
	}
	//compress and convert image

	var initImage image.Image

	File, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in open the image", "data": err.Error()})
	}
	defer File.Close()

	//Check type of image
	if extension == "image/png" {
		initImage, err = png.Decode(File)
	} else if extension == "image/jpeg" || extension == "image/jpg" {
		initImage, err = jpeg.Decode(File)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Image format not supported"})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in get the image", "data": err.Error()})
	}

	options := &jpeg.Options{
		Quality: 70,
	}

	var byteContainer bytes.Buffer

	byteContainer.ReadFrom(File)

	err = jpeg.Encode(&byteContainer, initImage, options)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in compress the image", "data": err.Error()})
	}

	//upload to AWS
	uploader := s3manager.NewUploader(s.AmazonSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket to be used
		Key:    aws.String(ID + ".jpeg"),  // Name of the file to be saved
		Body:   &byteContainer,            // File
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in sending to AWS S3", "data": err.Error()})
	}

	//Call to User service to change
	var url string = "https://" + AWS_S3_BUCKET + "." + "s3-" + AWS_S3_REGION + ".amazonaws.com/" + ID + ".jpeg"

	//User gRPC Client
	if success, err := client.ChangeAvatar(url, ID); !success || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error in Updated your user info", "data": err.Error()})
	}

	return c.SendStatus(fiber.StatusAccepted)
}
