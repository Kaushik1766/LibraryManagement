package com.kaushik;

import com.kaushik.model.Student;
import com.kaushik.service.Management;

public class Main {
    public static void main(String[] args) {
        Management M = new Management();
        while (true) {
            if (M.getUser() == null) {
                System.out.println("1. Register\n2. Login\n3. Signout\n");

            } else {
                if (M.getUser() instanceof Student) {

                } else {

                }
            }
        }
    }
}
