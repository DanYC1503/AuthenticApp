import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { ClientPageComponent } from './HomePage/client-page/client-page.component';
import { AdminPageComponent } from './HomePage/admin-page/admin-page.component';
import { RecoverPasswordComponent } from './InitPage/recover-password/recover-password.component';
import { CsrfInterceptor } from './interceptors/csrf.interceptor';
import { FormsModule } from '@angular/forms';
import { UsernameInterceptor } from './interceptors/username.interceptor';
import { LoginComponentComponent } from './InitPage/login-component/login-component.component';
import { RegistrarUsuarioComponent } from './InitPage/registrar-usuario/registrar-usuario.component';
import { UpdateUserPageComponent } from './HomePage/update-user-page/update-user-page.component';



@NgModule({
  declarations: [
    AppComponent,
    LoginComponentComponent,
    ClientPageComponent,
    AdminPageComponent,
    RecoverPasswordComponent,
    RegistrarUsuarioComponent,
    UpdateUserPageComponent

  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: CsrfInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: UsernameInterceptor, multi: true },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
