import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReactiveFormsModule } from '@angular/forms';  // Dodaj ReactiveFormsModule
import { AppRoutingModule } from './app-routing.module';
import { JwtHelperService, JwtModule } from '@auth0/angular-jwt';  // Dodaj import za JwtModule
import { ProjectService } from './project.service';
import { ProjectDetailsComponent } from './project-details/project-details.component';
import { ProjectCreateComponent } from './project-create/project-create.component';
import { NotificationComponent } from './notification/notification.component';
import { NotificationService } from './notification.service';
import { ProjectListComponent } from './project-list/project-list.component';
import { UserListComponent } from './user-list/user-list.component';

// Funkcija za dobijanje tokena iz localStorage
export function tokenGetter() {
  return localStorage.getItem('token');
}

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    RegisterComponent,
    ProjectDetailsComponent,
    ProjectCreateComponent,
    NotificationComponent,
    ProjectListComponent,
    UserListComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule,  // Dodaj AppRoutingModule ovde
    ReactiveFormsModule,  // Dodaj ReactiveFormsModule ovde
    JwtModule.forRoot({  // Inicijalizuj JwtModule sa konfiguracijom
      config: {
        tokenGetter: tokenGetter,
        allowedDomains: ['localhost:8000'],  // Dodaj domen vašeg API-ja
        disallowedRoutes: []  // Opcionalno možete dodati rute koje ne koriste JWT
      }
    })
  ],
  providers: [ProjectService, JwtHelperService,NotificationService],  // Dodaj JwtHelperService u providere
  bootstrap: [AppComponent]
})
export class AppModule { }
