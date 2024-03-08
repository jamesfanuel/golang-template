package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os" // http-swagger middleware
	"strings"
	"time"

	"crypto/md5"
	"encoding/hex"

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

type Schedule struct {
	PlayDate    MyNullString `json:"play_date"`
	TheaterCode MyNullString `json:"theater_code"`
	City        MyNullString `json:"city"`
	TheaterName MyNullString `json:"theater_name"`
	Studio      MyNullString `json:"studio"`
	MovieCode   MyNullString `json:"movie_code"`
	MovieName   MyNullString `json:"movie_name"`
	ShowType    MyNullString `json:"show_type"`
	ShowTime    MyNullString `json:"show_time"`
}

type FinalSchedule struct {
	ScheduleDate string             `json:"schedule_date"`
	TheaterCode  string             `json:"theater_code"`
	Studio       string             `json:"studio"`
	Schedule     []ScheduleByOutlet `json:"schedule"`
}

type ScheduleByOutlet struct {
	ShowID    MyNullString `json:"show_id"`
	ShowTime  MyNullString `json:"show_time"`
	MovieCode MyNullString `json:"movie_code"`
	ShowType  MyNullString `json:"show_type"`
}

type ScheduleMD5 struct {
	MD5 string `json:"md5"`
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
		DefaultZone: "http://localhost:8763/eureka/",
		// DefaultZone:           "http://172.16.2.21:8763/eureka/",
		App:                   "x1-schedule",
		Port:                  6001,
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
	fmt.Printf("viper : %s = %s \n", "Service Port", "6001")

	// db, err = sql.Open("mysql", "dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_digsig")
	db, err = sql.Open("mysql", DBUsername+":"+DBPassword+"@tcp("+DBHost+":"+DBPort+")/"+DBName)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer db.Close()

	logfile, err := os.Create("logs/x1-schedule.log")

	if err != nil {
		log.Println(err.Error())
		log.Fatal(err)
	}

	defer logfile.Close()
	log.SetOutput(logfile)

	router := mux.NewRouter()
	router.HandleFunc("/TrafficIntegration", getSchedule).Methods("GET")
	// router.HandleFunc("/get/{id}", getMemberLogById).Methods("GET")
	router.HandleFunc("/create", insertSchedule).Methods("POST")
	http.ListenAndServe(":6001", router)
}

