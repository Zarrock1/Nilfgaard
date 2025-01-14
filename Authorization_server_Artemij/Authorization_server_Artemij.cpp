#include <iostream>
#include <fmt/core.h>
#include <string>
#include <vector>
#include <mongocxx/client.hpp>
#include <mongocxx/instance.hpp>
#include <mongocxx/uri.hpp>
#include <bsoncxx/builder/stream/document.hpp>
#include <bsoncxx/json.hpp>
#include <bsoncxx/builder/basic/document.hpp>
#include <bsoncxx/builder/basic/kvp.hpp>
#include <jwt-cpp/jwt.h>
#include <httplib.h>
#include <memory>
#include <fstream>
#include <chrono>
#include <random>
#include <sstream>




// AuthorizationServer отвечает за авторизацию пользователей и управление правами доступа.
class AuthorizationServer {
private:
    mongocxx::client client; // Подключение к MongoDB.
    mongocxx::database db = client["zarrock"];// База данных MongoDB.
    mongocxx::collection users; // Коллекция пользователей в базе данных.
    mongocxx::collection logs; // Коллекция для логирования событий.
    mongocxx::collection users_collection = db["users_collection"];
    std::string secret_key = "NilfgaardAssirevarAnahid"; // Секретный ключ для JWT.
    httplib::Server svr; // HTTP-сервер для обработки API-запросов.

    std::string generateAccessToken(const std::string& userId) //генерирует токен доступа (Access Token) для пользователя с указанным userId. Токен подписывается с использованием алгоритма HS256 и секретного ключа
    {
        auto token = jwt::create()
            .set_type("JWT")
            .set_algorithm("HS256")
            .set_payload_claim("userId", jwt::claim(userId))
            .set_issuer("AuthorizationServer")
            .set_expires_at(std::chrono::system_clock::now() + std::chrono::minutes(30))
            .set_issued_at(std::chrono::system_clock::now())
            .sign(jwt::algorithm::hs256{ secret_key });
        return token;
    }

    std::string generateRefreshToken() //генерирует токен обновления (Refresh Token), который представляет собой случайную строку длиной 64 символа, состоящую из цифр и букв
    {
        static const char characters[] =
            "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
        std::random_device rd;
        std::mt19937 generator(rd());
        std::uniform_int_distribution<> distribution(0, sizeof(characters) - 2);
        std::stringstream refreshToken;
        for (int i = 0; i < 64; ++i) {
            refreshToken << characters[distribution(generator)];
        }
        return refreshToken.str();
    }

