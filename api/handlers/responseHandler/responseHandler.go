package responseHandler

import (
	"Atlan_Collect_Challenge/entity"
	"Atlan_Collect_Challenge/usecase/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type formId struct {
	formId int `json:"form_id"`
}

type userId struct {
	userId int `json:"user_id"`
}

type message struct {
	Message string `json:"message"`
}

type ResponseReq struct {
	QuestionId   int    `json:"question_id"`
	Response     string `json:"response"`
	ResponseType string `json:"response_type"`
}

type ResponseRequest struct {
	ResponseId int           `json:"response_id"`
	Responses  []ResponseReq `json:"responses"`
}

type ResponseHandler struct {
	Service *response.Service
}

func NewResponseHandler(Service *response.Service) *ResponseHandler {
	return &ResponseHandler{
		Service: Service,
	}
}

// Collect godoc
// @Summary Add response.
// @Description Endpoint to submit a new response.
// @Tags response
// @Param formId path string true "Form id"
// @Param userId path string true "User id"
// @Param form body ResponseRequest true "ResponseRequest"
// @Produce json
// @Accept json
// @Success 200 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /response/submit/{formId}/{userId} [post]
func (s *ResponseHandler) CreateResponse(c *gin.Context) {
	var formId = c.Param("formId")
	var userId = c.Param("userId")
	fId, err := strconv.Atoi(formId)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusBadRequest, gin.H{"incorrect Form id": false})
		return
	}
	uId, err := strconv.Atoi(userId)

	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusBadRequest, gin.H{"incorrect User id": false})
		return
	}

	var body entity.Response
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &entity.Response{
		Responses: body.Responses,
	}

	res, err := s.Service.AddResponse(response, int8(fId), int8(uId))

	if res == false {
		c.JSON(http.StatusInternalServerError, gin.H{"Unexpected error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Response submitted.": true})

}
