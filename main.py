from sanic import Sanic
from sanic.response import json
import requests
from dotenv import load_dotenv
import os

load_dotenv()

ALIGO_URL = 'https://apis.aligo.in'
ALIGO_ID = os.getenv('ALIGO_ID')
ALIGO_KEY = os.getenv('ALIGO_KEY')
ALIGO_PHONE = os.getenv('ALIGO_PHONE')

app = Sanic(__name__)


@app.get('/')
async def main(request):
    return json(
        {
            "ip": request.ip,
         }
    )


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
            aligo_suffix = '/send/'
            response = requests.post(ALIGO_URL + aligo_suffix, data=aligo_data)
            aligo_body = response.json()
            if int(aligo_body.get('result_code')) < 0:
                print('알리고에서 문자전송에 실패한 케이습니다.')
                message = {
                    'message': aligo_body.get('message')
                }
                status_code = 400
            else:
                message = {
                    'message': '문자전송에 성공하였습니다.',
                    'data': [
                        {'msg_id': aligo_body.get('msg_id')}
                    ]
                }
        except Exception as error:
            message = {
                'message': '알리고 서비스에서 에러가 발생했습니다. 잠시 후에 다시 시도해주세요.',
                'data': [
                    {'error_message': f'{str(error)}'}
                ]
            }
            status_code = 500
        finally:
            return json(message, status_code)
