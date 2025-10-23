import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class CsrfInterceptor implements HttpInterceptor {

  constructor() {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const token = localStorage.getItem('XSRF-TOKEN');

    //Clone request and add to header
    if(token){
      const cloned = request.clone({
        setHeaders:{
          'X-XSRF-TOKEN': token
        }
      });
      return next.handle(cloned)
    }
    return next.handle(request)

  }
}
