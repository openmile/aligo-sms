from functools import wraps
from sanic.response import json


def authorized():
    def decorator(f):
        @wraps(f)
        async def decorated_function(request, *args, **kwargs):
            token = request.token
            is_authorized = token is not None  # TODO : Authorization Key Setting

            if is_authorized:
                # the user is authorized.
                # run the handler method and return the response
                response = await f(request, *args, **kwargs)
                return response
            else:
                # the user is not authorized.
                return json({'message': 'Forbidden'}, 403)

        return decorated_function

    return decorator
