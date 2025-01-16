package com.kaushik.model;

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

    public boolean activationStatus() {
        return this.isDeleted;
    }
}