    // Метод для записи лога в базу данных.
    void logEvent(const std::string& event, const std::string& details) {
        bsoncxx::builder::stream::document log_builder;
        log_builder << "event" << event
            << "details" << details
            << "timestamp" << bsoncxx::types::b_date(std::chrono::system_clock::now());
        logs.insert_one(log_builder.view());

        // Добавляем вывод события в консоль
        std::cout << "Log Event: " << event << ", Details: " << details << std::endl;
    }

public:
    AuthorizationServer(const std::string& uri)
        : client(mongocxx::uri{ uri }), db(client["authorization_server"]), users(db["users"]), logs(db["logs"]) {

        // API-эндпоинт для генерации токена.
        svr.Get("/api/token", [this](const httplib::Request& req, httplib::Response& res) {
            std::string userId = req.get_param_value("userId");
            if (userId.empty()) {
                res.status = 400;
                res.set_content("User ID is required.", "text/plain");
                return;
            }
            std::string token = generateTokenAndSaveUser(userId);
            logEvent("TokenGenerated", fmt::format("Token generated for userId: {}", userId));
            res.set_content(fmt::format(token), "application/json");
            });

        // API-эндпоинт для получения прав доступа пользователя.
        svr.Get("/api/permissions", [this](const httplib::Request& req, httplib::Response& res) {
            std::string userId = req.get_param_value("userId");
            if (userId.empty()) {
                res.status = 400;
                res.set_content("User ID is required.", "text/plain");
                return;
            }
            std::string permission = getPermission(userId);
            logEvent("GetPermissions", fmt::format("Permissions retrieved for userId: {}", userId));
            res.set_content(permission, "text/plain");
            });

        // API-эндпоинт для установки прав доступа пользователя.
        svr.Post("/api/permissions", [this](const httplib::Request& req, httplib::Response& res) {
            auto userId = req.get_param_value("userId");
            auto permission = req.get_param_value("permission");
            if (userId.empty() || permission.empty()) {
                res.status = 400;
                res.set_content("User ID and permission are required.", "text/plain");
                return;
            }
            setPermission(userId, permission);
            logEvent("SetPermissions", fmt::format("Permission '{}' set for userId: {}", permission, userId));
            res.status = 200;
            res.set_content("Permission set successfully.", "text/plain");
            });

        // API-эндпоинт для получения информации о пользователе.
        svr.Get("/api/userinfo", [this](const httplib::Request& req, httplib::Response& res) {
            std::string userId = req.get_param_value("userId");
            if (userId.empty()) {
                res.status = 400;
                res.set_content("User ID is required.", "text/plain");
                return;
            }
            try {
                std::string userInfo = getUserInfo(userId);
                logEvent("GetUserInfo", fmt::format("User info retrieved for userId: {}", userId));
                res.set_content(userInfo, "application/json");
            }
            catch (const std::exception& e) {
                res.status = 404;
                res.set_content(e.what(), "text/plain");
            }
            });

        // API-эндпоинт для обновления имени пользователя.
        svr.Post("/api/updatename", [this](const httplib::Request& req, httplib::Response& res) {
            std::string userId = req.get_param_value("userId");
            std::string newName = req.get_param_value("newName");
            if (userId.empty() || newName.empty()) {
                res.status = 400;
                res.set_content("User ID and new name are required.", "text/plain");
                return;
            }
            updateUserName(userId, newName);
            logEvent("UpdateUserName", fmt::format("User name updated for userId: {}, new name: {}", userId, newName));
            res.status = 200;
            res.set_content("User name updated successfully.", "text/plain");
            });

        // API-эндпоинт для проверки подключения к внешним сервисам.
        svr.Get("/api/check_connection", [this](const httplib::Request& req, httplib::Response& res) {
            bool connectedToGithub = checkConnection("https://api.github.com");
            bool connectedToYandex = checkConnection("https://yandex.ru");

            std::string status = fmt::format("Connected to GitHub: {}\nConnected to Yandex: {}",
                connectedToGithub ? "Yes" : "No", connectedToYandex ? "Yes" : "No");
            logEvent("CheckConnection", "External services connection status checked.");
            res.set_content(status, "text/plain");
            });

        // Генерация ссылки для авторизации через код
        svr.Get("/api/generate_code", [this](const httplib::Request& req, httplib::Response& res) {
            std::string login_token = req.get_param_value("token");
            if (login_token.empty()) {
                res.status = 400;
                res.set_content("400", "text/plain");
                return;
            }

            // Генерация случайного 5-6 значного кода
            std::hash<std::string> hash_function;
            size_t hashed_value = hash_function(login_token + std::to_string(std::time(nullptr)));
            std::string code = std::to_string(hashed_value % 900000 + 100000);

            // Время истечения
            std::time_t current_time = std::time(nullptr);
            std::time_t expiration_time = current_time + 60; // Код устаревает через 1 минуту

            // Запись в базу данных
            bsoncxx::builder::stream::document filter_builder, update_builder;
            filter_builder << "login_token" << login_token;
            update_builder << "$set" << bsoncxx::builder::stream::open_document
                << "auth_code" << bsoncxx::builder::stream::open_document
                << "code" << code
                << "expiration_time" << static_cast<int64_t>(expiration_time)
                << bsoncxx::builder::stream::close_document
                << bsoncxx::builder::stream::close_document;

            try {
                auto result = users_collection.update_one(filter_builder.view(), update_builder.view(), mongocxx::options::update{}.upsert(true));
                if (!result || result->modified_count() == 0) {
                    res.status = 500;
                    res.set_content("Error: Failed to update user record", "text/plain");
                    return;
                }
            }
            catch (const std::exception& e) {
                res.status = 500;
                res.set_content(std::string("500") + e.what(), "text/plain");//Internal server error
                return;
            }

            // Ответ с кодом
            res.set_content(code, "text/plain");
            });

        // Валидация кода и генерация нового токена
        svr.Post("/api/validate_code", [this](const httplib::Request& req, httplib::Response& res) {
            std::string login_token = req.get_param_value("token");
            std::string code = req.get_param_value("code");

            if (login_token.empty() || code.empty()) {
                res.status = 400;
                res.set_content("Error: Missing token or code", "text/plain");
                return;
            }

            // Поиск записи в базе данных
            bsoncxx::builder::stream::document filter_builder;
            filter_builder << "login_token" << login_token
                << "auth_code.code" << code;

            try {
                auto record = users_collection.find_one(filter_builder.view());
                if (!record) {
                    res.status = 401;
                    res.set_content("401", "text/plain");//Invalid code or token

                    return;
                }

                auto auth_code = record->view()["auth_code"].get_document().view();
                std::time_t expiration_time = auth_code["expiration_time"].get_int64();
                std::time_t current_time = std::time(nullptr);

                // Проверка на истечение срока действия
                if (current_time > expiration_time) {
                    res.status = 401;
                    res.set_content("401", "text/plain");// Code expired
                    return;
                }

                // Генерация нового токена
                std::string new_token = "new_generated_token"; // Реальная генерация токена по вашему методу

                // Обновление записи в базе данных
                bsoncxx::builder::stream::document update_builder;
                update_builder << "$set" << bsoncxx::builder::stream::open_document
                    << "token" << new_token
                    << bsoncxx::builder::stream::close_document;

                users_collection.update_one(filter_builder.view(), update_builder.view());

                // Возврат нового токена
                res.set_content(new_token, "text/plain");
            }
            catch (const std::exception& e) {
                res.status = 500;
                res.set_content(std::string("Error: ") + e.what(), "text/plain");
            }
            });



        // Генерация ссылки для авторизации через Яндекс
        svr.Get("/api/auth_links/yandex", [this](const httplib::Request& req, httplib::Response& res) {
            const std::string yandex_client_id = "20412cf318f7479194f92d08d43d099e";
            const std::string redirect_uri = "http://78.136.201.177:8080/callback";//


            std::string yandex_auth_url = fmt::format(
                "https://oauth.yandex.ru/authorize?response_type=code&client_id={}&redirect_uri={}",
                yandex_client_id, redirect_uri);

            std::string response_body = fmt::format(
                "Yandex Auth Link: {}",
                 yandex_auth_url);

            logEvent("GenerateAuthLinks", "Generated Yandex auth link.");
            res.set_content(response_body, "text/plain");
            });

        // Генерация ссылки для авторизации через github

        svr.Get("/api/auth_links/github", [this](const httplib::Request& req, httplib::Response& res) {
            const std::string github_client_id = "Ov23liBuZYozkvpMXGMq";
            const std::string redirect_uri = "http://78.136.201.177:8080/callback";//

            std::string github_auth_url = fmt::format(
                "https://github.com/login/oauth/authorize?client_id={}&redirect_uri={}",
                github_client_id, redirect_uri);

            

            std::string response_body = fmt::format(
                "GitHub Auth Link: {}",
                github_auth_url);

            logEvent("GenerateAuthLinks", "Generated GitHub auth link.");
            res.set_content(response_body, "text/plain");
            });
    }

