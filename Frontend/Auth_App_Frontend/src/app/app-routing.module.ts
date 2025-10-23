import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponentComponent } from './InitPage/login-component/login-component.component';
import { ClientPageComponent } from './HomePage/client-page/client-page.component';
import { AdminPageComponent } from './HomePage/admin-page/admin-page.component';
import { RecoverPasswordComponent } from './InitPage/recover-password/recover-password.component';
import { RegistrarUsuarioComponent } from './InitPage/registrar-usuario/registrar-usuario.component';
import { UpdateUserPageComponent } from './HomePage/update-user-page/update-user-page.component';


const routes: Routes = [
  { path: 'login', component: LoginComponentComponent }, // define route
  { path: '', redirectTo: '/login', pathMatch: 'full' }, // optional: default redirect
  { path: 'ClientDashboard', component: ClientPageComponent }, // optional: default redirect
  { path: 'AdminDashboard', component: AdminPageComponent }, // optional: default redirect
  { path: 'passwordRecovery', component: RecoverPasswordComponent },
  { path: 'register/user', component: RegistrarUsuarioComponent },
  { path: 'update/User', component: UpdateUserPageComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
