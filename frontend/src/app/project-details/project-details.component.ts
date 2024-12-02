import { Component, OnInit } from '@angular/core';
import { ProjectService } from '../project.service';
import { JwtHelperService } from '@auth0/angular-jwt';

@Component({
  selector: 'app-project-details',
  templateUrl: './project-details.component.html',
  styleUrls: ['./project-details.component.css']
})
export class ProjectDetailsComponent implements OnInit {

  project: any = null;  // Inicijalizujemo sa null, što znači da nema podataka dok ih ne učitamo
  userId: string | null = null;

  constructor(
    private projectService: ProjectService,
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
          // Pozivamo API da uzmemo podatke o projektu
          this.projectService.getProjectDetails(this.userId).subscribe(
            response => {
              this.project = response;  // Pretpostavljamo da API vraća niz objekata, uzimamo prvi
            },
            error => {
              console.error('Greška pri učitavanju podataka o projektu:', error);
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
}