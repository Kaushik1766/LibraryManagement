package com.kaushik;

import java.util.Scanner;

import com.kaushik.model.Student;
import com.kaushik.service.Management;

public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        Management M = new Management(sc);
        while (true) {
            if (M.getUser() == null) {
                System.out.println("1. Register\n2. Login\n3. Signout\n0. Exit\n");
                int choice = sc.nextInt();
                switch (choice) {
                    case 0:
                        System.exit(0);
                        break;
                    case 1:
                        M.registerUser();
                        break;
                    case 2:
                        M.login();
                        break;
                    case 3:
                        M.logout();
                        break;
                    default:
                        System.out.println("Invalid choice.");
                        break;
                }
            } else {
                if (M.getUser() instanceof Student) {
                    System.out.println("1. Borrow Book\n2. Return Book\n3. Search Book\n4. Signout\n0. Exit\n");
                    int choice = sc.nextInt();
                    switch (choice) {
                        case 0:
                            System.exit(0);
                            break;
                        case 1:
                            M.borrow();
                            break;
                        case 2:
                            M.returnBook();
                            break;
                        case 3:
                            M.SearchBook();
                            break;
                        case 4:
                            M.logout();
                            break;
                        default:
                            System.out.println("Invalid choice.");
                            break;
                    }
                } else {
                    System.out.println("1. Add Book\n2. Search Book\n3. Signout\n0. Exit\n");
                    int choice = sc.nextInt();
                    switch (choice) {
                        case 0:
                            System.exit(0);
                            break;
                        case 1:
                            M.addBook();
                            break;
                        case 2:
                            M.SearchBook();
                            break;
                        case 3:
                            M.logout();
                            break;
                        default:
                            System.out.println("Invalid choice.");
                            break;
                    }
                }
            }
        }
    }
}