    // Генерация токена для пользователя и сохранение его в базе данных.
    std::string generateTokenAndSaveUser(const std::string& userId) {
        bsoncxx::builder::stream::document filter_builder;
        filter_builder << "userId" << userId;
        auto result = users.find_one(filter_builder.view());

        std::string username = userId;//вот тут поменял
        std::vector<std::string> useraccess = { "user:list:read", "user:fullName:write", "user:block:read" };

        if (!result) {
            bsoncxx::builder::stream::document document_builder;
            document_builder << "userId" << userId
                << "permission" << "read_write"
                << "username" << username;
            users.insert_one(document_builder.view());
            std::cout << "New user created with userId: " << userId << std::endl;
        }
        else {
            username = std::string(result->view()["username"].get_string().value);
        }

        auto token = jwt::create()
            .set_type("JWT")
            .set_algorithm("HS256")
            .set_payload_claim("userlogin", jwt::claim(username))
            .set_payload_claim("useraccess", jwt::claim(jwt::traits::kazuho_picojson::array_type(useraccess.begin(), useraccess.end())))
            .set_issuer("AssirevarAnahid")
            .set_expires_at(std::chrono::system_clock::now() + std::chrono::hours(24))
            .set_issued_at(std::chrono::system_clock::now())
            .sign(jwt::algorithm::hs256{ secret_key });

        // Сохраняем информацию из токена в базу данных.
        bsoncxx::builder::stream::document token_info;
        token_info << "userId" << userId
            << "token" << token
            << "issued_at" << bsoncxx::types::b_date(std::chrono::system_clock::now());
        users.update_one(filter_builder.view(), bsoncxx::builder::stream::document{} << "$set" << token_info.view() << bsoncxx::builder::stream::finalize);

        return token;
    }


