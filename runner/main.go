package main

import (
	"log"
	"github.com/hibiken/asynq"
	"encoding/json"
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/joho/godotenv"
	"os/exec"
	"strings"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"bytes"
)
const (

	TypeJobProcessing = "job:process"
)

type JobPayload struct {
    ID     string
    Body 	Payload
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


type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type Job struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Workflow     string             `json:"workflow" bson:"workflow"`
	Repository   string             `json:"repository" bson:"repository"`
	RepositoryId string             `json:"repository_id" bson:"repository_id"`
	Logs         []string           `json:"logs,omitempty" bson:"logs,omitempty"`
	Finished     bool               `json:"finished" bson:"finished"`
	Owner        OwnerType          `bson:"owner,omitempty" json:"owner,omitempty"`
	CreateAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type OwnerType struct {
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Login  string `bson:"login,omitempty" json:"login,omitempty"`
	Avatar string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}


var (
	JobsCollection *mongo.Collection
	MI             MongoInstance
)

func parseYAMLFile(path string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
/*Setup opens a database connection to mongodb*/
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

func UpdateJob(id string, logs []string, workflow string) error {
	// TODO: clean this up by expecting a context from the caller to better propagate cancelation
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)

	update := bson.M{"$set": bson.M{"logs": logs, "finished": true, "updated_at": time.Now(), "workflow": workflow}}

	_, err = JobsCollection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	return err
}

type Log struct {
    Timestamp int64 `json:"timestamp"`
    Data string `json:"data"`
}



func HandleJobProcessing(ctx context.Context, t *asynq.Task) error {
    var p JobPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }

	body := p.Body

	path := "./tmp/" + body.Repository.Name 

	exec.Command("rm", "-rf", path).Run()

	// clone repo and checkout to commit

	log.Println("Cloning " + body.Repository.FullName)

	Url := "https://github.com/" + body.Repository.FullName

	branch := strings.Split(body.Ref, "/")[2]

	if branch == "" {
		branch = "master"
	}

	cmd := exec.Command("git", "clone", "-b", branch, Url, path)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(stdout))

	HermesConfigFilePath := path + "/Hermes.yaml"
	// HermesConfigFilePath := "../test//Hermes.yaml"

	exist := fileExists(HermesConfigFilePath)

	// TODO: Move this part to runner app
	// TODO: Add queuing
	if exist == false {
		return fmt.Errorf("Hermes.yaml not found")
	}

	HermesFile, err := parseYAMLFile(HermesConfigFilePath)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if HermesFile["schema"] == nil {
		return fmt.Errorf("field schema in Hermes.yaml not found")
	}
	if HermesFile["schema"] == "docker" {

		ctx := context.Background()
		cmd = exec.CommandContext(ctx, "buildah", "bud", "--log-level", "info", "-t", strings.ToLower(body.Repository.FullName), path)
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			fmt.Println(err)
		}
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
		}
		var buf []string
		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			o := Log{Timestamp: time.Now().Unix(), Data: string(bytes.Trim(tmp, "\x00")) }
   			 om, _ := json.Marshal(o)
			line := fmt.Sprintf("%s", string(om) )
			buf = append(buf, line)
			if err != nil {
				break
			}
		}

		// Cleanup
		imageName := "localhost/" + strings.ToLower(body.Repository.FullName)

		exec.Command("buildah", "rmi", imageName).Run()
		exec.Command("rm", "-rf", path).Run()
		workflow := HermesFile["name"].(string)

		if err := UpdateJob(p.ID, buf, workflow); err != nil {
			log.Fatalf("could not update job: %v", err)
		}

	
	}

    return nil
}

func main() {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	redisConnection := asynq.RedisClientOpt{
		Addr: os.Getenv("REDIS"), // Redis server address
	}

	Setup()

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: 4,
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})

	mux := asynq.NewServeMux()

    mux.HandleFunc(TypeJobProcessing, HandleJobProcessing)
    // ...register other handlers...

    if err := worker.Run(mux); err != nil {
        log.Fatalf("could not run server: %v", err)
    }
}