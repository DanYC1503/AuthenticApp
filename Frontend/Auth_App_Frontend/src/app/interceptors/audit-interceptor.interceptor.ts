import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor
} from '@angular/common/http';
import { Observable } from 'rxjs';

  @Injectable()
  export class AuditInterceptorInterceptor implements HttpInterceptor {

    constructor() {}

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
      const username = localStorage.getItem('USERNAME') || 'anonymous';

      const metadata = {
        browser: navigator.userAgent,
        language: navigator.language,
        platform: navigator.platform,
        timestamp: new Date().toISOString()
      };

      const cloned = req.clone({
        setHeaders: {
          'X-Username': username,
          'X-Metadata': JSON.stringify(metadata)
        }
      });


      return next.handle(cloned);
    }
  }
