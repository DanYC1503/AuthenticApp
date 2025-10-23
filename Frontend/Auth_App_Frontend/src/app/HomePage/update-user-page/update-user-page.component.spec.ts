import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UpdateUserPageComponent } from './update-user-page.component';

describe('UpdateUserPageComponent', () => {
  let component: UpdateUserPageComponent;
  let fixture: ComponentFixture<UpdateUserPageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [UpdateUserPageComponent]
    });
    fixture = TestBed.createComponent(UpdateUserPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
