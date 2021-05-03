import aiohttp
import asyncio

ALIGO_URL = "https://apis.aligo.in"


async def aligo_message(session: aiohttp.ClientSession, aligo_data, url_suffix=None):
    aligo_suffix = url_suffix or "/send/"
    async with session.post(
        ALIGO_URL + aligo_suffix,
        data=aligo_data,
        ssl=False,
    ) as response:
        return await response.json(content_type=None)
