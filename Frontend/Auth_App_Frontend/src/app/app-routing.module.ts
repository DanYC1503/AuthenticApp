import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponentComponent } from './LoginComponent/login-component/login-component.component';


const routes: Routes = [
  { path: 'login', component: LoginComponentComponent }, // define route
  { path: '', redirectTo: '/login', pathMatch: 'full' }, // optional: default redirect
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {

}
