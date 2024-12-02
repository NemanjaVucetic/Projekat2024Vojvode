import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  private apiUrl = 'http://localhost:8000/api/notification/byUserId';  // Your API base URL
  private discardUrl = 'http://localhost:8000/api/notification/discard';  // URL for discarding notification

  constructor(private http: HttpClient) { }

  // Get notifications by userId
  getNotificationsByUserId(userId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/${userId}`);
  }

  // Discard a notification by its ID
  discardNotification(notificationId: string): Observable<void> {
    return this.http.put<void>(`${this.discardUrl}/${notificationId}`, {});  // Sending a PUT request to update the notification status
  }
}
