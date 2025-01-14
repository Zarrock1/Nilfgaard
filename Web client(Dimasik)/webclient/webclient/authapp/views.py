from django.http import HttpResponseRedirect
from django.shortcuts import render, redirect
from django.http import JsonResponse
from django.core.cache import cache
import uuid
import requests
import redis


redis_host = '127.0.0.1' 
redis_port = 6379  


redis_client = redis.StrictRedis(host=redis_host, port=redis_port, decode_responses=True)


session_token = 'your_session_token'  
jwt_token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzYzNTAzMjIsImlhdCI6MTczNjI2MzkyMiwiaXNzIjoiQXNzaXJldmFyQW5haGlkIiwidXNlcmFjY2VzcyI6WyJ1c2VyOmxpc3Q6cmVhZCIsInVzZXI6ZnVsbE5hbWU6d3JpdGUiLCJ1c2VyOmJsb2NrOnJlYWQiXSwidXNlcmxvZ2luIjoiSm9obkRvZSJ9.0Ii31MqpXhpVTnKGiGz0gPR_eBSa6XDAalropi0enGk'

# Сохраняем токен в Redis
redis_client.set(session_token, jwt_token, ex=3600)  

print(f"Токен сохранён в Redis с ключом {session_token}")


def home(request):
    return JsonResponse({'message': 'Это главная страница'}, status=200)

def login(request):
    type_param = request.GET.get('type')
    if not type_param:
        return render(request, 'authapp/login.html')

    session_token = str(uuid.uuid4())
    
    user_id = 'user7' 
    api_url = f"http://78.136.201.177:8080/api/token?userId={user_id}"
    
    
    try:
        response = requests.get(api_url)
        if response.status_code == 200:
            access_token = response.content.decode('utf-8')  
            cache.set(session_token, {'access_token': access_token}, timeout=3600)
            resp = JsonResponse({'message': 'Авторизация успешна', 'access_token': access_token})
            resp.set_cookie('session_token', session_token, httponly=True)
            return resp
        else:
            return JsonResponse({'message': 'Ошибка авторизации с API'}, status=401)
    except requests.RequestException as e:
        return JsonResponse({'error': str(e)}, status=500)


def logout(request):
    session_token = request.COOKIES.get('session_token')
    if not session_token:
        return redirect('/login')

    
    all_param = request.GET.get('all')
    if all_param == 'true':
        cache.delete(session_token)  
        requests.post('http://127.0.0.1/logout', data={'refresh_token': refresh_token})
        response = redirect('/login')
        response.delete_cookie('session_token')  
        return response
    else:
        # Локальный выход
        cache.delete(session_token) 
        response = redirect('/')
        response.delete_cookie('session_token')  
        return response
    
def authorize_user(login_token):
    return {'status': 'waiting', 'message': f'Токен ключ {(login_token)}'}


def disciplins_view(request):
      disciplins = 'http://91.245.227.51:3000/testapi/disciplins'
      return HttpResponseRedirect(disciplins) 
def users_views(request):
    users = 'http://91.245.227.51:3000/testapi/users'
    return HttpResponseRedirect(users)
def token(request):
    token = 'http://78.136.201.177:8080/api/'
    return HttpResponseRedirect(token)






               