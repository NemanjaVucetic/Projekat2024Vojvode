import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { UserService } from '../user.service'; // Import the UserService

interface Member {
  id: string;
  name: string;
  surname: string;
  email: string;
  password: string;
  userRole: string;
  enabled: boolean;
}

interface Project {
  id: string;
  name: string;
  manager: string;
  members: Member[];
  deadline: string;
  maxMembers: number;
  minMembers: number;
}

@Component({
  selector: 'app-project-list',
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.css'],
})
export class ProjectListComponent implements OnInit {
  projects: Project[] = [];
  users: Member[] = []; // To hold the list of users to add to the project
  selectedUserId: string | null = null; // The selected user to add
  private apiUrl = 'http://localhost:8000/api/project/'; // API URL
  private userId: string | null = null;

  constructor(private http: HttpClient, private userService: UserService) {}

  ngOnInit(): void {
    this.loadProjects();
    this.getUserIdFromToken(); // Get userId from token
    this.loadUsers(); // Load users to be displayed for selection
  }

  loadProjects(): void {
    this.http.get<Project[]>(this.apiUrl).subscribe((projects) => {
      this.projects = projects.filter((project) => project.manager === this.userId);
    });
  }

  loadUsers(): void {
    this.userService.getUsers().subscribe((users) => {
      this.users = users;
    });
  }

  getUserIdFromToken(): void {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decodedToken = JSON.parse(atob(token.split('.')[1])); // Decode JWT token
        this.userId = decodedToken?.userId || null;
      } catch (error) {
        console.error('Error decoding token', error);
      }
    }
  }

  // Method to remove a member from a project
  removeUserFromProject(projectId: string, memberId: string): void {
    const removeUrl = `${this.apiUrl}${projectId}/removeUser/${memberId}`;
    this.http.put(removeUrl, {}).subscribe(
      () => {
        // Remove the member from the local list after successful removal
        this.projects.forEach((project) => {
          if (project.id === projectId) {
            project.members = project.members.filter((member) => member.id !== memberId);
          }
        });
      },
      (error) => {
        console.error('Error removing member:', error);
      }
    );
  }

  // Method to add a user to a project
  addUserToProject(projectId: string, userId: string): void {
    const addUrl = `${this.apiUrl}${projectId}/addUser/${userId}`;
    this.http.put(addUrl, {}).subscribe(
      () => {
        // After adding the user, update the members list in the local data
        this.projects.forEach((project) => {
          if (project.id === projectId) {
            this.userService.getUsers().subscribe((users) => {
              const user = users.find((u) => u.id === userId);
              if (user) {
                project.members.push(user); // Add the selected user to the project
              }
            });
          }
        });
      },
      (error) => {
        console.error('Error adding member:', error);
      }
    );
  }
}
