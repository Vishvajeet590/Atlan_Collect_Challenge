# Atlan_Collect_Challenge

To read about the design part you can refer to this [Noton notebook.](https://jumbled-chime-d61.notion.site/Collect-Plugin-a3b47d6cad224ef59ad841fa4e152c4b) 
 - - - -
Steps to run Atlan Collect

1. Clone this repo and other two repos which are [Atlan_Collect_KGS](https://github.com/Vishvajeet590/Atlan_Collect_KGS) and [Atlan_Collect_Ingest](https://github.com/Vishvajeet590/Atlan_Collect_Ingest)

        git clone https://github.com/Vishvajeet590/Atlan_Collect_Challenge.git
        git clone https://github.com/Vishvajeet590/Atlan_Collect_KGS.git
        git clone https://github.com/Vishvajeet590/Atlan_Collect_Ingest.git

2. Navigate to Atlan_Collect_Challenge directory and run docker-compose.yml file by the command

         docker-compose up

   This will take a few minutes. Ensure 2 postgres DBs and 1 rabbitmq instances are up and runne on their respective ports specified in the docker-compose.yml


3. Navigate to Atlan_Collect_KGS directory and build the image form docker file named as Dockerfile.kgs with this command.

         docker build -t kgs -f Dockerfile.kgs ./
         docker run --rm -p 8080:8080 -e KEY_DATABASE_URL='postgres://vishwajeet:docker@localhost:5431/KeyStore-1?&pool_max_conns=10' --network="host" kgs

   on running this image it will add keys to our Key Db.


4. Now navigate to the Atlan_Collect_Ingest directory build the docker image from file Dockerfile.ingest and run the image

          docker build -t ingest -f Dockerfile.ingest ./
          docker run --rm -p 8081:8080 -e JOB_DATABASE_URL='postgres://vishwajeet:docker@localhost:5432/Hermes?&pool_max_conns=10'  --network="host" ingest
   This is our Ingest Server which will handle the plugin functionality.


5. Last Step, navigate back to the Atlan_Collect_Challenge directory find the Dockerfile.main build it and run it

        docker build -t collect -f Dockerfile.main ./![Screenshot_20220427_224916](https://user-images.githubusercontent.com/42716731/165584675-80cf8d5a-765b-4189-ad2e-b554557851bd.png)

        docker run --rm -p 8080:8080 -e DATABASE_URL='postgres://vishwajeet:docker@localhost:5432/Hermes?&pool_max_conns=10' -e KEY_DATABASE_URL='postgres://vishwajeet:docker@localhost:5431/KeyStore-1?&pool_max_conns=10' --network="host" collect

Now everything is set up go to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to try the end points.
 - - - -
 #### Endpoints####
There are total 5 endpoints 
 - [POST] /form/create           
 - [GET] /form/{formId}
 - [POST] /response/submit/{formId}/{userId}
 - [POST] /response/action
 - [GET] /job/status/{jobId}

* To create post use ```/form/create``` endpoint with json in body. An example form json is shown bellow.
```json
{
  "form_name": "Payments Error",
  "owner_id": 4910,
  "question": [
    {
      "question": "What is your querry ?",
      "question_type": "text/plain"
    },
    {
      "question": "At which date did it occured.",
      "question_type": "text/plain"
    },
    {
      "question": "What was the ammount?",
      "question_type": "text/plain"
    }
  ]
}
```

* Inorder to create response we must know the question Id so we retrive the form by using ```/form/{formId}``` endpoint, and then to create response use ```/response/submit/{formId}/{userId}``` with json in body as shown in example bellow.
```json
{
  "responses": [
    {
      "question_id": 4,
      "response": "My Payment is struck, ammount has been deducted from my account",
      "response_type": "text/plain"
    },
    {
      "question_id": 5,
      "response": "02/04/2022",
      "response_type": "text/plain"
    },
    {
      "question_id": 6,
      "response": "200000",
      "response_type": "text"
    }
  ]
}
```

* Now to perform any action use ```/response/action``` endpoint with json body, again example is bellow
```json
{
  "form_Id": 24,
  "oAuth_code": "4/0AX4XfWinPtL1DQ4CLYNeT5q1dvkDzLYMFJF8FwT0lSPuH-50SSuxm5wtx2Bgwtgsp-by0Q",
  "plugin_code": 1
}
```
you can get OAuth code from here [google sigin link](https://accounts.google.com/o/oauth2/auth/oauthchooseaccount?access_type=offline&client_id=1065888810890-3itaq54a957tapf45okdbrvra37soqop.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive&state=state-token&flowName=GeneralOAuthFlow). As this application not verified by google so sign in is limited to accounts allowed by appliction owner in order to signin use the provided google account id and password or drop an email to add your account

Here the plugin_code is what action is needed to be performed, currently 1 -> Send item to sheets 2 -> mimick sending confirmation messages.

<img src="https://user-images.githubusercontent.com/42716731/165583912-ddcacc5a-2a83-48aa-b607-d2c9b18602ba.png" height="300">
<img src="https://user-images.githubusercontent.com/42716731/165584408-cdf9f1de-3269-4cba-9a99-f3517044d065.png" height="300">
<img src="https://user-images.githubusercontent.com/42716731/165584724-557cabcb-e62f-4d0b-ad5f-ea4250b6b131.png"  height="300">

This is the oAuth code required to add file to sheet according google drive api.
![Screenshot_20220427_224950](https://user-images.githubusercontent.com/42716731/165584875-5b6791d2-d4fd-4a5c-bc04-103f9c6e59a6.png)

* You can check the status of the job by using ```/job/status/{jobId}``` endpoint 

