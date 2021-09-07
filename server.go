package main

import (
	"github.com/fuellab/simple-sms/aligo"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func readENV(key string) string {
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

var aligoKey, aligoID, aligoPhone string

func init() {
	aligoKey = readENV("ALIGO_KEY")
	aligoID = readENV("ALIGO_ID")
	aligoPhone = readENV("ALIGO_PHONE")
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		AppName:      "FUELLAB SMS V2.0.0",
		ServerHeader: "FUELLAB Simple SMS",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusNotFound
			body := "Not Found"

			if e, ok := err.(*fiber.Error); ok && e.Code == fiber.StatusInternalServerError {
				code = e.Code
				body = e.Message
			}

			err = c.Status(code).JSON(fiber.Map{"message": body, "status": "fail"})
			return nil
		},
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
		sendData.UserId = aligoID
		sendData.Sender = aligoPhone

		aligoRes := aligo.PostAligo(sendData)

		if aligoRes.ResultCode != "1" {
			return c.Status(400).JSON(fiber.Map{
				"message": aligoRes.Message, "status": "fail"})
		}

		detail := make(map[string]string)
		detail["msg_id"] = aligoRes.MsgId

		return c.JSON(fiber.Map{
			"message": "문자전송에 성공하였습니다.", "status": "success", "data": detail,
		})
	})

	app.Get("/*", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	err := app.Listen(":9001")
	if err != nil {
		return
	}
}
