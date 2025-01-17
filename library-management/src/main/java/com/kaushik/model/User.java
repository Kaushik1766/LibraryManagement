package com.kaushik.model;

import at.favre.lib.crypto.bcrypt.BCrypt.HashData;

abstract class User {
    String name;
    String userId;
    String password;
    String email;
    boolean isDeleted = false;

    abstract void displayInfo(String userId);

    public void deactivate() {
        this.isDeleted = true;
    }

    public void activate() {
        this.isDeleted = false;
    }

    public String getEmail() {
        return email;
    }

    public String getPassword() {
        return password;
    }

    public boolean activationStatus() {
        return this.isDeleted;
    }
}
