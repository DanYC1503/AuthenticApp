import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';


@Injectable()
export class UsernameInterceptor implements HttpInterceptor {
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let body = req.body;

    // Only modify if body exists and doesn't have username
    if (body && !body.username) {
      const defaultUsername = localStorage.getItem('USERNAME') || 'defaultUser';
      body = { ...body, username: defaultUsername };
    }

    const cloned = req.clone({ body });
    return next.handle(cloned);
  }
}
