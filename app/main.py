import os
import asyncio
from asyncio import TimeoutError
import aiohttp
from dotenv import load_dotenv
from sanic import Sanic
from sanic.response import json, text
from app.aligo import aligo_message
from app.auth import authorized

load_dotenv()

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
    try:
        if body["a_key"] != ALIGO_KEY:
            raise KeyError

        aligo_data = {
            "key": ALIGO_KEY,
            "user_id": ALIGO_ID,
            "sender": ALIGO_PHONE,
            "receiver": body.get("receiver"),
            "msg": body.get("msg"),
            "msg_type": body.get("msg_type") or "SMS",
            "title": body.get("title"),
        }
        async with aiohttp.ClientSession() as session:
            aligo_task = asyncio.create_task(
                aligo_message(
                    session,
                    aligo_data,
                )
            )
            aligo_body = await asyncio.wait_for(aligo_task, timeout=10.0)

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
                "data": [
                    {
                        "msg_id": aligo_body.get("msg_id"),
                    }
                ],
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
    except TimeoutError:
        return json(
            {
                "message": "문자전송 서비스에서 처리시간이 오래 걸리고 있습니다. 서비스를 확인해주세요.",
                "status": "fail",
            },
            409,
        )
    except (Exception, TypeError) as error:
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
