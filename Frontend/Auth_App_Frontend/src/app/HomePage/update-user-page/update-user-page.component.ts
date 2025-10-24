import { HttpClient } from '@angular/common/http';
import { Component, EventEmitter, Output } from '@angular/core';
import { environment } from 'environments/environment';
import { AuthServiceService } from 'src/app/services/auth-service.service';
import { UserService } from 'src/app/services/user.service';
import { ClientPageComponent } from '../client-page/client-page.component';

@Component({
  selector: 'app-update-user-page',
  templateUrl: './update-user-page.component.html',
  styleUrls: ['./update-user-page.component.css']
})
export class UpdateUserPageComponent {
  user: any = {};
  loading = false;
  message = '';
  showConfirmBox = false;
  confirmEmail = '';

  // fields
  full_name = '';
  email = '';
  phone_number = '';
  date_of_birth = '';
  address = '';
  username: string = localStorage.getItem('USERNAME') || '';
  successPopup = false;
  @Output() actionChanged = new EventEmitter<'home' | 'update' | 'delete'>();


  constructor(
    private authService: AuthServiceService,
    private userService: UserService
  ) {}

  ngOnInit(): void {
    const username = localStorage.getItem('USERNAME');
    if (!username) return;

    this.userService.getUserInfo(username).subscribe({
      next: (res) => {
        this.user = res.user;
        console.log('User info loaded:', this.user);

        this.full_name = this.user.full_name || '';
        this.email = this.user.email;
        this.phone_number = this.user.phone_number || '';
        this.date_of_birth = this.user.date_of_birth || '';
        this.address = this.user.address || '';
        this.username = this.user.username || '';
      },
      error: (err) => console.error('Error fetching user info:', err)
    });
  }

  confirmUpdate() {
    if (!this.username) {
      this.message = 'No se encontró el usuario actual en la sesión.';
      return;
    }
    if (!this.full_name.trim() && !this.email.trim()) {
      this.message = 'No hay campos para actualizar.';
      return;
    }
    this.showConfirmBox = true;
  }

  requestTokenUpdate(payload: any) {
    this.loading = true;

    this.authService.requestUpdateToken(payload).subscribe({
      next: (res: any) => {
        console.log('Token update response:', res);

        const updateToken = res.updateAuthToken;
        if (!updateToken) {
          this.message = 'Error al generar el token de actualización.';
          this.loading = false;
          return;
        }

        // Now update user info
        this.updateUserInfo(updateToken);
        this.showConfirmBox = false;
      },
      error: (err) => {
        console.error('Error requesting update token:', err);
        this.message = 'Error al validar el correo electrónico.';
        this.loading = false;
      }
    });
  }

  updateUserInfo(updateToken: string) {
  const body: any = { username: this.username, email: this.email };

  if (this.full_name.trim()) body.full_name = this.full_name.trim();
  if (this.phone_number.trim()) body.phone_number = this.phone_number.trim();
  if (this.address.trim()) body.address = this.address.trim();
  if (this.date_of_birth) {
    const date = new Date(this.date_of_birth);
    body.date_of_birth = date.toISOString().split('T')[0];
  }

  this.userService.updateUserInfo(body, updateToken).subscribe({
    next: (res) => {
      console.log('Update success:', res);
      this.loading = false;
      this.successPopup = true; // show the success popup
    },
    error: (err) => {
      console.error('Update failed:', err);
      this.message = 'Error al actualizar la información.';
      this.loading = false;
    }
  });
}


  cancelConfirm() {
    this.showConfirmBox = false;
    this.confirmEmail = '';
  }
  closeSuccessPopup() {
  this.successPopup = false;
  this.reloadUser();
  this.actionChanged.emit('home'); // tell parent to switch to home
}
reloadUser() {
  const username = localStorage.getItem('USERNAME');
  if (!username) return;

  this.userService.getUserInfo(username).subscribe({
    next: (res) => {
      this.user = res.user;
    },
    error: (err) => console.error('Error fetching user info:', err)
  });
}
}
