import { Component } from '@angular/core';
import { AuthServiceService } from 'src/app/services/auth-service.service';

@Component({
  selector: 'app-registrar-usuario',
  templateUrl: './registrar-usuario.component.html',
  styleUrls: ['./registrar-usuario.component.css']
})
export class RegistrarUsuarioComponent {
  showPassword: boolean = false;

  // Form fields
  id_number: string = '';
  full_name: string = '';
  email: string = '';
  password: string = '';
  phone_number: string = '';
  date_of_birth: string = '';
  address: string = '';
  username: string = '';

  constructor(private authService: AuthServiceService) {}

  onRegister() {
     if (this.username === '' || this.email === '' || this.password === '') {
      alert("Llene todos los datos porfavor");
      return;
    }
    const payload = {
      id_number: this.id_number,
      full_name: this.full_name,
      email: this.email,
      password: this.password,
      phone_number: this.phone_number,
      date_of_birth: this.date_of_birth,
      address: this.address,
      username: this.username
    };


    this.authService.registerUser(payload);
  }
}
