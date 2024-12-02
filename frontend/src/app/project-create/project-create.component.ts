import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-project-create',
  templateUrl: './project-create.component.html',
  styleUrls: ['./project-create.component.css']
})
export class ProjectCreateComponent {
  projectForm: FormGroup;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private http: HttpClient
  ) {
    this.projectForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(255)]],
      deadline: ['', [Validators.required]],
      maxMembers: [22, [Validators.required, Validators.min(1)]],
      minMembers: [1, [Validators.required, Validators.min(1)]]
    });
  }

  get formControls() {
    return this.projectForm.controls;
  }

  // Metoda za slanje podataka na server
  onSubmit() {
    if (this.projectForm.valid) {
      this.isSubmitting = true;
      const projectData = this.projectForm.value;
      
      // Pošaljite POST zahtev na backend
      this.http.post('http://localhost:8000/api/project/', projectData)
        .subscribe(
          (response) => {
            console.log('Project created successfully', response);
            this.isSubmitting = false;
            this.projectForm.reset();
          },
          (error) => {
            console.error('Error creating project', error);
            this.isSubmitting = false;
          }
        );
    }
  }
}
