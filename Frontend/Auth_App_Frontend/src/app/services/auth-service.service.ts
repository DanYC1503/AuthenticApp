import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { environment } from 'environments/environment';
import { BehaviorSubject } from 'rxjs';
import Swal from 'sweetalert2';

@Injectable({
  providedIn: 'root'
})
export class AuthServiceService {
  constructor(private http: HttpClient, private router: Router) { }

  private showAdminMenuSubject = new BehaviorSubject<boolean>(this.getStoredAdminState());
  public showAdminMenu$ = this.showAdminMenuSubject.asObservable();

   private getStoredAdminState(): boolean {
    // Check if admin state is stored in localStorage
    const isAdmin = localStorage.getItem('IS_ADMIN');
    return isAdmin === 'true';
  }
  setShowAdminMenu(show: boolean) {
      this.showAdminMenuSubject.next(show);
      // Save to localStorage to persist across page reloads
      localStorage.setItem('IS_ADMIN', show.toString());
    }

    getShowAdminMenu(): boolean {
      return this.showAdminMenuSubject.value;
    }

    // Call this on logout
    clearAdminState() {
      this.setShowAdminMenu(false);
      localStorage.removeItem('IS_ADMIN');
    }
  registerUser(payload: any): void {
    this.http.post(`${environment.API_GATEWAY_URL}auth/register`, payload, { withCredentials: true })
      .subscribe({
        next: (res) => {
          console.log('Registro exitoso:', res);

          // SweetAlert2 success modal
          Swal.fire({
            title: '🎉 Registro exitoso',
            text: 'Tu cuenta se creó correctamente. Serás redirigido al inicio de sesión.',
            icon: 'success',
            confirmButtonColor: '#002D74',
            timer: 3000,
            timerProgressBar: true,
            showConfirmButton: false,
            didOpen: () => {
              const swalContainer = document.querySelector('.swal2-container');
              if (swalContainer) swalContainer.classList.add('animate-pulse');
            }
          });

          // Wait 3 seconds before redirecting
          setTimeout(() => {
            this.router.navigate(['/login']);
          }, 3000);
        },
        error: (err) => {
          console.error('Error en registro:', err);
          Swal.fire({
            title: '❌ Error',
            text: 'Ocurrió un error durante el registro. Intenta nuevamente.',
            icon: 'error',
            confirmButtonColor: '#dc2626'
          });
        }
      });
  }
}
