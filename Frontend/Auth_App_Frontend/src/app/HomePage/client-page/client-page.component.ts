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
  const queryParams = new URLSearchParams(window.location.search);
  const usernameParam = queryParams.get('username');
  const emailParam = queryParams.get('email');
  const userTypeParam = queryParams.get('userType');

  // If OAuth redirect, save info to localStorage
  if (usernameParam) {
    localStorage.setItem('USERNAME', usernameParam);
    localStorage.setItem('EMAIL', emailParam || '');
    localStorage.setItem('USER_TYPE', userTypeParam || '');
    // Optionally clear query params from URL
    window.history.replaceState({}, document.title, '/ClientDashboard');
  }

  // Determine username to fetch from backend
  const storedUsername = localStorage.getItem('USERNAME');

  if (storedUsername) {
    this.userService.getUserInfo(storedUsername).subscribe({
      next: (res) => {
        this.user = res.user;
        console.log('User info loaded:', this.user);
        localStorage.setItem('EMAIL', this.user.email);
        // Set user type for admin menu: prefer param from OAuth if exists
        const type = userTypeParam || this.user.user_type || localStorage.getItem('USER_TYPE');
        if (type) {
          this.authService.setShowAdminMenu(type === 'admin');
        }
      },
      error: (err) => console.error('Error fetching user info:', err)
    });
  }

  // Subscribe reactively to admin menu state
  this.authService.showAdminMenu$.subscribe(state => {
    this.showAdminMenu = state;
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
