package images_generator_mvc

import (
	"io"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-identity/images-generator/src/images_generator/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	GENERATE_IMAGES_ENDPOINT = "/generate-images"
)

type ImagesController struct {
	logger         *zap.Logger
	generateImages use_cases.GenerateImagesInterface
}

func NewImagesController(logger *zap.Logger, generateImages use_cases.GenerateImagesInterface) *ImagesController {
	return &ImagesController{
		logger:         logger,
		generateImages: generateImages,
	}
}

func (controller *ImagesController) RegisterRoutes(app fiber.Router) {
	app.Post(GENERATE_IMAGES_ENDPOINT, controller.GenerateImages)
}

// GenerateImages godoc
// @Summary      Generate images
// @Description  Generates images based on the request payload
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        baseImageBytes formData file false "Base image file"
// @Param        otherParam formData requests.GenerateImagesRequest false "Other parameter"
// @Success      200  {object}  responses.GenerateImagesResponse
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /generate-images [post]
func (controller *ImagesController) GenerateImages(c *fiber.Ctx) error {
	controller.logger.Info("Generate images request received")

	request := new(requests.GenerateImagesRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request", zap.Error(err))
		return err
	}

	// Manually parse the file
	fileHeader, err := c.FormFile("baseImageBytes")
	if err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err == nil {
			defer file.Close()
			request.BaseImageBytes, _ = io.ReadAll(file)
			request.BaseImageMimeType = fileHeader.Header.Get("Content-Type")
		}
	}

	response, err := controller.generateImages.Generate(request)
	if err != nil {
		controller.logger.Error("Failed to generate images", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Generate images request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
