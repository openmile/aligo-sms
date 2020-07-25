import aiohttp

ALIGO_URL = 'https://apis.aligo.in'


async def aligo_sms(data):
    async with aiohttp.ClientSession() as session:
        aligo_suffix = '/send/'
        async with session.post(ALIGO_URL + aligo_suffix, data=data, ssl=False) as response:
            aligo_body = await response.json(content_type=None)
            return aligo_body
