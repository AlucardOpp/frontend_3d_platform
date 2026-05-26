package interfaces

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/AhEhIOhYou/etomne/pkg/server/application"
	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/utils"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces/filemanager"
	"github.com/gin-gonic/gin"
)

type Model struct {
	modelApp application.ModelAppInterface
	userApp  application.UserAppInterface
	fileApp  application.FileAppInterface
	fm       filemanager.ManagerFileInterface
	rd       auth.AuthInterface
	tk       auth.TokenInterface
}

func NewModel(mApp application.ModelAppInterface, uApp application.UserAppInterface, fApp application.FileAppInterface, fm filemanager.ManagerFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Model {
	return &Model{
		modelApp: mApp,
		userApp:  uApp,
		fileApp:  fApp,
		fm:       fm,
		rd:       rd,
		tk:       tk,
	}
}

//	@Summary	Save model
//	@Tags		Model
//	@Produce	json
//	@Param		data	body		entities.ModelRequest	true	"Model data"
//	@Success	201		{object}	entities.Model
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/model [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) SaveModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	var modelReq entities.ModelRequest

	if err := c.ShouldBindJSON(&modelReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	validateErr := modelReq.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	newModel := modelReq.NewModel()
	newModel.Prepare()
	newModel.UserID = userID

	savedModel, err := m.modelApp.SaveModel(newModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	for _, id := range modelReq.FilesId {
		if modelReq.FullDescriptionFileId != nil && id == *modelReq.FullDescriptionFileId {
			continue
		}
		err = m.modelApp.AddFileToModel(id, savedModel.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}
	}

	if modelReq.FullDescriptionFileId != nil && *modelReq.FullDescriptionFileId > 0 {
		if err := m.validateFullDescriptionFile(*modelReq.FullDescriptionFileId, userID); err != nil {
			if m.respondFullDescriptionError(c, err) {
				return
			}
		}

		if err := m.modelApp.UpdateFullDescriptionFileID(savedModel.ID, modelReq.FullDescriptionFileId); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}

		savedModel.FullDescriptionFileID = modelReq.FullDescriptionFileId
	}

	c.JSON(http.StatusCreated, savedModel)
}

//	@Summary	Update model
//	@Tags		Model
//	@Produce	json
//	@Param		model_id	path		int						true	"Model ID"
//	@Param		data		body		entities.ModelRequest	true	"Model updated data"
//	@Success	200			{object}	entities.Model
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/update/{model_id} [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) UpdateModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelID, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	updatableModel, err := m.modelApp.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(updatableModel.UserID, user, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	var modelReq entities.ModelRequest

	if err := c.ShouldBindJSON(&modelReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	hasTitle := len(modelReq.Title) != 0
	hasDescription := len(modelReq.Description) != 0
	hasKeywords := len(modelReq.Keywords) != 0
	hasFullDescription := modelReq.FullDescriptionFileId != nil

	// Если все поля пустые, возвращаем модель без изменений
	if !hasTitle && !hasDescription && !hasKeywords && !hasFullDescription {
		c.JSON(http.StatusOK, updatableModel)
		return
	}

	if hasTitle {
		updatableModel.Title = modelReq.Title
	}
	if hasDescription {
		updatableModel.Description = modelReq.Description
	}
	// Обновляем keywords в следующих случаях:
	// 1. Если keywords не пустые (пользователь хочет их изменить)
	// 2. ИЛИ если обновляются другие поля (title или description) - это позволяет очистить keywords
	// Frontend всегда отправляет keywords в запросе при обновлении модели
	if hasKeywords || hasTitle || hasDescription {
		updatableModel.Keywords = modelReq.Keywords
	}

	if hasFullDescription {
		if *modelReq.FullDescriptionFileId == 0 {
			if err := m.removeFullDescriptionFile(updatableModel); err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
				return
			}
			updatableModel.FullDescriptionFileID = nil
		} else {
			if err := m.validateFullDescriptionFile(*modelReq.FullDescriptionFileId, userID); err != nil {
				if m.respondFullDescriptionError(c, err) {
					return
				}
			}

			oldFileID := updatableModel.FullDescriptionFileID
			if oldFileID != nil && *oldFileID != *modelReq.FullDescriptionFileId {
				if err := m.removeFullDescriptionFile(updatableModel); err != nil {
					c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
					return
				}
			}

			if err := m.modelApp.UpdateFullDescriptionFileID(updatableModel.ID, modelReq.FullDescriptionFileId); err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
				return
			}
			updatableModel.FullDescriptionFileID = modelReq.FullDescriptionFileId
		}
	}

	updatableModel.BeforeUpdate()

	validateErr := updatableModel.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	updatedModel, err := m.modelApp.UpdateModel(updatableModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, updatedModel)
}

