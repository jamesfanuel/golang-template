package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os" // http-swagger middleware
	"os/exec"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/cnf/structhash"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	eureka "github.com/xuanbo/eureka-client"
)

func getEnv() string {
	return os.Getenv("APP_ENV")
}

type MyNullString struct {
	sql.NullString
}

func newRedisClient(host string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	return client
}

func (s MyNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

type Response struct {
	DataIntis []DataInti
}

type DataInti struct {
	MovieKey                string `json:"movieKey"`
	MovieCode               string `json:"movieCode"`
	MovieMTIX               string `json:"movieMTIX"`
	movieDistributor        string
	MoviePresentationTypeId int `json:"moviePresentationTypeId"`
	moviePresentationType   string
	movieCategoryID         int
	movieCategory           string
	movieDirector           string
	movieScriptWriter       string
	movieProducer           string
	movieSinopsis           string
	movieStars              string
	movieOfficialWebsite    string
	MovieGenre              string `json:"movieGenre"`
	MovieRating             string `json:"movieRating"`
	MovieDuration           string `json:"movieDuration"`
	movieTags               string
	movieShortTitle         string
	movieFeatures           string
	MovieTitle              string `json:"movieTitle"`
	movieSinopsisEng        string
	movieStatusId           int
	MovieStatus             string `json:"movieStatus"`
	movieCreateTime         string
	movieUpdateTime         string
}

type Slots struct {
	SlotNames    []string `json:"slot_names"`
	SlotOrder    []int    `json:"slot_order"`
	SlotAutoplay []int    `json:"slot_autoplay"`
	SlotRepeat   []int    `json:"slot_repeat"`
}

type FinalMovie struct {
	SlotNames    []string    `json:"slot_names"`
	SlotOrder    []int       `json:"slot_order"`
	SlotAutoplay []int       `json:"slot_autoplay"`
	SlotRepeat   []int       `json:"slot_repeat"`
	MovieInfo    []MovieInfo `json:"movie_info"`
}

type MovieInfo struct {
	MovieCode string `json:"movie_code"`
	Title     string `json:"title"`
	LsfRating string `json:"lsf_rating"`
	Duration  string `json:"duration"`
	Genre     string `json:"genre"`
	Show3D    int    `json:"show3d"`
	Status    int    `json:"status"`
}

type ScheduleMD5 struct {
	MD5 string `json:"md5"`
}

func (movieInfo *MovieInfo) SetMovieCode(movie_code string) {
	movieInfo.MovieCode = movie_code
}

func (movieInfo *MovieInfo) SetTitle(title string) {
	movieInfo.Title = title
}

func (movieInfo *MovieInfo) SetRating(rating string) {
	movieInfo.LsfRating = rating
}

func (movieInfo *MovieInfo) SetDuration(duration string) {
	movieInfo.Duration = duration
}

func (movieInfo *MovieInfo) SetGenre(genre string) {
	movieInfo.Genre = genre
}

func (movieInfo *MovieInfo) Set3d(show3d int) {
	movieInfo.Show3D = show3d
}

func (movieInfo *MovieInfo) SetStatus(status int) {
	movieInfo.Status = status
}

var db *sql.DB
var err error

var ctx = context.Background()

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:6001
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	client := eureka.NewClient(&eureka.Config{
		// DefaultZone: "http://localhost:8763/eureka/",
		DefaultZone:           "http://172.16.2.21:8763/eureka/",
		App:                   "x1-movie",
		Port:                  6011,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})

	// start client, register、heartbeat、refresh
	client.Start()

	viper.SetConfigName("prod.yaml")
	viper.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("DB.HOST", "k8s.devel.intra.db.cinema21.co.id")

	DBHost, ok := viper.Get("DB.HOST").(string)
	DBPort, ok := viper.Get("DB.PORT").(string)
	DBUsername, ok := viper.Get("DB.USERNAME").(string)
	DBPassword, ok := viper.Get("DB.PASSWORD").(string)
	DBName, ok := viper.Get("DB.NAME").(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	fmt.Printf("viper : %s = %s \n", "Database Host", DBHost)
	fmt.Printf("viper : %s = %s \n", "Database Port", DBPort)
	fmt.Printf("viper : %s = %s \n", "Database Username", DBUsername)
	fmt.Printf("viper : %s = %s \n", "Database Password", DBPassword)
	fmt.Printf("viper : %s = %s \n", "Database Name", DBName)
	fmt.Printf("viper : %s = %s \n", "Service Port", "6011")

	// db, err = sql.Open("mysql", "dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_digsig")
	db, err = sql.Open("mysql", DBUsername+":"+DBPassword+"@tcp("+DBHost+":"+DBPort+")/"+DBName)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer db.Close()

	logfile, err := os.Create("logs/x1-movie.log")

	if err != nil {
		log.Println(err.Error())
		log.Fatal(err)
	}

	defer logfile.Close()
	log.SetOutput(logfile)

	router := mux.NewRouter()
	router.HandleFunc("/TrafficIntegration", getMovie).Methods("GET")
	router.HandleFunc("/TrafficIntegration/md5", getMD5).Methods("GET")
	// router.HandleFunc("/get/{id}", getMemberLogById).Methods("GET")
	router.HandleFunc("/save", saveToRedis).Methods("POST")
	http.ListenAndServe(":6011", router)
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	// out, err := os.Create(filepath)
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Get the data
	// os.Setenv("HTTP_PROXY", "http://idproxy.cinema21.co.id:9908")
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// Check server response
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("bad status: %s", resp.Status)
	// }

	// Writer the body to file
	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func saveToRedis(w http.ResponseWriter, r *http.Request) {

	var redisHost = "localhost:6379"
	// var redisHost = "172.16.1.138:6379"
	var redisPassword = ""

	rdb := newRedisClient(redisHost, redisPassword)
	// fmt.Println("redis client initialized")

	currentTime := time.Now().AddDate(0, -6, 0)

	/// feed movie
	downloadFile("file/movie.json", "https://manager.cinema21.co.id/DIApi/feed/movies/getAll/"+currentTime.Format("2006-01-02"))
	// file, err := os.Open("file/movie.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	body, err := ioutil.ReadFile("file/movie.json")
	if err != nil {
		log.Fatal(err)
	}

	// res, err := http.Get("https://manager.cinema21.co.id/DIApi/feed/movies/getAll/" + currentTime.Format("2006-01-02"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	/// slot
	res2, err2 := http.Get("http://172.16.2.21:6010/api/v1/getSlotInfo")
	if err2 != nil {
		log.Fatal(err2)
	}
	defer res2.Body.Close()

	body2, err2 := ioutil.ReadAll(res2.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	/// initiate variable
	var finalMovie FinalMovie

	var dataInti []DataInti
	json.Unmarshal([]byte(body), &dataInti)

	var slots Slots
	json.Unmarshal([]byte(body2), &slots)

	var dtMovieInfo []MovieInfo

	for i, _ := range dataInti {
		dt := MovieInfo{}
		dt.SetMovieCode(dataInti[i].MovieCode)
		dt.SetTitle(dataInti[i].MovieTitle)
		switch dataInti[i].MovieRating {
		case "S":
			dt.SetRating("0")
		case "R":
			dt.SetRating("1")
		case "D":
			dt.SetRating("2")
		}
		dt.SetDuration(dataInti[i].MovieDuration)
		dt.SetGenre(dataInti[i].MovieGenre)
		switch {
		case dataInti[i].MoviePresentationTypeId >= 60:
			dt.Set3d(1)
		case dataInti[i].MoviePresentationTypeId < 60:
			dt.Set3d(0)
		}
		dt.SetStatus(1)
		dtMovieInfo = append(dtMovieInfo, dt)
	}

	finalMovie.MovieInfo = dtMovieInfo
	finalMovie.SlotNames = slots.SlotNames
	finalMovie.SlotAutoplay = slots.SlotAutoplay
	finalMovie.SlotOrder = slots.SlotOrder
	finalMovie.SlotRepeat = slots.SlotRepeat

	result, err := db.Query("SELECT value from setting_tbl where varname = 'movie.md5'")
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer result.Close()

	var ScheduleMD5 ScheduleMD5

	for result.Next() {
		err := result.Scan(&ScheduleMD5.MD5)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
	}

	variable := fmt.Sprintf("%s", structhash.Md5(finalMovie, 1))

	if variable != ScheduleMD5.MD5 {
		// store data using SET command
		key := "movie"
		data, _ := json.Marshal(finalMovie)
		ttl := time.Duration(604800) * time.Second

		op1 := rdb.Set(context.Background(), key, data, ttl)
		if err := op1.Err(); err != nil {
			fmt.Printf("unable to SET data. error: %v", err)
		}

		op2 := rdb.Get(context.Background(), "movie")
		if err := op2.Err(); err != nil {
			fmt.Printf("unable to GET data. error: %v", err)
			return
		}
		redisMovie, err := op2.Result()
		if err != nil {
			fmt.Printf("unable to GET data. error: %v", err)
			return
		}
		log.Println("get operation success. result:", redisMovie)

		hash := md5.Sum([]byte(redisMovie))

		stmt2, err := db.Prepare("UPDATE setting_tbl SET value = ? WHERE varname='movie.md5'")
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
		_, err = stmt2.Exec(hex.EncodeToString(hash[:]))
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		json.NewEncoder(w).Encode("Movie Success Insert to Redis")
	} else {
		json.NewEncoder(w).Encode("No Data to Update ")
	}
}

// @Summary Get Product
// @Router  /get [GET]
func getMovie(w http.ResponseWriter, r *http.Request) {
	var redisHost = "localhost:6379"
	// var redisHost = "172.16.1.138:6379"
	var redisPassword = ""

	rdb := newRedisClient(redisHost, redisPassword)

	// get data
	op2 := rdb.Get(context.Background(), "movie")
	if err := op2.Err(); err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return
	}
	res, err := op2.Result()
	if err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return
	}
	log.Println("get operation success. result:", res)

	fmt.Fprintf(w, res)
}

// @Summary Get Product
// @Router  /get [GET]
func getMD5(w http.ResponseWriter, r *http.Request) {
	var redisHost = "localhost:6379"
	// var redisHost = "172.16.1.138:6379"
	var redisPassword = ""

	rdb := newRedisClient(redisHost, redisPassword)
	// fmt.Println("redis client initialized")

	// get data
	op2 := rdb.Get(context.Background(), "movie")
	if err := op2.Err(); err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return
	}
	res, err := op2.Result()
	if err != nil {
		fmt.Printf("unable to GET data. error: %v", err)
		return
	}
	log.Println("get operation success. result:", res)

	hash := md5.Sum([]byte(res))

	fmt.Fprintf(w, hex.EncodeToString(hash[:]))
}
