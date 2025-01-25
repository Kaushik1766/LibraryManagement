package com.kaushik.service;

import java.util.Scanner;

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

    public void borrow() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Book Id: ");
        String bookId = sc.nextLine();
        try {
            if (user instanceof Student) {
                Book b = db.searchBook(bookId);
                Student student = (Student) user;
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
        } finally {
            sc.close();
        }
    }

    public void returnBook() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Book Id: ");
        String bookId = sc.nextLine();
        try {
            if (user instanceof Student) {
                Book b = db.searchBook(bookId);
                Student student = (Student) user;
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
        } finally {
            sc.close();
        }
    }

    public void addBook() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Book Name: ");
        String name = sc.nextLine();
        System.out.println("Enter ISBN: ");
        String isbn = sc.nextLine();
        System.out.println("Enter Author: ");
        String author = sc.nextLine();
        System.out.println("Enter Category: ");
        String category = sc.nextLine();
        System.out.println("Enter Price: ");
        double price = sc.nextDouble();
        System.out.println("Enter Path: ");
        String path = sc.nextLine();
        System.out.println("Enter Availability: ");
        boolean available = sc.nextBoolean();
        sc.close();
        Book book = new Book(name, isbn, author, category, price, path, available, db.getLastBookId());
        try {
            if (user instanceof Admin) {
                db.addBook(book, user);
                System.out.println("Book added successfully.");
            } else {
                throw new Exception("Logged in as Student.");
            }
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }

    }

    public void registerUser() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Name: ");
        String name = sc.nextLine();
        System.out.println("Enter Email: ");
        String email = sc.nextLine();
        System.out.println("Enter Password: ");
        String password = sc.nextLine();
        System.out.println("Enter Role: ");
        String role = sc.nextLine();
        role.toLowerCase();

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
        } finally {
            sc.close();
        }
    }

    public void login() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Email: ");
        String email = sc.nextLine();
        System.out.println("Enter Password: ");
        String password = sc.nextLine();
        System.out.println("Enter Role: ");
        String role = sc.nextLine();
        role.toLowerCase();

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
        } finally {
            sc.close();
        }
    }

    public void SearchBook() {
        Scanner sc = new Scanner(System.in);
        System.out.println("Enter Book Id: ");
        String bookId = sc.nextLine();
        try {
            Book b = db.searchBook(bookId);
            System.out.println("Book found: " + b.getName());
        } catch (Exception e) {
            System.out.println(e.getMessage());
        } finally {
            sc.close();
        }

    }

    public void logout() {
        this.user = null;
    }
}
