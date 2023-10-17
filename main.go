package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os" // http-swagger middleware

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

func (s MyNullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

type Product struct {
	ID              int          `json:"id"`
	ProductCode     string       `json:"product_code"`
	ProductName     string       `json:"product_name"`
	ProductType     int          `json:"product_type"`
	SubCategoryCode MyNullString `json:"subcategory_code"`
	CategoryCode    MyNullString `json:"category_code"`
	ActiveStatus    int          `json:"active_status"`
}

var db *sql.DB
var err error

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9211
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://172.16.2.21:8762/eureka/",
		App:                   "ds3-product-v2",
		Port:                  9211,
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

	viper.SetConfigName("dev.yaml")
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

	// db, err = sql.Open("mysql", "dsserver:xxi2121.@tcp(k8s.devel.intra.db.cinema21.co.id:3306)/db_digsig")
	db, err = sql.Open("mysql", DBUsername+":"+DBPassword+"@tcp("+DBHost+":"+DBPort+")/"+DBName)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer db.Close()

	logfile, err := os.Create("ds3-product-v2.log")

	if err != nil {
		log.Println(err.Error())
		log.Fatal(err)
	}

	defer logfile.Close()
	log.SetOutput(logfile)

	router := mux.NewRouter()
	router.HandleFunc("/get", getProduct).Methods("GET")
	router.HandleFunc("/get/{id}", getProductById).Methods("GET")
	router.HandleFunc("/create", insertProduct).Methods("POST")
	router.HandleFunc("/update/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/delete/{id}", deleteProduct).Methods("DELETE")
	http.ListenAndServe(":9211", router)
}

// @Summary Get Product
// @Router  /get [GET]
func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s for %s", r.RemoteAddr, r.URL)
	w.Header().Set("Content-Type", "application/json")
	var products []Product
	result, err := db.Query("SELECT id, product_code, product_name, product_type, subcategory_code, category_code, active_status from ds_product")
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.ProductCode, &product.ProductName, &product.ProductType, &product.SubCategoryCode, &product.CategoryCode, &product.ActiveStatus)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}

// @Summary Get Product By Id
// @Router  /get/{id} [GET]
func getProductById(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s for %s", r.RemoteAddr, r.URL)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT id, product_code, product_name, product_type, subcategory_code, category_code, active_status from ds_product where id = ?", params["id"])
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	defer result.Close()

	var product Product

	for result.Next() {
		err := result.Scan(&product.ID, &product.ProductCode, &product.ProductName, &product.ProductType, &product.SubCategoryCode, &product.CategoryCode, &product.ActiveStatus)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(product)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s for %s", r.RemoteAddr, r.URL)
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO ds_product (product_code, product_name, product_type, subcategory_code, category_code, product_end_date, active_status, created_time, created_by) VALUES(?,?,?,?,?,?,1,NOW(),?)")
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	keyVal := map[string]any{}
	json.Unmarshal(body, &keyVal)
	product_code := keyVal["product_code"]
	product_name := keyVal["product_name"]
	product_type := keyVal["product_type"]
	subcategory_code := keyVal["subcategory_code"]
	category_code := keyVal["category_code"]
	product_end_date := keyVal["product_end_date"]
	created_by := keyVal["created_by"]
	_, err = stmt.Exec(product_code, product_name, product_type, subcategory_code, category_code, product_end_date, created_by)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "New product was created")
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s for %s", r.RemoteAddr, r.URL)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE ds_product SET product_code=?, product_name=?, product_type=?, subcategory_code=?, category_code=?, product_end_date=?, active_status=1, updated_time=NOW(), updated_by=? WHERE id=?;")
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	keyVal := map[string]any{}
	json.Unmarshal(body, &keyVal)
	new_product_code := keyVal["product_code"]
	new_product_name := keyVal["product_name"]
	new_product_type := keyVal["product_type"]
	new_subcategory_code := keyVal["subcategory_code"]
	new_category_code := keyVal["category_code"]
	new_product_end_date := keyVal["product_end_date"]
	updated_by := keyVal["updated_by"]
	_, err = stmt.Exec(new_product_code, new_product_name, new_product_type, new_subcategory_code, new_category_code, new_product_end_date, updated_by, params["id"])
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product with ID = %s was updated", params["id"])
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM ds_product WHERE id = ?")
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product with ID = %s was deleted", params["id"])
}
