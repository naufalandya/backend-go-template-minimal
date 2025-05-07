package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	pbapi "modular_monolith/protobuf/api"
	"modular_monolith/server/api/legal/functions"
	legal "modular_monolith/server/api/legal/models"
	"modular_monolith/server/api/legal/provider"
	"modular_monolith/server/api/legal/utils"
	"modular_monolith/server/client"
	common "modular_monolith/server/functions"
	global "modular_monolith/server/models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadFile(c *fiber.Ctx) error {

	var input legal.UploadFileRequest

	tglPenetapanStr := c.FormValue("tgl_penetapan")
	tglBerlakuStr := c.FormValue("tgl_berlaku")
	layout := "2006-01-02"

	tglPenetapan, err := time.Parse(layout, tglPenetapanStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "Invalid tgl_penetapan~ (；⌣̀_⌣́)",
		})
	}

	tglBerlaku, err := time.Parse(layout, tglBerlakuStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "Invalid tgl_berlaku~ (；⌣̀_⌣́)",
		})
	}

	folderIDStr := c.FormValue("folderId")
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "Invalid folderId~ (`･ω･´)",
		})
	}

	pemrakarsaStr := c.FormValue("pemrakarsa_id")
	pemrakarsaID, err := strconv.Atoi(pemrakarsaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: "Invalid pemrakarsa_id~ ヽ(ﾟДﾟ)ﾉ",
		})
	}

	// parse revoke
	var revoke []legal.RelatedFile
	revokeStr := c.FormValue("revoke")
	if revokeStr != "" {
		if err := json.Unmarshal([]byte(revokeStr), &revoke); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
				Code:    fiber.StatusBadRequest,
				Status:  false,
				Message: "Invalid revoke format~ (°ロ°) !",
			})
		}
	}

	fmt.Println("revoke", revoke)

	var modify []legal.RelatedFile
	modifyStr := c.FormValue("modify")
	if modifyStr != "" {
		if err := json.Unmarshal([]byte(modifyStr), &modify); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
				Code:    fiber.StatusBadRequest,
				Status:  false,
				Message: "Invalid modify format~ (°ロ°) !",
			})
		}
	}

	input = legal.UploadFileRequest{
		FolderID:      folderID,
		Title:         c.FormValue("title"),
		Nomor:         c.FormValue("nomor"),
		Tahun:         c.FormValue("tahun"),
		TglPenetapan:  &tglPenetapan,
		TglBerlaku:    &tglBerlaku,
		Pemrakarsa_ID: pemrakarsaID,
		Revoke:        revoke,
		Modify:        modify,
	}

	if err := common.FuckOffHackerByJSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Suspicious input detected~ (｀_´) : %s", err.Error()),
		})
	}

	if errs := common.ValidateStruct(input); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Message: fmt.Sprintf("Validation failed~ (｀_´) : %s", errs[0]),
		})
	}

	err = provider.CheckFolderExists(input.FolderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(global.Apiresponse{
			Message: "Folder is not available",
			Code:    fiber.StatusNotFound,
			Status:  false,
			Data:    nil,
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Message: "File is required",
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Data:    nil,
		})
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" && ext != ".docx" {
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Message: "Only .docx and .pdf files are allowed",
			Code:    fiber.StatusBadRequest,
			Status:  false,
			Data:    nil,
		})
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/app/uploads"
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to create upload directory: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}

	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to save file: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}

	file, err := os.Open(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to reopen saved file: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}
	defer file.Close()

	fileType := utils.GetFileType(fileHeader.Filename)
	var text string
	switch fileType {
	case "pdf":
		text, err = functions.ExtractTextFromPDF(file)
	case "docx":
		text, err = functions.ExtractTextFromDocx(file)
	default:
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(global.Apiresponse{
			Message: "Unsupported file type",
			Code:    fiber.StatusUnsupportedMediaType,
			Status:  false,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to extract text: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}

	docID, err := provider.InsertTextToElastic(text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Failed to insert text into Elasticsearch: %v", err),
			Code:    fiber.StatusInternalServerError,
			Status:  false,
			Data:    nil,
		})
	}

	userID := c.Locals("id")
	fmt.Printf("userID type: %T\n", userID)
	var userIDInt int
	if userIDFloat, ok := userID.(float64); ok {
		userIDInt = int(userIDFloat)
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Status:  false,
			Message: "Failed to parse user ID",
			Data:    nil,
			Code:    fiber.StatusBadRequest,
		})
	}

	request := &legal.UploadFileRequest{
		FolderID:     input.FolderID,
		Title:        input.Title,
		Nomor:        input.Nomor,
		Tahun:        input.Tahun,
		TglBerlaku:   input.TglBerlaku,
		TglPenetapan: input.TglPenetapan,
	}

	uuidv7, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Status:  false,
			Message: "Failed uploading file",
			Data:    nil,
			Code:    fiber.StatusBadRequest,
		})
	}

	metadata := &legal.File{
		Id:         uuidv7.String(),
		FileName:   fileHeader.Filename,
		UploadedBy: int32(userIDInt),
		Size:       int32(fileHeader.Size),
		Type:       fileType,
		Path:       filePath,
		IsTrash:    false,
		Download:   fmt.Sprintf("https://alobro.my.id/images/%s", fileHeader.Filename),
		IndexText:  docID,
	}

	fileID, err := provider.InsertFile(request, metadata)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Status:  false,
			Message: "Failed to insert file",
			Code:    fiber.StatusInternalServerError,
		})
	}

	var modifyProto []*pbapi.RelatedFile
	for _, m := range modify {
		modifyProto = append(modifyProto, &pbapi.RelatedFile{
			Reason:       m.Reason,
			TargetFileId: m.TargetFileID,
		})
	}

	var revokeProto []*pbapi.RelatedFile
	for _, m := range revoke {
		revokeProto = append(revokeProto, &pbapi.RelatedFile{
			Reason:       m.Reason,
			TargetFileId: m.TargetFileID,
		})
	}

	res, err := client.Clients.LegalServiceClient.UploadFile(
		c.Context(), &pbapi.UploadFileRequest{
			FileId: fileID,
			Revoke: revokeProto,
			Modify: modifyProto,
		},
	)

	if err != nil {
		log.Printf("gRPC UploadFile error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(global.Apiresponse{
			Message: "Failed to upload file to legal service (╯°□°）╯︵ ┻━┻",
			Code:    fiber.StatusInternalServerError,
			Status:  false,
		})
	}

	if res.GetStatus() == false {
		log.Printf("gRPC UploadFile failed: %v", res.GetMessage())
		return c.Status(fiber.StatusBadRequest).JSON(global.Apiresponse{
			Message: fmt.Sprintf("Legal service rejected the file: %s", res.GetMessage()),
			Code:    fiber.StatusBadRequest,
			Status:  false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(global.Apiresponse{
		Message: "Text extracted and inserted into Elasticsearch successfully",
		Code:    fiber.StatusOK,
		Status:  true,
		Data: map[string]interface{}{
			"id":       docID,
			"filePath": filePath,
		},
	})

}
