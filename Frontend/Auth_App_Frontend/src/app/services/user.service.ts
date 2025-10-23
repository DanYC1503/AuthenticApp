import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private userData: any = null;

  constructor(private http: HttpClient) {}

  /** Fetch user info from API */
  getUserInfo(username: string): Observable<any> {
  return this.http.get(`${environment.API_GATEWAY_URL}users/info?username=${encodeURIComponent(username)}`, {
    withCredentials: true
  });
}

  /** Save user data after fetching */
  setUserData(data: any): void {
    this.userData = data;
  }

  /** Retrieve cached user data */
  getUserData(): any {
    return this.userData;
  }
}
