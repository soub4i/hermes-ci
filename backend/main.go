package main

import (
	"Heremes-ci-server/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/hibiken/asynq"

)


var (
	JobsCollection *mongo.Collection
	MI             MongoInstance
)
const (

	TypeJobProcessing = "job:process"
)
type JobPayload struct {
    ID     string
    Body 	Payload
}

var client *asynq.Client

func JobProcessingTask(id string, body Payload) (*asynq.Task, error) {
    payload, err := json.Marshal(JobPayload{ID: id, Body: body})
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(TypeJobProcessing, payload), nil
}

func Setup() {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")

	MI = MongoInstance{
		Client: client,
		DB:     client.Database(os.Getenv("DB")),
	}

	JobsCollection = MI.DB.Collection("jobs")
}

func main() {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	client = asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS")})
    defer client.Close()

	Setup()
	router := mux.NewRouter()

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*",
		},
	})

	// Set routing rules
	router.HandleFunc("/github/{id}", handleGitHubWebhook).Methods("POST")
	router.HandleFunc("/github/{id}", GetJobs).Methods("GET")
	router.HandleFunc("/jobs/{id}", GetJob).Methods("GET")

	//Use the default DefaultServeMux.
	err := http.ListenAndServe(":"+os.Getenv("PORT"), corsOpts.Handler(router))
	if err != nil {
		log.Fatal(err)
	}
}

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}



func CreateJob(b models.Job) (string, error) {
	// TODO: clean this up by expecting a context from the caller to better propagate cancelation
	ctx := context.Background()
	result, err := JobsCollection.InsertOne(ctx, b)
	if err != nil {
		return "", err
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)

	return fmt.Sprintf("%v", oid.Hex()), nil
}

type Payload struct {
	Ref        string `json:"ref"`
	Url        string `json:"html_url"`
	Repository struct {
		Name      string `json:"name"`
		FullName  string `json:"full_name"`
		Timestamp string `json:"timestamp"`
		Owner     struct {
			Name   string `json:"name"`
			Login  string `json:"login"`
			Avatar string `json:"avatar_url"`
		} `json:"owner"`
	} `json:"repository"`
	Commit struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	} `json:"head_commit"`
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	// TODO: clean this up by expecting a context from the caller to better propagate cancelation
	ctx := context.Background()

	vars := mux.Vars(r)
	id := vars["id"]

	var jobs []models.Job
	var job models.Job
	cursor, err := JobsCollection.Find(ctx, bson.M{"repository_id": id})

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&job)
		if err != nil {
			log.Fatal(err)
			return
		}
		jobs = append(jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	// TODO: clean this up by expecting a context from the caller to better propagate cancelation
	ctx := context.Background()
	vars := mux.Vars(r)
	id := vars["id"]

	objectId, err := primitive.ObjectIDFromHex(id)

	var job models.Job
	err = JobsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&job)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func handleGitHubWebhook(_ http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}

	decoder := json.NewDecoder(r.Body)
	var body Payload
	err := decoder.Decode(&body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(body)

	// Insert job into database

	job := models.Job{
		Name:         body.Commit.Message,
		Repository:   body.Repository.FullName,
		CreateAt:     time.Now(),
		Finished:     false,
		RepositoryId: id,
		Owner: models.OwnerType{
			Name:   body.Repository.Owner.Name,
			Login:  body.Repository.Owner.Login,
			Avatar: body.Repository.Owner.Avatar,
		},
	}

	idd, err := CreateJob(job)
	if err != nil {
		log.Println(err)
	}

	task, err :=  JobProcessingTask(idd, body)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
