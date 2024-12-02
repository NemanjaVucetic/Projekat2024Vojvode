import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProjectService {

  private apiUrl = 'http://localhost:8000/api/project/user/';  // API URL

  constructor(private http: HttpClient) { }

  getProjectDetails(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}${userId}`);
  }
}
