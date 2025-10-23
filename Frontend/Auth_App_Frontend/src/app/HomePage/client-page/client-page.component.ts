import { Component, OnInit } from '@angular/core';
import { AuthServiceService } from 'src/app/services/auth-service.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-client-page',
  templateUrl: './client-page.component.html',
  styleUrls: ['./client-page.component.css']
})
export class ClientPageComponent implements OnInit {
  user: any = null;
  isCollapsed: boolean = false;
  profileOpen: boolean = false;
  public showAdminMenu: boolean = false;
  selectedAction: 'home' | 'update' | 'delete' = 'home';

  constructor(private userService: UserService, private authService: AuthServiceService) {}

  ngOnInit(): void {
  // Subscribe to the admin menu state reactively
  this.authService.showAdminMenu$.subscribe(state => {
    this.showAdminMenu = state;
  });

  const username = localStorage.getItem('USERNAME');
  if (!username) return;

  this.userService.getUserInfo(username).subscribe({
    next: (res) => {
      this.user = res.user;
      console.log('User info loaded:', this.user);
      // Optional: update admin menu if user info has type
      if (this.user.user_type) {
        this.authService.setShowAdminMenu(this.user.user_type === 'admin');
      }
    },
    error: (err) => console.error('Error fetching user info:', err)
  });
}
selectAction(action: 'home' | 'update' | 'delete') {
    this.selectedAction = action;
  }


updateAdminMenu() {
  const isAdmin = this.user?.user_type === 'admin';
  this.authService.setShowAdminMenu(isAdmin);
}

}
