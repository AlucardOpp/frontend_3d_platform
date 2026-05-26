package interfaces

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/utils"
	"github.com/gin-gonic/gin"
)

var (
	errFullDescriptionFileInvalid = errors.New("full description file format is not supported")
	errFullDescriptionFileAccess  = errors.New("full description file is not available")
)

func (m *Model) validateFullDescriptionFile(fileID, userID uint64) error {
	file, err := m.fileApp.GetFile(fileID)
	if err != nil {
		return err
	}

	if file.OwnerId != userID {
		return errFullDescriptionFileAccess
	}

	if !utils.IsFullDescriptionExtension(file.Extension) {
		return errFullDescriptionFileInvalid
	}

	return nil
}

func (m *Model) deleteFileFromStorage(file *entities.File) error {
	return m.fm.DeleteFile(file.Url)
}

func (m *Model) removeFullDescriptionFile(model *entities.Model) error {
	if model.FullDescriptionFileID == nil || *model.FullDescriptionFileID == 0 {
		return nil
	}

	file, err := m.fileApp.GetFile(*model.FullDescriptionFileID)
	if err != nil {
		return err
	}

	if err := m.deleteFileFromStorage(file); err != nil {
		return err
	}

	if err := m.fileApp.DeleteFile(file.ID); err != nil {
		return err
	}

	return m.modelApp.UpdateFullDescriptionFileID(model.ID, nil)
}

func (m *Model) buildModelData(model entities.Model, user *entities.PublicUser) (entities.ModelData, error) {
	files, err := m.modelApp.GetFilesByModel(model.ID)
	if err != nil {
		return entities.ModelData{}, err
	}

	filteredFiles := utils.FilterOutFullDescriptionFile(files, model.FullDescriptionFileID)

	modelData := entities.ModelData{
		Model: model,
		User:  *user,
		Files: utils.SortFiles(filteredFiles),
	}

	if model.FullDescriptionFileID != nil && *model.FullDescriptionFileID > 0 {
		fullDescription, err := m.fileApp.GetFile(*model.FullDescriptionFileID)
		if err == nil {
			modelData.FullDescription = fullDescription
		}
	}

	return modelData, nil
}

func (m *Model) respondFullDescriptionError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch err {
	case errFullDescriptionFileInvalid:
		c.JSON(http.StatusBadRequest, constants.FullDescriptionFileFormatError)
		return true
	case errFullDescriptionFileAccess:
		c.JSON(http.StatusBadRequest, constants.FileNotAvaliable)
		return true
	default:
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return true
	}
}
