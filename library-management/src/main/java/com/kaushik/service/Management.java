package com.kaushik.service;

import com.kaushik.database.Database;
import com.kaushik.model.Admin;
import com.kaushik.model.Book;
import com.kaushik.model.Student;

public class Management {
    private Database db;
    private Object user = null;
    // private String secret = "asdfasdf";

    public Management() {
        this.db = new Database();
    }

    public Object getUser() {
        return user;
    }

    public void borrow(Student student, String bookId) {
        try {
            if (user instanceof Student) {
                Book b = db.searchBook(bookId);
                if (b.getBorrowedBy() != null) {
                    System.out.println("Book already borrowed");
                } else {
                    b.setBorrowedBy(student);
                    System.out.println("Book with id" + bookId + "has been borrowed by " + student.getName() + ".");
                }
            } else {
                throw new Exception("Logged in as Admin.");
            }
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }

    public void returnBook(Student student, String bookId) {
        try {
            if (user instanceof Student) {
                Book b = db.searchBook(bookId);
                if (b.getBorrowedBy() != student) {
                    System.out.println("Student didn't borrow this book.");
                } else {
                    b.setBorrowedBy(null);
                    System.out.println("Book with id" + bookId + "has been returned.");
                }
            } else {
                throw new Exception("Logged in as admin.");
            }
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }

    public void addBook() {

    }

    public void registerUser(String name, String email, String password, String role) {
        try {
            String lastId = db.getLastUserId();
            if (role == "admin") {
                Admin admin = new Admin(name, lastId, password, email);
                db.addUser(admin);
            } else {
                Student student = new Student(name, lastId, password, email);
                db.addUser(student);
            }
            System.out.println("Registered Successfully.");
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }

    public void login(String email, String password, String role) {
        try {
            if (role == "admin") {
                Admin admin = (Admin) db.searchUserByEmail(email, role);
                if (admin.getPassword() == password) {
                    user = admin;
                    System.out.println("Logged in successfully.");
                } else {
                    throw new Exception("Invalid Password.");
                }
            } else {
                Student student = (Student) db.searchUserByEmail(email, role);
                if (student.getPassword() == password) {
                    user = student;
                    System.out.println("Logged in successfully.");
                } else {
                    throw new Exception("Invalid Password.");
                }
            }
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }

    public void logout() {
        user = null;
    }
}
