package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	client := []Client{}
  getDB().Where(map[string]interface{}{"broker": "2", "status": 1}).Find(&client)

	for _, c := range client {
		downloadInforme(&c)
	}
}

type Client struct {
	ID      int
	Name    string
	Sinacor string
	Cpfcnpj string
	Broker  int
}

func getDB() *gorm.DB {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func downloadInforme(c *Client) {
	req, err := http.NewRequest("GET", "https://api.xpi.com.br/corporate-tax-ri-apim/v1/download/"+c.Sinacor+"?tradingAccount="+c.Sinacor+"&year=2023&reportType=1&description=0&brand=3&origin=HubAdvisorXP", nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", os.Getenv("XP_BEARER"))
	req.Header.Add("ocp-apim-subscription-key", os.Getenv("XP_SUBSCRIPTION_KEY"))
	req.Header.Add("token", os.Getenv("XP_TOKEN"))
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("bad status: " + res.Status)
	}

	w, err := os.Create(c.Cpfcnpj + "-" + time.Now().Format("2006-01-01") + ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	bodyRes, _ := io.ReadAll(res.Body)

	body := string(bodyRes)

	cleanBody := body[1 : len(body)-1]

	decoded, _ := base64.StdEncoding.DecodeString(cleanBody)

	w.WriteString(string(decoded))
	if err != nil {
		fmt.Println(err)
	}
}
