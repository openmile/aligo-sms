package main

import (
	"github.com/fuellab/simple-sms/aligo"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func readEnv(key string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Read env file error %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}

func main() {
	aligoKey := readEnv("ALIGO_KEY")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "fail",
		})
	})

	app.Get("/send/sms", func(c *fiber.Ctx) error {
		sendData := new(aligo.SendData)
		if err := c.BodyParser(sendData); err != nil {
			log.Print(err)
			return c.Status(403).JSON(fiber.Map{
				"message": "Invalid body",
				"status":  "fail",
			})
		}

		if aligoKey != sendData.Key {
			return c.Status(403).JSON(fiber.Map{
				"message": "알리고 키 인증에 실패했습니다.",
				"status":  "fail",
			})
		}

		aligoRes := aligo.PostAligo(sendData)

		if aligoRes.ResultCode < 0 {
			return c.Status(400).JSON(fiber.Map{
				"message": aligoRes.Message,
				"status": "fail",
			})
		}
		return c.JSON(aligoRes)
	})

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
