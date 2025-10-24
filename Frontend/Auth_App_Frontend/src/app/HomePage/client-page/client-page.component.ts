import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
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
  selectedAction: 'home' | 'update' | 'delete' | 'listUsers' = 'home';

  constructor(private userService: UserService, private authService: AuthServiceService, private router: Router) {}

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
      localStorage.setItem('EMAIL', this.user.email);
      if (this.user.user_type) {
        this.authService.setShowAdminMenu(this.user.user_type === 'admin');
      }
    },
    error: (err) => console.error('Error fetching user info:', err)
  });
}
selectAction(action: 'home' | 'update' | 'delete' | 'listUsers') {
  this.selectedAction = action;

  // Reload user info only when going to home
  if (action === 'home') {
    this.reloadUser();
  }
}

reloadUser() {
  const username = localStorage.getItem('USERNAME');
  if (!username) return;

  this.userService.getUserInfo(username).subscribe({
    next: (res) => {
      this.user = res.user;
      console.log('User info reloaded:', this.user);
    },
    error: (err) => console.error('Error fetching user info:', err)
  });
}


updateAdminMenu() {
  const isAdmin = this.user?.user_type === 'admin';
  this.authService.setShowAdminMenu(isAdmin);
}
logout() {
  // Clear cookies, localStorage, sessionStorage
  localStorage.clear();
  sessionStorage.clear();

  // Optional: if using a cookie service
  // this.cookieService.deleteAll('/', window.location.hostname);

  // Redirect to login page
  this.router.navigate(['/login']);
}

}
