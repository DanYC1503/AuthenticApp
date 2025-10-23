import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'environments/environment';
import { map, tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class CsrfService {
  private baseUrl = environment.API_GATEWAY_URL;
  private csrfToken: string | null = null;

  constructor(private http: HttpClient) {}

  loadCsrfToken() {
    // Call backend to make sure cookie is set
    return this.http.get(`${this.baseUrl}api/csrf-token`, { withCredentials: true });
  }

  getTokenFromCookie(): string | null {
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : null;
  }
}
