from sanic import Sanic
from sanic.response import json, text
from dotenv import load_dotenv
import os
from app.aligo import aligo_sms
from app.auth import authorized
import asyncio

load_dotenv()

ALIGO_URL = 'https://apis.aligo.in'
ALIGO_ID = os.getenv('ALIGO_ID')
ALIGO_KEY = os.getenv('ALIGO_KEY')
ALIGO_PHONE = os.getenv('ALIGO_PHONE')

app = Sanic(__name__)


@app.get('/')
@authorized()
async def main(request):
    ip = request.headers.get('X-Real-IP')
    return text(ip)


@app.route('/send/sms', methods=['POST'])
async def send_sms(request):
    message, status_code = {'message': ''}, 200
    body = request.json
    if body.get('a_key') != ALIGO_KEY:
        return json({'message': '알리고 키 인증에 실패했습니다.'}, 403)
    else:
        # bytes_test(body.get('msg'))
        aligo_data = {
            'key': ALIGO_KEY,
            'user_id': ALIGO_ID,
            'sender': ALIGO_PHONE,
            'receiver': body.get('receiver'),
            'msg': body.get('msg'),
            'msg_type': 'SMS',
        }
        try:
            # TODO: requests will be deprecated
            '''
            aligo_suffix = '/send/'
            response = requests.post(ALIGO_URL + aligo_suffix, data=aligo_data)
            aligo_body = response.json()
            '''
            # TODO : Search Use and Exception Cases
            # loop = app.loop
            aligo_task = asyncio.create_task(aligo_sms(aligo_data))  # == (app.)loop.create_task()
            aligo_body = await aligo_task
            # try aligo_body = await asyncio.wait_for(aligo_sms(aligo_data), timeout=3.0) except asyncio.TimeoutError
            # try aligo_body = await asyncio.wait_for(aligo_task, timeout=3.0) except asyncio.TimeoutError

            if int(aligo_body.get('result_code')) < 0:
                print('알리고에서 문자전송에 실패한 케이스입니다.')
                message = {
                    'message': aligo_body.get('message')
                }
                message.update({'status': 'fail'})  # TODO: this will be deprecated
                status_code = 400
            else:
                message = {
                    'message': '문자전송에 성공하였습니다.',
                    'data': [
                        {'msg_id': aligo_body.get('msg_id')}
                    ]
                }
                message.update({'status': 'success'})  # TODO: this will be deprecated
        except Exception as error:
            message = {
                'message': '문자전송 서비스에서 에러가 발생했습니다. 잠시 후에 다시 시도해주세요.',
                'data': [
                    {'error_message': f'{str(error)}'}
                ]
            }
            message.update({'status': 'fail'})  # TODO: this will be deprecated
            status_code = 500
        finally:
            return json(message, status_code)
