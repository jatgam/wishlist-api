package main

import (
	"flag"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	"github.com/jatgam/wishlist-api/config"
	"github.com/jatgam/wishlist-api/db"
	"github.com/jatgam/wishlist-api/jwt"
	"github.com/jatgam/wishlist-api/microservice"
	"github.com/jatgam/wishlist-api/models"
	"github.com/jatgam/wishlist-api/routes"
	"github.com/jatgam/wishlist-api/service/sgmail"
	"github.com/jatgam/wishlist-api/validation"
)

const (
	port        = "4000"
	host        = "localhost"
	metricsPort = "4001"
	healthPort  = "4002"
)

func main() {
	var exposePorts bool

	flag.BoolVar(&exposePorts, "expose-ports", false, "Expose Ports outside the docker network")
	flag.Parse()

	serviceConfig := config.GetConfig()
	db := db.Connect(serviceConfig.DB)
	db.AutoMigrate(&models.UserModel{}, &models.ItemModel{})
	defer db.Close()

	router := microservice.NewMicroservice(metricsPort, healthPort, 0, true)
	router.StartHealthRouter()

	// Prometheus metrics can grow if a new metric is created for every url.
	// We need to remove URL parameters from metrics, for things like item ids and usernames
	router.Prometheus.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.String()
		if len(c.Request.URL.Query()) > 0 {
			parts := strings.Split(url, "?")
			url = parts[0]
		}
		for _, p := range c.Params {
			if p.Key == "itemID" {
				url = strings.Replace(url, p.Value, ":itemID", 1)
				break
			}
			if p.Key == "pwResetToken" {
				url = strings.Replace(url, p.Value, ":pwResetToken", 1)
				break
			}
			if p.Key == "rank" {
				url = strings.Replace(url, p.Value, ":rank", 1)
				break
			}
		}
		return url
	}

	ginjwt := jwt.CreateJWTMiddleware(serviceConfig.Secret, serviceConfig.JWTRealmName)

	// Custom Validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passcomplexity", validation.ComplexityValidator)
		v.RegisterValidation("notblank", validation.NotBlank)
	}

	sgmail.SetupMail(serviceConfig.EMail.SendGridAPIKey, serviceConfig.EMail.FromName, serviceConfig.EMail.FromAddress, serviceConfig.EMail.Debug)

	routes.SetupRoutes(&router.RouterGroup, ginjwt)

	if exposePorts {
		router.Run("0.0.0.0" + ":" + port)
	} else {
		router.Run(host + ":" + port)
	}

}
