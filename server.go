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
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusForbidden

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = ctx.Status(code).JSON(fiber.Map{
				"message": "Forbidden", "status": "fail",
			})
			return nil
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden", "status": "fail",
		})
	})

	app.Post("/send/sms", func(c *fiber.Ctx) error {
		sendData := new(aligo.SendData)
		if err := c.BodyParser(sendData); err != nil {
			log.Print(err)
			return c.Status(403).JSON(fiber.Map{
				"message": "Invalid body", "status": "fail",
			})
		}
		if sendData.AKey != aligoKey {
			return c.Status(403).JSON(fiber.Map{
				"message": "알리고 키 인증에 실패했습니다.", "status": "fail",
			})
		}

		sendData.Key = aligoKey
		sendData.UserId = readEnv("ALIGO_ID")
		sendData.Sender = readEnv("ALIGO_PHONE")

		aligoRes := aligo.PostAligo(sendData)

		if aligoRes.ResultCode != "1" {
			return c.Status(400).JSON(fiber.Map{
				"message": aligoRes.Message,
				"status":  "fail",
			})
		}

		detail := make(map[string]string)
		detail["msg_id"] = aligoRes.MsgId

		return c.JSON(fiber.Map{
			"message": "문자전송에 성공하였습니다.", "status": "success", "data": detail,
		})
	})

	err := app.Listen(":9001")
	if err != nil {
		return
	}
}
