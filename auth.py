from functools import wraps
from sanic.response import json


def authorized():
    def decorator(f):
        @wraps(f)
        async def decorated_function(request, *args, **kwargs):
            white_list = ['127.0.0.1']
            ip = request.ip
            is_authorized = ip in white_list

            if is_authorized:
                # the user is authorized.
                # run the handler method and return the response
                response = await f(request, *args, **kwargs)
                return response
            else:
                # the user is not authorized.
                return json({'message': 'Not Authorized'}, 403)

        return decorated_function

    return decorator
