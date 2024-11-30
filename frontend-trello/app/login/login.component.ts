import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent {
  email: string = '';
  password: string = '';
  errorMessage: string | null = null;

  constructor(private http: HttpClient, private router: Router) {}

  onLogin(): void {
    const loginData = { email: this.email, password: this.password };

    this.http.post('http://localhost:8000/api/user/login', loginData).subscribe({
      next: (response: any) => {
        // Handle successful login
        console.log('Login successful', response);
        localStorage.setItem('token', response.token); // Assuming the token is returned
        this.router.navigate(['/project']); // Redirect to dashboard or desired page
      },
      error: (error) => {
        // Handle error
        console.error('Login failed', error);
        this.errorMessage = 'Invalid email or password. Please try again.';
      },
    });
  }
}
