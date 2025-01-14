import asyncio, os, requests, uuid, json, httpx
from aiogram import Bot, Dispatcher, types
from aiogram.filters import Command
from dotenv import load_dotenv

load_dotenv()
NGINX_URL = "http://127.0.0.1:80"
TELEGRAM_BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")

bot = Bot(token=TELEGRAM_BOT_TOKEN)
dp = Dispatcher()

##############################################################################################################

@dp.message(Command("start"))
async def help_command(message: types.Message):
    help_text = (
        "Список команды:\n"
        "/start - Выводит список команд\n"
        "/login - Проверка статуса пользователя, авторизация\n"
        "/login github - Авторизация через GitHub(coming soon)\n"
        "/login yandex - Авторизация через Яндекс ID(coming soon)\n"
        "/login code - Авторизация через код\n"
        "/disciplines - Список дисциплин\n"
        "/disciplines <ID> - Подробная информация о дисциплине\n"
        "/users - Список пользователей\n"
        "/users <ID> - Подробная информация о пользователе\n"
    )
    await message.answer(help_text)

@dp.message(Command("identification"))
async def start(message: types.Message):
        tokenVHOD = str(uuid.uuid4())
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{NGINX_URL}/login/{message.chat.id}:anon:{tokenVHOD}")
            response.raise_for_status()

@dp.message(Command("login"))
async def login_handler(message: types.Message):
    args = message.text.split()

    if len(args) == 1:
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(f"{NGINX_URL}/identification/{message.chat.id}")
                response.raise_for_status()
                status, tokenVHOD = response.text.split(":")
                status = status.replace("anon", "Анонимный")
                status = status.replace("avtor", "Авторизированный")
                await message.answer(f"Статус - {status[1:]}\nТокен входа - {tokenVHOD[:-1]}")

        except httpx.HTTPStatusError as e:
            if e.response.status_code == 404:
                await message.answer("Пользователь не найден. Идентифицируйтесь командой /identification.")
                await message.answer(
                    "Вы не авторизованы. Пожалуйста, авторизуйтесь, выбрав один из вариантов:\n"
                    "/login github\n"
                    "/login yandex\n"
                    "/login code"
                )
                return

        if status[1:] == "anon" or status[1:] == "Анонимный":
            await message.answer(
                "Вы не авторизованы. Пожалуйста, авторизуйтесь, выбрав один из вариантов:\n"
                "/login github(coming soon)\n"
                "/login yandex(coming soon)\n"
                "/login code"
            )

    elif len(args) == 2:
        auth_type = args[1]

        if auth_type.lower != "github" or auth_type.lower != "yandex" or auth_type.lower != "code":
                    await message.answer("Неверный тип авторизации. Доступные типы: github(coming soon), yandex(coming soon), code")
                    return

        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(f"{NGINX_URL}/identification/{message.chat.id}")
                response.raise_for_status()
                status, tokenVHOD = response.text.split(":")
                status = status.replace("anon", "Анонимный")
                status = status.replace("avtor", "Авторизированный")
                await message.answer(f"Статус - {status[1:]}\n")

        except httpx.HTTPStatusError as e:
            if e.response.status_code == 404:
                await message.answer("Пользователь не найден. Идентифицируйтесь командой /identification.")
                return

        if auth_type.lower() == "github":
            await message.answer("Регистрация через github")
            return
        
        if auth_type.lower() == "yandex":
            await message.answer("Регистрация через yandex")
            return
            
        if auth_type.lower() == "code":
            await message.answer("Регистрация через код")
            return
        return            

@dp.message(Command("disciplines"))
async def list_disciplines(message: types.Message):
    args = message.text.split()

    if len(args) == 1:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{NGINX_URL}/disciplines/")
            response.raise_for_status()
            data = response.json()
            output = ""
            for item in data:
                output += f"ID: {item['id']}\n"
                output += f"Название: {item['name']}\n"
                output += f"Описание: {item['discription']}\n\n"
            await message.answer(output)
            return

    if len(args) == 2:
        if args[1].isdigit():
            async with httpx.AsyncClient() as client:
                response = await client.get(f"{NGINX_URL}/disciplinesNONALL/{args[1]}")
                response.raise_for_status()
                data = response.json()
                output = ""
                output += f"Название: {data['name']}\n"
                output += f"Описание: {data['discription']}\n"
                response2 = await client.get(f"{NGINX_URL}/usersNONALL/{data['prepod_id']}")
                response2.raise_for_status()
                data2 = response2.json()
                output = f"Преподаватель: {data2['username']}\n" + output
                await message.answer(output)
        else:
           await message.answer("Неверный формат команды, используйте /disciplines <id>")
         
@dp.message(Command("users"))
async def list_users(message: types.Message):
    args = message.text.split()
    
    if len(args) == 1:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{NGINX_URL}/users/")
            response.raise_for_status()
            data = response.json()
            output = ""
            if isinstance(data, list):
                for item in data:
                    output += f"ФИО: {item['username']}\n"
                    output += f"ID: {item['id']}\n\n"
            else:
                output = "Неверный формат данных"
        await message.answer(output)
        return

    if len(args) == 2:
        async with httpx.AsyncClient() as client:
            response = await client.get(f"{NGINX_URL}/usersNONALL/{args[1]}")
            response.raise_for_status()
            data = response.json()
            output = ""
            if isinstance(data, dict):
                output += f"ФИО: {data['username']}\n"
                output += f"ID: {data['id']}\n\n"
            else:
                output = "Неверный формат данных"
        await message.answer(output)
        return

##############################################################################################################

@dp.message()
async def NoneComm(message: types.Message):
    await message.answer("Нет такой команды")

async def main():
    await dp.start_polling(bot)

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Остановка')