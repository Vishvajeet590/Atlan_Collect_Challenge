package formHandler

import (
	"Atlan_Collect_Challenge/entity"
	"Atlan_Collect_Challenge/usecase/form"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FormHandler struct {
	Service *form.Service
}

type formId struct {
	formId int `json:"form_id"`
}

type userId struct {
	userId int `json:"user_id"`
}

type message struct {
	Message string `json:"message"`
}

type formReq struct {
	FormName string        `json:"form_name" binding:"required"`
	FormId   int8          `json:"form_id"`
	OwnerId  int           `json:"owner_id" binding:"required"`
	Question []questionReq `json:"question" binding:"required"`
}

type questionReq struct {
	QuestionId   int    `json:"question_id"`
	Question     string `json:"question" binding:"required"`
	QuestionType string `json:"question_type" binding:"required"`
}

var TAG = "FormHandler"

func NewFormHandler(Service *form.Service) *FormHandler {
	return &FormHandler{
		Service: Service,
	}
}

// Collect godoc
// @Summary Create Forms.
// @Description Endpoint to create form.
// @Tags form
// @Accept json
// @Produce json
// @Param form body formReq true "Form"
// @Success 200 {object} formId
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /form/create [post]
func (s *FormHandler) CreateForm(c *gin.Context) {
	var body entity.Form
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form := &entity.Form{FormName: body.FormName, OwnerId: body.OwnerId, Question: body.Question}
	//TODO CREATE ID
	res, formId, err := s.Service.CreateForm(form)
	if err != nil {
		return
	}
	if res == false {
		c.JSON(http.StatusInternalServerError, gin.H{"Unexpected error": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Form created. formId = ": formId})

}

// Collect godoc
// @Summary Get Forms.
// @Description Endpoint to retrive form.
// @Tags form
// @Param formId path string true "Form id"
// @Produce json
// @Success 200 {object} formReq
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /form/{formId} [get]
func (s *FormHandler) GetForm(c *gin.Context) {
	var formId = c.Param("formId")
	id, err := strconv.Atoi(formId)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusBadRequest, gin.H{"incorrect Form id": false})
		return
	}

	form, err := s.Service.GetForm(int8(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Unexpected error : ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Form": form})
}
