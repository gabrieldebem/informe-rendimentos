package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	client := Client{
		Broker: 2,
	}
	getDB().First(&client, 1)
  fmt.Println(client.Name)

	downloadInforme(&client)
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
	req, err := http.NewRequest("GET", "https://api.xpi.com.br/corporate-tax-ri-apim/v1/download/"+"5682826"+"?tradingAccount="+"5682826"+"&year=2023&reportType=1&description=0&brand=3&origin=HubAdvisorXP", nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", os.Getenv("XP_BEARER"))
	req.Header.Add("ocp-apim-subscription-key", os.Getenv("XP_SUBSCRIPTION_KEY"))
	req.Header.Add("token", os.Getenv("XP_TOKEN"))

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	w, err := os.Create(c.Cpfcnpj + ".pdf")

	if err != nil {
		fmt.Println(err)
	}

	defer w.Close()

	io.Copy(w, res.Body)

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