//	@Summary	Get a list of models with the specified quantity and position
//	@Tags		Model
//	@Param		_page	query	int	false	"Query page param"
//	@Param		_limit	query	int	false	"Query limit param"
//	@Param		user_id	query	int	false	"Query user ID param"
//	@Param		keywords	query	string	false	"Query keywords param (comma-separated)"
//	@Success	201		{Array}	[]entities.ModelData
//	@Failure	400		string	string
//	@Failure	401		string	string
//	@Failure	500		string	string
//	@Router		/model	[get]
func (m *Model) GetModelList(c *gin.Context) {

	limit, _ := strconv.Atoi(c.Query("_limit"))

	page, _ := strconv.Atoi(c.Query("_page"))

	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)

	keywords := c.Query("keywords")

	var readyModels []entities.ModelData

	rawModels, err := m.modelApp.GetAllModels(page, limit, userID, keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	for _, model := range rawModels {

		user, err := m.userApp.GetUser(model.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}

		modelData, err := m.buildModelData(model, user.PublicUser())
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}

		readyModels = append(readyModels, modelData)
	}

	c.JSON(http.StatusOK, readyModels)
}

//	@Summary	Get model by ID
//	@Tags		Model
//	@Param		model_id	path		int	true	"Model ID"
//	@Success	201			{object}	entities.ModelData
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/{model_id}  [get]
func (m *Model) GetModel(c *gin.Context) {
	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	model, err := m.modelApp.GetModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(model.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelData, err := m.buildModelData(*model, user.PublicUser())
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, &modelData)
}

//	@Summary	Delete model by ID
//	@Tags		Model
//	@Param		model_id	path	int	true	"Model ID"
//	@Success	200			string	string
//	@Failure	400			string	string
//	@Failure	401			string	string
//	@Failure	500			string	string
//	@Router		/model/{model_id} [delete]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) DeleteModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelID, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	model, err := m.modelApp.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(model.UserID, user, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	if err := m.removeFullDescriptionFile(model); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	files, err := m.modelApp.GetFilesByModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	// Delete files from storage
	for _, file := range files {
		err = m.fm.DeleteFile(file.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}
	}

	err = m.modelApp.DeleteModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, constants.DeletedSuccessful)
}

//	@Summary	Save model file
//	@Tags		Model
//	@Produce	json
//	@Param		model_id	path		int		true	"Model ID"
//	@Param		file		formData	file	true	"Body with files"
//	@Success	201			{object}	entities.File
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/addfile/{model_id}	[post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) SaveModelFile(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelID, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	url, err := m.fm.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if url == "" {
		c.JSON(http.StatusUnprocessableEntity, constants.FileURLError)
		return
	}

	ext := filepath.Ext(file.Filename)

	File := entities.File{
		OwnerId:   userID,
		Title:     file.Filename,
		Url:       url,
		Extension: ext,
	}

	File.Prepare()

	validateErr := File.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	modelFile, err := m.modelApp.SaveModelFile(&File, modelID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusCreated, modelFile)
}
