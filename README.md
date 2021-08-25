# Simple SMS with ALI.Go
알리고 SMS 발송을 처리하는 간단한 웹 서버. 미리 등록해놓은 IP에서만 알리고 API 콜을 할 수 있어서 작은 서버용으로 만듦 	:innocent:

**Go1.17** : Version 1.14 or higher is required

[**ALIGO**](https://smartsms.aligo.in/admin/api/info.html) : 문자전송 서비스

[**Fiber**](https://github.com/gofiber/fiber) : Fiber is an Express inspired web framework built on top of [Fasthttp](https://github.com/valyala/fasthttp).

[**Viper**](https://github.com/spf13/viper) : Viper is a complete configuration solution for Go applications including 12-Factor apps.

## Set ENV variables
Write **.env** file in the working directory


## Set systemd at Linux
vi /etc/systemd/system/sms-server.service

systemctl daemon-reload

systemctl start sms-server.service
