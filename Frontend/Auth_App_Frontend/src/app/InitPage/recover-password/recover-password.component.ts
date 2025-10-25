import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthServiceService } from 'src/app/services/auth-service.service';

@Component({
  selector: 'app-recover-password',
  templateUrl: './recover-password.component.html',
  styleUrls: ['./recover-password.component.css']
})

export class RecoverPasswordComponent implements OnInit {

  email: string = '';
  token: string = '';
  newPassword: string = '';
  confirmPassword: string = '';
  isRedirected: boolean = false; // toggles between request/reset
  isResetMode = false;

  constructor(private route: ActivatedRoute, private authService: AuthServiceService, private router: Router) {}

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      const urlToken = params['token'];
      if (urlToken) {
        this.token = urlToken;       // save token
        this.isResetMode = true;     // show password input form
      }
    });
  }

  // Request password token
  onRecoverPassword() {
    if (!this.email) {
      alert('Por favor ingrese su correo electrónico');
      return;
    }

    // Save email locally so it's available after redirect

    this.authService.requestPasswordToken({ email: this.email }).subscribe({

      next: () => {
        // Save the email *after* the request succeeds
        localStorage.setItem('resetEmail', this.email);
        alert('Si el correo existe, se ha enviado un enlace de recuperación');
      },
      error: (err) => {
        console.error(err);
      }
    });
  }


  // Reset password using token
    onResetPassword() {
      if (this.newPassword !== this.confirmPassword) {
        alert('Las contraseñas no coinciden');
        return;
      }
      const savedEmail = localStorage.getItem('resetEmail');
      if (!savedEmail) {
        alert('No se encontró el correo asociado. Por favor repita el proceso de recuperación.');
        return;
      }
      const payload = {
        email: savedEmail,
        new_password: this.newPassword
      };

      // Print to check what’s being sent
      console.log('[ResetPassword] Payload being sent:', payload);
      console.log('[ResetPassword] Token being sent:', this.token);

      this.authService.resetPassword(payload, this.token).subscribe({
        next: (res) => {
          console.log('[ResetPassword] Server response:', res);
          alert('Contraseña cambiada exitosamente');
          localStorage.removeItem('resetEmail');
          this.isResetMode = false;
          this.router.navigate(['/login']);
        },
        error: (err) => {
          console.error('[ResetPassword] Error:', err);
          alert('Ocurrió un error. Intente nuevamente.');
        }
      });
    }


}
