import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthServiceService } from 'src/app/services/auth-service.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-delete-user-page',
  templateUrl: './delete-user-page.component.html',
  styleUrls: ['./delete-user-page.component.css']
})
export class DeleteUserPageComponent {
  showConfirmBox = false;
  confirmEmail = '';
  successPopup = false;
  loading = false;
  message = '';

  constructor(private userService: UserService,
    private authService: AuthServiceService,
    private router: Router
  ){}

  requestDelete(payload: any) {
  this.loading = true;

  this.authService.requestDeleteToken(payload).subscribe({
    next: (res: any) => {
      console.log('Token response:', res);

      const deleteToken = res.deleteAuthToken;
      if (!deleteToken) {
        this.message = 'Error al generar el token de eliminación.';
        this.loading = false;
        return;
      }

      // Call deleteUser with both payload and token
      this.userService.deleteUser(payload, deleteToken).subscribe({
        next: () => {
          console.log('User deleted successfully');
          this.loading = false;
          this.successPopup = true;
        },
        error: (err) => {
          console.error('Error deleting user:', err);
          this.loading = false;
        }
      });

      this.showConfirmBox = false;
    },
    error: (err) => {
      console.error('Error requesting delete token:', err);
      this.message = 'Error al validar el correo electrónico.';
      this.loading = false;
    }
  });
}


  closeDeleteSuccessPopup() {
  this.successPopup = false;
  this.logout()
  }
  logout() {
    // Clear cookies, localStorage, sessionStorage
    localStorage.clear();
    sessionStorage.clear();
    // Redirect to login page
    this.router.navigate(['/login']);
  }

}
