package jobHandler

import (
	"Atlan_Collect_Challenge/entity"
	"Atlan_Collect_Challenge/usecase/jobStatus"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strconv"
)

type jobId struct {
	formId int `json:"job_id"`
}

type userId struct {
	userId int `json:"user_id"`
}

type message struct {
	Message string `json:"message"`
}

type job struct {
	JobId         int    `json:"job_id"`
	JobStatus     string `json:"job_status"`
	JobStatusCode int    `json:"job_status_code"`
	PluginCode    int    `json:"job_plugin_code"`
}

type JobRequestSwag struct {
	FormId     int    `json:"form_Id"`
	OAuthCode  string `json:"OAuth_code"`
	PluginCode int    `json:"plugin_code"`
}

type JobHandler struct {
	Service *jobStatus.Service
	Ch      *amqp.Channel
}

var TAG = "FormHandler"

func NewJobHandler(Service *jobStatus.Service, channel *amqp.Channel) *JobHandler {
	return &JobHandler{
		Service: Service,
		Ch:      channel,
	}
}

// Collect godoc
// @Summary Take action on all response.
// @Description Endpoint to perform action accordingly on all response of a form.
// @Tags response
// @Produce json
// @Accept json
// @Param form body JobRequestSwag true "Action Request"
// @Success 200 {object} message
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /response/action [post]
func (s *JobHandler) CreateJob(c *gin.Context) {

	var body entity.JobRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jobReq := &entity.JobRequest{
		FormId:     body.FormId,
		OAuthCode:  body.OAuthCode,
		PluginCode: body.PluginCode,
	}

	jobid, err := s.Service.AddJob(jobReq.PluginCode)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusBadRequest, gin.H{"Unexpected error : ": err.Error()})
		return
	}

	err = s.SendJob(jobReq.FormId, jobid, jobReq.PluginCode, jobReq.OAuthCode)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusInternalServerError, gin.H{"Unexpected error : ": err.Error()})
		return
	}

	resultStr := fmt.Sprintf("Job created. Check the progress with job id = %d", jobid)
	c.JSON(http.StatusOK, gin.H{resultStr: true})

}

// Collect godoc
// @Summary Get status of the action.
// @Description Endpoint to get the status of the action.
// @Schemes
// @Tags response
// @Produce json
// @Accept json
// @Param jobId path string true "Job id"
// @Success 200 {object} job
// @Failure 400 {object} message
// @Failure 500 {object} message
// @Router /job/status/{jobId} [get]
func (s *JobHandler) GetJobStatus(c *gin.Context) {
	jobId := c.Param("jobId")
	id, err := strconv.Atoi(jobId)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusBadRequest, gin.H{"Bad jobId error : ": err.Error()})
		return
	}

	job, err := s.Service.GetJob(id)
	if err != nil {
		log.Printf("%s : ")
		c.JSON(http.StatusInternalServerError, gin.H{"Unexpected error : ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"job": job})

}

//Helper Function to send msg to rabbit mq

func (s *JobHandler) SendJob(formId, jobId, pluginCode int, oAuthCode string) error {
	/*
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println("Failed Initializing Broker Connection")
			panic(err)
		}

		ch, err := conn.Channel()
		if err != nil {
			fmt.Println(err)
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"PluginQueue",
			false,
			false,
			false,
			false,
			nil,
		)
		// We can print out the status of our Queue here
		// this will information like the amount of messages on
		// the queue
		fmt.Println(q)
		// Handle any errors if we were unable to create the queue
		if err != nil {
			fmt.Println(err)
		}*/

	// attempt to publish a message to the queue!

	//msg format "formid,jobId,code,pluginCode"

	//key := fmt.Sprintf("%d",jobId)
	messageStr := fmt.Sprintf("%d,%d,%s,%d", formId, jobId, oAuthCode, pluginCode)
	log.Printf("\n%s\n", messageStr)
	err := s.Ch.Publish(
		"",
		"PluginQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageStr),
		},
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully Published Message to Queue")
	return nil
}
