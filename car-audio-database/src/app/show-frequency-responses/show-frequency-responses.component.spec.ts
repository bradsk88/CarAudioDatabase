import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowFrequencyResponsesComponent } from './show-frequency-responses.component';

describe('ShowFrequencyResponsesComponent', () => {
  let component: ShowFrequencyResponsesComponent;
  let fixture: ComponentFixture<ShowFrequencyResponsesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ShowFrequencyResponsesComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ShowFrequencyResponsesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
