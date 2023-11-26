package routes

import (
	"github.com/dg222599/go-url-shortener/database"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v2"
)

func ResolveURL(ctx * fiber.Ctx) error{ 
 
	   url:= ctx.Params("url")

	   r1:= database.CreateClient(0)

	   defer r1.Close()

	   value,err := r1.Get(database.Ctx,url).Result()

	   if err== redis.Nil{

		 	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":"Short URL is not present on the database",
			})
	   } else if err!=nil{
		 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":"Can't Connect to the DB",
			})
	   }
	   
	   rateIncrement := database.CreateClient(1)
	   defer rateIncrement.Close()

	   _ =rateIncrement.Incr(database.Ctx,"counter")

	   return ctx.Redirect(value,301)
}