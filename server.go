package main

import (
	"github.com/fuellab/simple-sms/aligo"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

func readEnv(key string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading env file: %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}

func main() {
	aligoKey := readEnv("ALIGO_KEY")

	app := fiber.New(fiber.Config{
		Prefork:      true,
		AppName:      "FUELLAB SMS v2.0.0",
		ServerHeader: "fuellab simple sms",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "fail",
		})
	})

	app.Post("/send/sms", func(c *fiber.Ctx) error {
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

		theCode, _ := strconv.Atoi(aligoRes.ResultCode)
		if theCode < 0 {
			return c.Status(400).JSON(fiber.Map{
				"message": aligoRes.Message,
				"status":  "fail",
			})
		}

		detail := make([]map[string]interface{}, 1)
		detail[0]["msg_id"] = aligoRes.MsgId

		return c.JSON(fiber.Map{
			"message": "문자전송에 성공하였습니다.",
			"status":  "success",
			"data":    detail,
		})
	})

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