    void deleteUser(const std::string& userId) {
        bsoncxx::builder::stream::document filter_builder;
        filter_builder << "userId" << userId;
        auto result = users.delete_one(filter_builder.view());
        if (!result || result->deleted_count() == 0) {
            throw std::runtime_error("User not found or already deleted.");
        }
    }

    // Проверка подключения к указанному URL.
    bool checkConnection(const std::string& url) {
        httplib::Client cli(url);
        auto res = cli.Get("/");
        return res && res->status == 200;
    }

    // Получение прав доступа пользователя.
    std::string getPermission(const std::string& userId) {
        bsoncxx::builder::stream::document filter_builder;
        filter_builder << "userId" << userId;
        auto result = users.find_one(filter_builder.view());
        if (result) {
            auto permission_field = result->view()["permission"];
            if (permission_field) {
                return std::string(permission_field.get_string().value);
            }
        }
        return "No Permission";
    }

    // Установка прав доступа пользователя.
    void setPermission(const std::string& userId, const std::string& permission) {
        bsoncxx::builder::stream::document filter_builder, update_builder;
        filter_builder << "userId" << userId;
        update_builder << "$set" << bsoncxx::builder::stream::open_document
            << "permission" << permission
            << bsoncxx::builder::stream::close_document;
        users.update_one(filter_builder.view(), update_builder.view(), mongocxx::options::update().upsert(true));
    }

    // Получение информации о пользователе.
    std::string getUserInfo(const std::string& userId) {
        bsoncxx::builder::stream::document filter_builder;
        filter_builder << "userId" << userId;
        auto result = users.find_one(filter_builder.view());
        if (result) {
            return bsoncxx::to_json(result->view());
        }
        else {
            throw std::runtime_error("User not found.");
        }
    }

    // Обновление имени пользователя.
    void updateUserName(const std::string& userId, const std::string& newName) {
        bsoncxx::builder::stream::document filter_builder, update_builder;
        filter_builder << "userId" << userId;
        update_builder << "$set" << bsoncxx::builder::stream::open_document
            << "username" << newName
            << bsoncxx::builder::stream::close_document;
        users.update_one(filter_builder.view(), update_builder.view());
    }

    // Запуск HTTP-сервера.
    void startServer() {
        constexpr const char* ip_address = "0.0.0.0";
        constexpr int port = 8080;
        std::cout << "Authorization server is running on http://" << ip_address << ":" << port << std::endl;
        svr.listen(ip_address, port);
    }
};

int main() {
    try {
        mongocxx::instance instance{}; // Инициализация MongoDB.

        const std::string cloud_uri = "mongodb+srv://zarrock1:123@myclouddb.0i3hs.mongodb.net/?retryWrites=true&w=majority&appName=MyCloudDB";

        AuthorizationServer authServer(cloud_uri); // Создание объекта сервера авторизации.

        mongocxx::client test_conn{ mongocxx::uri{cloud_uri} };
        mongocxx::database test_db = test_conn["admin"];
        const auto ping_cmd = bsoncxx::builder::basic::make_document(bsoncxx::builder::basic::kvp("ping", 1));
        test_db.run_command(ping_cmd.view());
        std::cout << "Pinged your deployment. You successfully connected to MongoDB!" << std::endl;

        std::cout << R"(
GET /api/token?userId=<userId>        - Generate a token for a user.
GET /api/permissions?userId=<userId>  - Get permissions for a user.
POST /api/permissions?userId=<userId>&permission=<permission> - Set permissions for a user.
GET /api/userinfo?userId=<userId>     - Get user info in JSON format.
POST /api/updatename?userId=<userId>&newName=<newName> - Update user name.
)";

        authServer.startServer(); // Запуск сервера авторизации.
    }
    catch (const std::exception& e) {
        std::cerr << "Exception: " << e.what() << std::endl;
    }

    return 0;
}
