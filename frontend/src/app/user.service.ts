import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'http://localhost:8000/api/user/';

  constructor(private http: HttpClient) {}

  register(userData: { name: string, surname: string, email: string, password: string, userRole: string }): Observable<any> {
    return this.http.post<any>(this.apiUrl, userData);
  }

  getUsers(): Observable<any[]> {  // Specify that it returns an Observable of any[]
    return this.http.get<any[]>(this.apiUrl);
  }
}
