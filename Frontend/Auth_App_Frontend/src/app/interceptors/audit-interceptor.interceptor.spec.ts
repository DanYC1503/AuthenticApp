import { TestBed } from '@angular/core/testing';

import { AuditInterceptorInterceptor } from './audit-interceptor.interceptor';

describe('AuditInterceptorInterceptor', () => {
  beforeEach(() => TestBed.configureTestingModule({
    providers: [
      AuditInterceptorInterceptor
      ]
  }));

  it('should be created', () => {
    const interceptor: AuditInterceptorInterceptor = TestBed.inject(AuditInterceptorInterceptor);
    expect(interceptor).toBeTruthy();
  });
});
