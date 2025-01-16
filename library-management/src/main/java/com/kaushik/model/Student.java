package com.kaushik.model;

public class Student extends User {

    public Student(String name, String userId, String password, String email) {
        this.name = name;
        this.userId = userId;
        this.password = password;
        this.email = email;
    }

    public String getUserId() {
        return userId;
    }

    public void displayInfo(String userId) {
        System.out.println("Name: " + name);
        System.out.println("UserId" + userId);
        System.out.println("Email: " + email);
    }
}
