package routes

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/dg222599/go-url-shortener/database"
	"github.com/dg222599/go-url-shortener/helpers"
)

type request struct{
	URL string 				`json:"url"`
	CustomShort string		`json:"short_url"`
	Expiry time.Duration	`json:"expiry"`
}

type response struct{
	URL string 				`json:"url"`
	CustomShort string		`json:"short_url"`
	Expiry time.Duration	`json:"expiry"`
	XRateRemaining int 		`json:"rate_limit"`
	XRateLimitReset time.Duration 	`json:"rate_limit_reset"`
}

func ShortenURL(ctx *fiber.Ctx) error{
	 body := &request{}

	 

	 if err := ctx.BodyParser(&body); err!=nil{
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Can't parse the JSON",
		})
	 }

	 fmt.Println(body)

	 r2  := database.CreateClient(1)
	 defer r2.Close()

	 value,err := r2.Get(database.Ctx,ctx.IP()).Result()
	 limit,_ := r2.TTL(database.Ctx,ctx.IP()).Result()

	 if err == redis.Nil{
		 _  = r2.Set(database.Ctx,ctx.IP(),os.Getenv("API_QOUTA"),30*60*time.Second).Err()
	 } else if err == nil{
		valInt,_ := strconv.Atoi(value)
		if valInt <= 0{
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"Rate Limit Exceeded",
				"rate_limit_resent":limit/time.Nanosecond/time.Minute,
			})
		}
	 }

	 if !govalidator.IsURL(body.URL){
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid URL",
		})
	 }

	 if !helpers.RemoveDomainError(body.URL){
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Please provide other URL than this Browser address",
		})
	 }

	 body.URL = helpers.EnforceHTTP(body.URL)

	 var id string
	 
	 if body.CustomShort == ""{
		 id = helpers.Base62Encode(rand.Uint64())
	 } else {
		id = body.CustomShort
	 }

	 r:=database.CreateClient(0)

	 defer r.Close()

	 val , _ := r.Get(database.Ctx,id).Result()

	 if val != ""{
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":"Custome URL is already in use",
		})
	 }

	 if body.Expiry == 0{
		body.Expiry = 24
	 }

	 err = r.Set(database.Ctx,id,body.URL,body.Expiry*3600*time.Second).Err()

	 fmt.Println(err)

	 if err!=nil{
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "I am the error",
		}) 
	 }

	 defaultAPIQouta := os.Getenv("API_QOUTA")
	 if err !=nil{
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Unable to connect to teh server",
		})
	 }

	 defaultAPIQoutaInt,_ := strconv.Atoi(defaultAPIQouta)

	 resp:=response{
		URL: body.URL,
		CustomShort:"",
		Expiry : body.Expiry,
		XRateRemaining: defaultAPIQoutaInt,
		XRateLimitReset: 30,
	 }

	 remainingQouta,_ :=r2.Decr(database.Ctx,ctx.IP()).Result()

	 resp.XRateRemaining = int(remainingQouta)

	 resp.XRateRemaining = int(limit / time.Nanosecond / time.Minute)

	 resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	 return ctx.Status(fiber.StatusOK).JSON(resp)



}