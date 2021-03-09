# simple-sms with ALIGO
알리고 SMS 전송을 다루는 간단한 웹 서비스. 미리 등록해놓은 IP에서만 알리고 API 콜을 할 수 있어서 작은 서버용으로 만든 것.

Python >= 3.8

[ALIGO](https://smartsms.aligo.in/admin/api/info.html)
문자전송 서비스

[Sanic Framework](https://sanic.readthedocs.io/)
Web Framework that allows the usage of the async/await syntax, faster than flask

[AIOHTTP](https://docs.aiohttp.org/en/stable/)
Asynchronous HTTP Client/Server for asyncio and Python. In this project, only use Client

[Uvicorn](https://www.uvicorn.org/)
ASGI server, built on uvloop and httptools

## Run service
Write your env config(aligo key, id, phone) on .env file or export in bash

- python -m venv my-venv/project-name
- source my-venv/project-name/bin/activate
- pip install -r requirements.txt
- uvicorn app.main:app

## Extra setting (systemd)
- vi /etc/systemd/system/project-name.socket & project-name.service
- systemctl daemon-reload
- systemctl enable socket & service
- systemctl start socket & service
