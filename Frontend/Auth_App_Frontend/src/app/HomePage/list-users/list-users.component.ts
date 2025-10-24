import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { UserService } from 'src/app/services/user.service';

interface User {
  username: string;
  full_name: string;
  email: string;
  phone_number: string;
  date_of_birth: string;
  address: string;
  create_date: string;
  account_status: string;
  oauth_provider: string;
  oauth_id: string;
  user_type: string;
}

@Component({
  selector: 'app-list-users',
  templateUrl: './list-users.component.html',
  styleUrls: ['./list-users.component.css']
})
export class ListUsersComponent implements OnInit {
  users: User[] = [];
  loading: boolean = false;
  error: string = '';
  email: string = '';
  username: string = '';
  
  constructor(private http: HttpClient, private userService: UserService) {}

  ngOnInit() {
    this.email = localStorage.getItem('EMAIL') || '';
    if (this.email) {
      this.loadUsers();
    } else {
      this.error = 'No email found in local storage';
    }
  }

  loadUsers() {
    this.loading = true;
    this.error = '';

    this.userService.listUsers(this.email).subscribe({
      next: (response) => {
        this.users = Array.isArray(response) ? response : [response];
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load users: ' + err.message;
        this.loading = false;
      }
    });
  }
  disableUser(user: User) {
    if (confirm(`Are you sure you want to disable ${user.username}?`)) {
      this.userService.disableUser(this.username, user.username).subscribe({
        next: (response) => {
          console.log('User disabled successfully:', response);
          // Update the user status locally
          user.account_status = 'disabled';
          // Or reload the users list to get fresh data
          this.loadUsers();
        },
        error: (err) => {
          this.error = 'Failed to disable user: ' + err.message;
          console.error('Error disabling user:', err);
        }
      });
    }
  }

  enableUser(user: User) {
    if (confirm(`Are you sure you want to enable ${user.username}?`)) {
      this.userService.enableUser(this.username, user.username).subscribe({
        next: (response) => {
          console.log('User enabled successfully:', response);
          // Update the user status locally
          user.account_status = 'active';
          // Or reload the users list to get fresh data
          this.loadUsers();
        },
        error: (err) => {
          this.error = 'Failed to enable user: ' + err.message;
          console.error('Error enabling user:', err);
        }
      });
    }
  }
}
