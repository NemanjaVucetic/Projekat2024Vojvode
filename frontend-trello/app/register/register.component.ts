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
  successMessage: string | null = null; // Poruka o uspehu
  errorMessage: string | null = null;   // Poruka o grešci

  constructor(
    private formBuilder: FormBuilder,
    private userService: UserService,
    private router: Router
  ) {
    this.registerForm = this.formBuilder.group({
      name: ['', [Validators.required]],
      surname: ['', [Validators.required]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      userRole: ['user', [Validators.required]] // Podrazumevani izbor je "user"
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
        this.successMessage = 'Uspešno ste se registrovali! Proverite svoj email za verifikaciju.';
        
        // Odmah preusmeri na login
        this.router.navigate(['/login']);
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Došlo je do greške prilikom registracije. Molimo pokušajte ponovo.';
        console.error(error);
      }
    });
  }
}
