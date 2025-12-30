import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DecisionTableLayoutComponent } from './decision-table-layout.component';

describe('DecisionTableLayoutComponent', () => {
  let component: DecisionTableLayoutComponent;
  let fixture: ComponentFixture<DecisionTableLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DecisionTableLayoutComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DecisionTableLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
