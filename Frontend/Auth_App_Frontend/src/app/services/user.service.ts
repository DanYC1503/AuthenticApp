import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private userData: any = null;
  private baseUrl = `${environment.API_GATEWAY_URL}users/`;

  constructor(private http: HttpClient) {}

  /** Fetch user info from API */
  getUserInfo(username: string): Observable<any> {
  return this.http.get(`${environment.API_GATEWAY_URL}users/info?username=${encodeURIComponent(username)}`, {
    withCredentials: true
  });
  }
  /** Update user info */
  updateUserInfo(payload: any, updateToken: string): Observable<any> {
  return this.http.put(
    `${this.baseUrl}update`,
    payload,
    {
      headers: { 'X-Update-Auth': updateToken },
      withCredentials: true
    }
  );
  }
  deleteUser(payload: any, deleteToken: string): Observable<any> {
  return this.http.delete(
    `${this.baseUrl}delete`,
    {
      headers: { 'X-Delete-Auth': deleteToken },
      body: payload,          // <--- include the payload here
      withCredentials: true
    }
  );
  }

  /** Save user data after fetching */
  setUserData(data: any): void {
    this.userData = data;
  }

  /** Retrieve cached user data */
  getUserData(): any {
    return this.userData;
  }

  listUsers(email: string): Observable<any> {
    return this.http.get(`${environment.API_GATEWAY_URL}users/list/users?email=${encodeURIComponent(email)}`, {
      withCredentials: true
    });
  }
  // In your user.service.ts
  disableUser(username: string, clientUsername: string): Observable<any> {
    const payload = {
      username: username,
      client_username: clientUsername
    };

    return this.http.put(`${environment.API_GATEWAY_URL}users/disable/user`, payload, {
      withCredentials: true
    });
  }

  enableUser(username: string, clientUsername: string): Observable<any> {
    const payload = {
      username: username,
      client_username: clientUsername
    };

    return this.http.put(`${environment.API_GATEWAY_URL}users/enable/user`, payload, {
      withCredentials: true
    });
  }
}

