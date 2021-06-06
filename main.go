package main

import (
	"fmt"
	"os"
	"encoding/json"
	"path/filepath"
	"log"
	"tars_dht22/models"
	"tars_dht22/configer"
	"time"
	"github.com/d2r2/go-dht"
)

//config reader
func Config_reader(cfg_file string) configer.Configuration {

	file, err := os.Open(cfg_file)
	if err != nil {
		fmt.Println("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := configer.Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("can't decode config JSON: ", err)
	}

	return Config
}

func inserter(temp_val float32,hum_val float32) {
        //log.Println("# Inserting values")
        dt := time.Now()
        var lastInsertId int
        err := models.Db.QueryRow("INSERT INTO rasp_weather (temp_val,hum_val,r_date) VALUES($1,$2,$3) returning r_id;", temp_val,hum_val, dt).Scan(&lastInsertId)
        checkErr(err)
        fmt.Println("last inserted id =", lastInsertId)

}


func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }


func main() {
	version := "0.0.1"
	fmt.Println("tars dht22 inserter  version:"+version)
//************************* read config ******************************************//
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

//log.Println(dir)
	cfg := Config_reader(filepath.Join(dir,"db.conf"))
	models.Initdb(cfg)

//*********************** parse config **********************************//
   //logging
	log_dir := "./log"
	if _, err := os.Stat(log_dir); os.IsNotExist(err) {
		os.Mkdir(log_dir, 0644)
	}
	file, err := os.OpenFile(filepath.Join(log_dir,cfg.Log_file_name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("Logging to a file !")

	// Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
	// You may enable "boost GPIO performance" parameter, if your device is old
	// as Raspberry PI 1 (this will require root privileges). You can switch off
	// "boost GPIO performance" parameter for old devices, but it may increase
	// retry attempts. Play with this parameter.
	// Note: "boost GPIO performance" parameter is not work anymore from some
	// specific Go release. Never put true value here.
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 4, false, 10)
	if err != nil {
		log.Fatal(err)
	}
        inserter(temperature, humidity)
	// Print temperature and humidity
	//fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
	//	temperature, humidity, retried)
}
