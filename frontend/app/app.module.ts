import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReactiveFormsModule } from '@angular/forms';  
import { AppRoutingModule } from './app-routing.module';
import { JwtHelperService, JwtModule } from '@auth0/angular-jwt';  
import { ProjectCreateComponent } from './project-create/project-create.component';


export function tokenGetter() {
  return localStorage.getItem('jwtToken');
}

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    RegisterComponent,
    ProjectCreateComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule,  
    ReactiveFormsModule,  
    JwtModule.forRoot({  
      config: {
        tokenGetter: tokenGetter,
        disallowedRoutes: []  
      }
    })
  ],
  providers: [JwtHelperService],  
  bootstrap: [AppComponent]
})
export class AppModule { }
