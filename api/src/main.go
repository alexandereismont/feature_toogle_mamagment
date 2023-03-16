package main

import (
	"net/http"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/Unleash/unleash-client-go/v3/context"
	"github.com/gin-gonic/gin"
)

type BuyProduct struct {
	User    string `json:"user"`
	Country string `json:"country"`
	Product string `json:"product"`
	Sale    string `json:"sale"`
}

var boughtProduct = []BuyProduct{
	{User: "Alexander", Country: "Norway", Product: "T-shirt"},
}

func buyWithContext(c *gin.Context) {
	var newBuyProduct BuyProduct

	if err := c.BindJSON(&newBuyProduct); err != nil {
		return
	}

	contextMap := make(map[string]string)

	contextMap["country"] = newBuyProduct.Country
	ctx := context.Context{
		Properties: contextMap,
	}

	if unleash.IsEnabled("context_toogle", unleash.WithContext(ctx)) {
		newBuyProduct.Sale = newBuyProduct.Country + "-20%"
	}

	c.IndentedJSON(http.StatusCreated, newBuyProduct)
}

func buyWithCanery(c *gin.Context) {
	var newBuyProduct BuyProduct

	if err := c.BindJSON(&newBuyProduct); err != nil {
		return
	}

	if unleash.IsEnabled("canary") {
		c.IndentedJSON(http.StatusCreated, newBuyProduct)
	} else {
		c.IndentedJSON(http.StatusAccepted, "Not supported")
	}

}

func main() {
	initUnleash()
	router := gin.Default()
	router.POST("/buy/withcontext", buyWithContext)
	router.POST("/buy/withcanary", buyWithCanery)

	router.Run("localhost:8000")
}

func initUnleash() {
	unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName("my-application"),
		unleash.WithEnvironment("development"),
		unleash.WithUrl("http://localhost:4242/api"),
		unleash.WithCustomHeaders(http.Header{"Authorization": {"*:development.242ef0ddfd634e6b424a3b4b0e25dfc827c84781ae442147b38724eb"}}),
	)
}
