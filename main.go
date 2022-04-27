package main

import (
	"Atlan_Collect_Challenge/api/handlers/formHandler"
	jobHandler2 "Atlan_Collect_Challenge/api/handlers/jobHandler"
	"Atlan_Collect_Challenge/api/handlers/responseHandler"
	"Atlan_Collect_Challenge/repository"
	"Atlan_Collect_Challenge/usecase/form"
	"Atlan_Collect_Challenge/usecase/jobStatus"
	"Atlan_Collect_Challenge/usecase/response"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/streadway/amqp"
	_ "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/swaggerFiles"
	"os"

	swaggerFiles "github.com/swaggo/files"
	_ "net/http"

	_ "Atlan_Collect_Challenge/docs"
)

// @title API Atlan Collect
// @version version(1.0)
// @description Atlan collect with plugin feature

// @contact.name API supporter
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name license(Mandatory)
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {

	//env
	//KEY_DATABASE_URL=postgres://vishwajeet:vishvapriya123@localhost:5432/KeyStore-1?&pool_max_conns=10;
	//DATABASE_URL=postgres://vishwajeet:vishvapriya123@localhost:5432/Hermes?&pool_max_conns=10
	//postgres://vishwajeet:docker@host.docker.internal:5432/KeyStore-1?&pool_max_conns=10
	//postgres://vishwajeet:docker@host.docker.internal:5432/Hermes?&pool_max_conns=10

	DATABASE_URL := os.Getenv("DATABASE_URL")
	fmt.Printf(" DB = %s\n", DATABASE_URL)
	ctx := context.Background()
	cofig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, cofig)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}
	defer pool.Close()

	KEY_DATABASE_URL := os.Getenv("KEY_DATABASE_URL")
	key_cofig, err := pgxpool.ParseConfig(KEY_DATABASE_URL)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}

	key_pool, err := pgxpool.ConnectConfig(context.Background(), key_cofig)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}
	defer key_pool.Close()

	//RabbitMq
	/*amqp://guest:guest@localhost:5672/*/
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
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

	fmt.Println(q)
	if err != nil {
		fmt.Println(err)
	}

	formRepo := repository.NewFormDbSql(pool, key_pool)
	formService := form.NewService(formRepo)
	formHandler := formHandler.NewFormHandler(formService)

	responseRepo := repository.NewResponseDbSql(pool, key_pool)
	responseService := response.NewService(responseRepo)
	responseHandler := responseHandler.NewResponseHandler(responseService)

	jobRepo := repository.NewJobDbSql(pool)
	jobService := jobStatus.NewService(jobRepo)
	jobHandler := jobHandler2.NewJobHandler(jobService, ch)

	r := gin.Default()
	r.POST("/form/create", formHandler.CreateForm)
	r.GET("/form/:formId", formHandler.GetForm)
	r.POST("/response/submit/:formId/:userId", responseHandler.CreateResponse)
	r.POST("/response/action", jobHandler.CreateJob)
	r.GET("/job/status/:jobId", jobHandler.GetJobStatus)

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run()

}
