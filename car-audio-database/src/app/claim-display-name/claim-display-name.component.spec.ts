import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClaimDisplayNameComponent } from './claim-display-name.component';

describe('ClaimDisplayNameComponent', () => {
  let component: ClaimDisplayNameComponent;
  let fixture: ComponentFixture<ClaimDisplayNameComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ClaimDisplayNameComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ClaimDisplayNameComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