// @Summary Get Product
// @Router  /get [GET]
func getSchedule(w http.ResponseWriter, r *http.Request) {
	var redisHost = "localhost:6379"
	var redisPassword = ""

	rdb := newRedisClient(redisHost, redisPassword)
	fmt.Println("redis client initialized")

	currentTime := time.Now()
	// get data
	op2 := rdb.Get(context.Background(), currentTime.Format("2006-01-02")+"_"+strings.Split(r.URL.Query().Get("t"), "_")[0]+"_"+strings.Split(r.URL.Query().Get("t"), "_")[1])
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

// func GetMD5Hash(text string) string {
// 	hash := md5.Sum([]byte(text))
// 	return hex.EncodeToString(hash[:])
// }

// func OnPage(link string) string {
// 	res, err := http.Get(link)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	content, err := io.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(content)
// }

func insertSchedule(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	log.Printf("Download File")
	downloadFile("file/jadwal.txt", "https://m.cinemaxxi.net/ftp-jadwal/xl"+currentTime.Format("060102")+".txt")

	log.Printf("Request from %s for %s", r.RemoteAddr, r.URL)
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open("file/jadwal.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash1 := md5.New()
	_, err = io.Copy(hash1, file)

	if err != nil {
		panic(err)
	}

	// hash2 := md5.New()
	// hash2.Write([]byte(OnPage("https://m.cinemaxxi.net/ftp-jadwal/xl" + currentTime.Format("060102") + ".txt")))

	result, err := db.Query("SELECT value from setting_tbl where varname = 'schedule.md5'")
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

	if hex.EncodeToString(hash1.Sum(nil)) != ScheduleMD5.MD5 {

		file, err := os.Open("file/jadwal.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		del, err := db.Prepare("DELETE FROM temp_schedule")
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		_, err = del.Exec()

		stmt, err := db.Prepare("INSERT INTO temp_schedule (play_date, theater_code, city, theater_name, studio, movie_code, movie_name, show_type, show_time) VALUES(?,?,?,?,?,?,?,?,?)")
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		for scanner.Scan() {
			theater_code := scanner.Text()[0:7]
			city := scanner.Text()[7:32]
			theater_name := scanner.Text()[32:57]
			studio := scanner.Text()[57:59]
			movie_code := scanner.Text()[59:65]
			movie_name := scanner.Text()[65:95]
			show_type := scanner.Text()[95:97]
			show_time := scanner.Text()[97:len([]rune(scanner.Text()))]
			fmt.Println("Insert " + theater_code + "_" + studio + " " + movie_name)
			start_index := 0
			end_index := 5
			for i := 0; i < len([]rune(strings.TrimSpace(show_time)))/5; i++ {
				_, err = stmt.Exec(currentTime.Format("2006-01-02"), theater_code, city, theater_name, studio, movie_code, movie_name, show_type, show_time[start_index:end_index])
				start_index += 5
				end_index += 5
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		log.Printf("Start Insert Into Schedule")

		del2, err := db.Prepare("DELETE FROM schedule")
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		_, err = del2.Exec()

		trn, err := db.Prepare(`INSERT INTO schedule(show_date, cinema_code, screen_no, show_start, show_type, film_row_id, show_no)
		SELECT play_date
		, theater_code
		, studio
		, ADDTIME(show_time , "00:09") AS show_time
		, show_type
		, movie_code
		, ROW_NUMBER() OVER(PARTITION BY play_date, theater_code, studio ORDER BY play_date, theater_code, studio, show_time) AS show_no
		FROM temp_schedule
		ORDER BY play_date, theater_code , studio , show_time`)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		_, err = trn.Exec()

		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		log.Printf("Done Insert Into Schedule")

		var redisHost = "localhost:6379"
		var redisPassword = ""

		rdb := newRedisClient(redisHost, redisPassword)
		fmt.Println("redis client initialized")

		currentTime = time.Now()

		result, err := db.Query("SELECT show_date, cinema_code, screen_no FROM schedule WHERE show_date = ? GROUP BY show_date, cinema_code, screen_no", currentTime.Format("2006-01-02"))
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
		defer result.Close()

		for result.Next() {
			var FinalSchedule FinalSchedule
			var ScheduleByOutlets []ScheduleByOutlet

			err := result.Scan(&FinalSchedule.ScheduleDate, &FinalSchedule.TheaterCode, &FinalSchedule.Studio)
			if err != nil {
				log.Println(err.Error())
				panic(err.Error())
			}

			result2, err := db.Query(`SELECT
			CONCAT(replace(show_date,"-",""),".",cinema_code,".",screen_no,".",show_no) AS show_id
			, CONCAT (show_date, " ", show_start ) AS show_time
			, film_row_id AS movie_code
			, show_type
			FROM
			schedule WHERE cinema_code = ? AND screen_no = ? AND show_date = ? ORDER BY show_start`, &FinalSchedule.TheaterCode, &FinalSchedule.Studio, currentTime.Format("2006-01-02"))

			if err != nil {
				log.Println(err.Error())
				panic(err.Error())
			}

			defer result2.Close()
			for result2.Next() {
				var ScheduleByOutlet ScheduleByOutlet
				err := result2.Scan(&ScheduleByOutlet.ShowID, &ScheduleByOutlet.ShowTime, &ScheduleByOutlet.MovieCode, &ScheduleByOutlet.ShowType)
				if err != nil {
					log.Println(err.Error())
					panic(err.Error())
				}

				ScheduleByOutlets = append(ScheduleByOutlets, ScheduleByOutlet)
			}

			FinalSchedule.Schedule = append(FinalSchedule.Schedule, ScheduleByOutlets...)

			key := currentTime.Format("2006-01-02") + "_" + FinalSchedule.TheaterCode + "_" + FinalSchedule.Studio
			data, err := json.Marshal(FinalSchedule)
			ttl := time.Duration(64800) * time.Second

			// store data using SET command
			op1 := rdb.Set(context.Background(), key, data, ttl)
			if err := op1.Err(); err != nil {
				fmt.Printf("unable to SET data. error: %v", err)
				return
			}

			fmt.Println("Insert Key " + currentTime.Format("2006-01-02") + "_" + FinalSchedule.TheaterCode + "_" + FinalSchedule.Studio + " to Redis")
			log.Println("set operation success")
		}

		json.NewEncoder(w).Encode("Schedule Success Insert to Redis")

		stmt2, err := db.Prepare("UPDATE setting_tbl SET value = ? WHERE varname='schedule.md5'")
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
		_, err = stmt2.Exec(hex.EncodeToString(hash1.Sum(nil)))
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

	}

	json.NewEncoder(w).Encode("No Data to Update")

}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
