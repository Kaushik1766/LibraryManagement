package com.kaushik.model;

public class Book {
    String Name;
    String ISBN;
    String id;
    String Author;
    String Category;
    double Price;
    String path;
    boolean available;
    boolean isDeleted = false;
    Student borrowedBy = null;

    public Student getBorrowedBy() {
        return borrowedBy;
    }

    public void setBorrowedBy(Student borrowedBy) {
        this.borrowedBy = borrowedBy;
    }

    public String getId() {
        return id;
    }

    public String getName() {
        return Name;
    }

    public Book(String Name, String ISBN, String Author, String Category, double Price, String path,
            boolean available, String id) {
        this.Name = Name;
        this.ISBN = ISBN;
        this.Author = Author;
        this.Category = Category;
        this.Price = Price;
        this.path = path;
        this.available = available;
        this.id = id;
    }
}
