package com.kaushik.model;

import at.favre.lib.crypto.bcrypt.BCrypt.HashData;

public class Admin extends User {
    public Admin(String name, String userId, String password, String email) {
        this.name = name;
        this.userId = userId;
        this.password = password;
        this.email = email;
    }

    public String getUserId() {
        return this.userId;
    }

    public void displayInfo(String userId) {
        System.out.println("Name: " + name);
        System.out.println("UserId" + userId);
        System.out.println("Email: " + email);
    }

}
