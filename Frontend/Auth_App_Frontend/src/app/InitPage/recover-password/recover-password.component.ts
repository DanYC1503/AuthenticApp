import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'environments/environment';

@Component({
  selector: 'app-recover-password',
  templateUrl: './recover-password.component.html',
  styleUrls: ['./recover-password.component.css']
})
export class RecoverPasswordComponent {
  email: string = '';

  constructor(private http: HttpClient) {}

  onRecoverPassword() {
    if (!this.email) {
      alert('Por favor ingrese su correo electrónico');
      return;
    }

    this.http.get(`${environment.API_GATEWAY_URL}auth/passwordToken`, {
      params: { email: this.email },
      withCredentials: true
    }).subscribe({
      next: (res: any) => {
        alert('Si el correo existe, se ha enviado un enlace de recuperación');
        console.log('Password token response:', res);
      },
      error: (err) => {
        console.error('Error requesting password token:', err);
        alert('Ocurrió un error. Intente nuevamente.');
      }
    });
  }
}
