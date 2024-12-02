// user-list.component.ts
import { Component, OnInit } from '@angular/core';
import { UserService } from '../user.service';

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  users: any[] = [];  // Use 'any[]' to allow any structure

  constructor(private userService: UserService) {}

  ngOnInit(): void {
    this.userService.getUsers().subscribe(
      (response: any[]) => {  // Specify the response type as 'any[]'
        this.users = response;
      },
      (error) => {
        console.error('Error fetching user data', error);
      }
    );
  }
}

