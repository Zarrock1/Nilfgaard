#include <iostream>
#include <string>
#include <unordered_map>
#include <vector>
#include <json.hpp>

class AuthorizationServer {
private:
    std::unordered_map<std::string, std::string> userPermissions;
public:
    AuthorizationServer() {}

    void authenticateUser(const std::string& token) {
        std::cout << "Authenticating user with token: " << token << std::endl;
        // Эмуляция проверки с внешними сервисами
    }

    void setPermission(const std::string& userId, const std::string& permission) {
        userPermissions[userId] = permission;
    }

    std::string getPermission(const std::string& userId) {
        return userPermissions.count(userId) ? userPermissions[userId] : "No Permission";
    }
};

int main() {
    AuthorizationServer authServer;
    authServer.authenticateUser("example_token");
    authServer.setPermission("user1", "read_write");
    std::cout << "Permission for user1: " << authServer.getPermission("user1") << std::endl;

    return 0;
}
