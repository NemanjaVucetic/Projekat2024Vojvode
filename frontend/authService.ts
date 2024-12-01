import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(private http: HttpClient) {}

  // Dobavlja JWT token iz localStorage
  getToken(): string | null {
    const token = localStorage.getItem('jwt');
    console.log('Dobavljeni token iz localStorage:', token);  // Dodaj log za token
    return token;
  }

  // Dekodira JWT token
  getDecodedToken(): any {
    const token = this.getToken();
    if (token) {
      try {
        const decoded = jwtDecode(token);
        console.log('Dekodirani token:', decoded);  // Dodaj log za dekodirani token
        return decoded;
      } catch (error) {
        console.error('Nevažeći JWT token:', error);
        return null;
      }
    }
    return null;
  }

  // Dobavlja korisnički ID iz dekodiranog tokena
  getUserId(): string | null {
    const decodedToken = this.getDecodedToken();
    if (decodedToken) {
      console.log('Korisnički ID iz tokena:', decodedToken.userId);  // Dodaj log za korisnički ID
      return decodedToken.userId;
    }
    console.log('Nema korisničkog ID-a u dekodiranom tokenu.');  // Log ako nema korisničkog ID-a
    return null;
  }
}
