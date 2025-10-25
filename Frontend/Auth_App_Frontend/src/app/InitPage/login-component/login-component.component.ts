import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { environment } from 'environments/environment';
import { CsrfService } from 'src/app/services/csrf.service';
import { Router } from '@angular/router';
import { AuthServiceService } from 'src/app/services/auth-service.service';


@Component({
  selector: 'app-login-component',
  templateUrl: './login-component.component.html',
  styleUrls: ['./login-component.component.css']
})
export class LoginComponentComponent implements OnInit {
  usuario: string = '';
  password: string = '';
  showPassword: boolean = false;
  constructor(private csrfService: CsrfService,
    private http: HttpClient,
    private router: Router,
    private authService: AuthServiceService
  ) {}

  ngOnInit() {
    // 1️Clear previous CSRF cookie
    document.cookie.split(";").forEach((c) => {
      document.cookie = c
        .replace(/^ +/, "")
        .replace(/=.*/, `=;expires=${new Date(0).toUTCString()};path=/`);
    });

    // 2️Clear localStorage
    localStorage.removeItem('XSRF-TOKEN');
    localStorage.removeItem('EMAIL');
    localStorage.removeItem('SESSION_TOKEN');
    localStorage.removeItem('USERNAME');
    localStorage.removeItem('ISADMIN');

    // Now load the new CSRF token
    this.csrfService.loadCsrfToken().subscribe({
      next: () => {
        const token = this.csrfService.getTokenFromCookie();
        if (token) localStorage.setItem('XSRF-TOKEN', token);
      },
      error: (err) => {
        console.error('Failed to load CSRF token', err);
      }
    });
  }

  onLogin() {
    const body = {
      username: this.usuario,
      password: this.password
    };


  this.http.post(`${environment.API_GATEWAY_URL}auth/login`, body, { withCredentials: true })
    .subscribe({
      next: (res: any) => {
        console.log('Login success:', res);

        // The backend sets the cookie via Set-Cookie header,
        // the browser will automatically store it (thanks to withCredentials).
        localStorage.setItem('USERNAME', this.usuario);
        if (res.session_token) {
          localStorage.setItem('SESSION_TOKEN', res.session_token);
          localStorage.setItem('SESSION_EXPIRES', res.expires.toString());
          console.log('Session token stored manually');
        }

        // After login, get user type
        this.getUserType();
      },
      error: (err) => {
        console.error('Login failed:', err);
        alert('Invalid credentials or session error');
      }
    });
  }
  loginWithGoogle(): void {
    const backendUrl = `${environment.API_GATEWAY_URL}auth/google/login`;
    window.location.href = backendUrl; // Redirect to backend for OAuth
  }
  getUserType() {
    const username = localStorage.getItem('USERNAME');
    if (!username) return;

    const body = { username };

    this.http.post<{ type: string }>(
      `${environment.API_GATEWAY_URL}users/retrieve/type`,
      body,
      { withCredentials: true }
    ).subscribe({
      next: (res) => {
        console.log('User type:', res.type);

        const isAdmin = res.type === 'admin';
        // Update service and localStorage
        this.authService.setShowAdminMenu(isAdmin);

        // Navigate after setting admin state
        this.router.navigate(['/ClientDashboard']);
      },
      error: (err) => {
        console.error('Failed to retrieve user type', err);
        alert('Cannot determine user type');
      }
    });
  }
}
