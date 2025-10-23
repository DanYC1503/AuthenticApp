import { TestBed } from '@angular/core/testing';

import { UsernameInterceptor } from './username.interceptor';

describe('UsernameInterceptor', () => {
  beforeEach(() => TestBed.configureTestingModule({
    providers: [
      UsernameInterceptor
      ]
  }));

  it('should be created', () => {
    const interceptor: UsernameInterceptor = TestBed.inject(UsernameInterceptor);
    expect(interceptor).toBeTruthy();
  });
});
