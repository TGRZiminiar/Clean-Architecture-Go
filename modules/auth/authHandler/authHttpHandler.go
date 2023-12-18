package authhandler

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authusecase "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authUsecase"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/request"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/response"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type (
	AuthHttpHandlerService interface {
		CreateUser(c *fiber.Ctx) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authusecase.AuthUsecaseService
	}
)

type FileReq struct {
	File        *multipart.FileHeader `form:"file"`
	Destination string                `form:"destination"`
	Extension   string
	FileName    string
}

type FileRes struct {
	FileName string `json:"filename"`
	Url      string `json:"url"`
}

type DeleteFileReq struct {
	Destination string `json:"destination"`
}

type filesPub struct {
	bucket      string
	destination string
	file        *FileRes
}

func NewAuthHttpHandler(cfg *config.Config, authUsecase authusecase.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
	}
}

func (h *authHttpHandler) CreateUser(c *fiber.Ctx) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(auth.CreateUser)
	if err := wrapper.ParseJson(req); err != nil {
		return response.ErrorRes(c, http.StatusBadRequest, err.Error())
	}

	if errs := wrapper.Validate(req); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		return response.ErrorRes(c, http.StatusBadRequest, strings.Join(errMsgs, " and "))
	}

	token, user, err := h.authUsecase.CreateUser(h.cfg, ctx, req)
	if err != nil {
		return response.ErrorRes(c, http.StatusBadRequest, err.Error())
	}

	oneWeek := time.Now().Add(7 * 24 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		Expires:  oneWeek,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessChecker",
		Value:    user.Id.String(),
		HTTPOnly: false,
		Secure:   false,
		Expires:  oneWeek,
	})

	return response.SuccessRes(c, 201, user)
}

func (h *authHttpHandler) UploadImageUser(c *fiber.Ctx) error {

	req := make([]*FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return response.ErrorRes(c, fiber.StatusBadRequest, fmt.Sprintf("type of form might be wrong : %s", err.Error()))
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return response.ErrorRes(c, fiber.ErrBadRequest.Code, fmt.Sprintf("extension is not acceptable : %s", err.Error()))
		}
		if file.Size > int64(2097152) {
			return response.ErrorRes(c, fiber.ErrBadRequest.Code, fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(2097152)/math.Pow(1024, 2)))))
		}

		filename := utils.RandFileName(ext)
		req = append(req, &FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := UploadToStorage(req)
	if err != nil {
		return response.ErrorRes(c, fiber.ErrInternalServerError.Code, err.Error())
	}
	return response.SuccessRes(c, fiber.StatusCreated, res)
}

func UploadToStorage(req []*FileReq) ([]*FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *FileReq, len(req))
	resultsCh := make(chan *FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go uploadToStorageWorker(ctx, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}
	return res, nil
}

func uploadToStorageWorker(ctx context.Context, jobs <-chan *FileReq, results chan<- *FileRes, errs chan<- error) {
	for job := range jobs {
		cotainer, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}
		b, err := ioutil.ReadAll(cotainer)
		if err != nil {
			errs <- err
			return
		}

		// Upload an object to storage
		dest := fmt.Sprintf("./assets/images/%s", job.Destination)
		if err := os.WriteFile(dest, b, 0777); err != nil {
			if err := os.MkdirAll("./assets/images/"+strings.Replace(job.Destination, job.FileName, "", 1), 0777); err != nil {
				errs <- fmt.Errorf("mkdir \"./assets/images/%s\" failed: %v", err, job.Destination)
				return
			}
			if err := os.WriteFile(dest, b, 0777); err != nil {
				errs <- fmt.Errorf("write file failed: %v", err)
				return
			}
		}

		newFile := &filesPub{
			file: &FileRes{
				FileName: job.FileName,
				// host and post
				Url: fmt.Sprintf("http://%s:%d/%s", "mix", "5000", job.Destination),
			},
			destination: job.Destination,
		}

		errs <- nil
		results <- newFile.file
	}
}
