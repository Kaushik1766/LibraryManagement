package com.kaushik.database;

import com.kaushik.model.Student;

import java.util.ArrayList;

import com.kaushik.model.Admin;
import com.kaushik.model.Book;

public class Database {
    ArrayList<Admin> Admins;
    ArrayList<Student> Students;
    ArrayList<Book> Books;

    public void addUser(Object user) throws Exception {
        if (user instanceof Admin) {
            Admins.add((Admin) user);
        } else if (user instanceof Student) {
            Students.add((Student) user);
        } else {
            throw new Exception("Object of invalid type provided.");
        }
    }

    public Object searchUser(String userId, String role) throws Exception {
        if (role == "Admin") {
            for (Admin ad : Admins) {
                if (ad.getUserId() == userId && !ad.activationStatus()) {
                    return ad;
                }
            }
        } else {
            for (Student std : Students) {
                if (std.getUserId() == userId && !std.activationStatus()) {
                    return std;
                }
            }
        }
        throw new Exception("User not found.");
    }

    public Object searchUserByEmail(String email, String role) throws Exception {
        if (role == "Admin") {
            for (Admin ad : Admins) {
                if (ad.getEmail() == email && !ad.activationStatus()) {
                    return ad;
                }
            }
        } else {
            for (Student std : Students) {
                if (std.getEmail() == email && !std.activationStatus()) {
                    return std;
                }
            }
        }
        throw new Exception("User not found.");
    }

    public void deleteUser(String userId, String role) throws Exception {
        if (role == "Admin") {
            for (Admin ad : Admins) {
                if (ad.getUserId() == userId && !ad.activationStatus()) {
                    ad.deactivate();
                    return;
                }
            }
        } else {
            for (Student std : Students) {
                if (std.getUserId() == userId && !std.activationStatus()) {
                    std.deactivate();
                    return;
                }
            }
        }
        throw new Exception("User not found.");
    }

    public String getLastUserId() {
        if (Students.size() == 0) {
            return "1";
        } else {
            String lastId = Students.get(Students.size() - 1).getUserId();
            return String.valueOf(Integer.parseInt(lastId) + 1);
        }
    }

    public void addBook(Book book, Object user) throws Exception {
        if ((user instanceof Admin)) {
            Books.add(book);
        } else {
            throw new Exception("User is not an admin.");
        }
    }

    public Book searchBook(String bookId) throws Exception {
        for (Book b : Books) {
            if (b.getId() == bookId) {
                return b;
            }
        }
        throw new Exception("Book with the provided bookId does not exist.");
    }
}
