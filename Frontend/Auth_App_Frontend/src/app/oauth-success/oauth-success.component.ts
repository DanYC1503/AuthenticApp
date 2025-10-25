import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-oauth-success',
  template: '<p>Logging in...</p>',
})
export class OauthSuccessComponent implements OnInit {
  constructor(private route: ActivatedRoute, private router: Router) {}

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      const token = params['token'];
      if (token) {
        localStorage.setItem('SESSION_TOKEN', token);
        // Optionally set up user info here
        this.router.navigate(['/dashboard']); // or wherever
      } else {
        console.error('OAuth token missing!');
        this.router.navigate(['/login']);
      }
    });
  }
}
