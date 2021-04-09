from sanic import Sanic
from sanic.response import json, text
from dotenv import load_dotenv
import os
from app.aligo import aligo_sms
from app.auth import authorized
import asyncio

load_dotenv()

ALIGO_URL = "https://apis.aligo.in"
ALIGO_ID = os.getenv("ALIGO_ID")
ALIGO_KEY = os.getenv("ALIGO_KEY")
ALIGO_PHONE = os.getenv("ALIGO_PHONE")

app = Sanic(__name__)


@app.get("/")
@authorized()
async def main(request):
    ip = request.headers.get("X-Real-IP")
    return text(ip)


@app.route("/send/sms", methods=["POST"])
async def send_sms(request):
    body = request.json
    if body.get("a_key") != ALIGO_KEY:
        raise KeyError
    else:
        aligo_data = {
            "key": ALIGO_KEY,
            "user_id": ALIGO_ID,
            "sender": ALIGO_PHONE,
            "receiver": body.get("receiver"),
            "msg": body.get("msg"),
            "msg_type": body.get("msg_type") or "SMS",
            "title": body.get("title"),
        }
        try:
            aligo_task = asyncio.create_task(aligo_sms(aligo_data))
            aligo_body = await aligo_task

            # loop = app.loop
            # try:
            #     aligo_body = await asyncio.wait_for(aligo_task, timeout=3.0)
            # except asyncio.TimeoutError:
            #     print("time out")

            if int(aligo_body.get("result_code")) < 0:
                # 알리고에서 문자전송에 실패한 케이스입니다.
                message = {
                    "message": aligo_body.get("message"),
                    "status": "fail",
                }
                status_code = 400
            else:
                message = {
                    "message": "문자전송에 성공하였습니다.",
                    "data": [{"msg_id": aligo_body.get("msg_id")}],
                    "status": "success",
                }
                status_code = 200
        except KeyError:
            return json(
                {
                    "message": "알리고 키 인증에 실패했습니다.",
                    "status": "fail",
                },
                403,
            )
        except Exception as error:
            return json(
                {
                    "message": "문자전송 서비스에서 에러가 발생했습니다. 잠시 후에 다시 시도해주세요.",
                    "data": {"error_message": f"{str(error)}"},
                    "status": "fail",
                },
                500,
            )
        else:
            return json(message, status_code)
