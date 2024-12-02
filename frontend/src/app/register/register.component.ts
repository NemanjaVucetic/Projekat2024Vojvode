import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UserService } from '../user.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  registerForm: FormGroup;
  isLoading = false;
  successMessage: string | null = null; // Success message
  errorMessage: string | null = null;   // Error message

  constructor(
    private formBuilder: FormBuilder,
    private userService: UserService,
    private router: Router
  ) {
    this.registerForm = this.formBuilder.group({
      name: ['', [Validators.required]],
      surname: ['', [Validators.required]],
      username: ['', [Validators.required]],  // Added username field
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      userRole: ['user', [Validators.required]] // Default role is "user"
    });
  }

  get f() {
    return this.registerForm.controls;
  }

  onSubmit() {
    if (this.registerForm.invalid) {
      return;
    }

    this.isLoading = true;
    this.userService.register(this.registerForm.value).subscribe({
      next: (response) => {
        this.isLoading = false;
        this.successMessage = 'You have successfully registered! Please check your email for verification.';
        
        // Immediately redirect to login
        this.router.navigate(['/login']);
      },
      error: (error) => {
        this.isLoading = false;
        // Check for specific errors like email or username already taken
        if (error?.error?.message === 'Email or Username already taken') {
          this.errorMessage = 'This email or username is already in use. Please choose another one.';
        } else {
          this.errorMessage = 'An error occurred during registration. Please try again.';
        }
        console.error(error);
      }
    });
  }
}
