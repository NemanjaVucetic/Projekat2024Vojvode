import { Component, OnInit } from '@angular/core';
import { JwtHelperService } from '@auth0/angular-jwt';
import { NotificationService } from '../notification.service'; // Make sure you have a NotificationService

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.css']
})
export class NotificationComponent implements OnInit {

  notifications: any[] = [];  // Inicijalizujemo sa praznim nizom
  userId: string | null = null;

  constructor(
    private notificationService: NotificationService,  // Service for API calls
    private jwtHelper: JwtHelperService
  ) { }

  ngOnInit(): void {
    // Uzimamo token iz localStorage
    const token = localStorage.getItem('token');

    if (token) {
      try {
        // Dekodiramo token koristeći JwtHelperService
        const decodedToken = this.jwtHelper.decodeToken(token);
        console.log('Decoded Token:', decodedToken);

        // Izvlačimo userId iz dekodiranog tokena
        this.userId = decodedToken?.userId || null;

        // Proveravamo da li je userId prisutan
        if (this.userId) {
          // Pozivamo API da uzmemo notifikacije za korisnika
          this.notificationService.getNotificationsByUserId(this.userId).subscribe(
            (response: any[]) => {
              // Filter out the notifications with status 'discarded'
              this.notifications = response.filter(notification => notification.status !== 'discarded');
            },
            error => {
              console.error('Greška pri učitavanju notifikacija:', error);
            }
          );
        } else {
          console.error('User ID nije dostupan u tokenu');
        }
      } catch (error) {
        console.error('Greška pri dekodiranju tokena:', error);
      }
    } else {
      console.error('JWT token nije pronađen u localStorage');
    }
  }

  discardNotification(notificationId: string): void {
    // Pozivamo API endpoint za postavljanje statusa na "discarded"
    this.notificationService.discardNotification(notificationId).subscribe(
      (response) => {
        // Ako je poziv uspešan, ažuriramo status notifikacije u lokalnoj listi
        const updatedNotification = this.notifications.find(notification => notification.id === notificationId);
        if (updatedNotification) {
          updatedNotification.status = 'discarded';  // Postavljamo status na "discarded"
        }
      },
      error => {
        console.error('Greška pri ažuriranju statusa notifikacije:', error);
      }
    );
  }
}
