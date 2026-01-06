import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SimulationWindowComponent } from './simulation-window.component';

describe('SimulationWindowComponent', () => {
  let component: SimulationWindowComponent;
  let fixture: ComponentFixture<SimulationWindowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SimulationWindowComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SimulationWindowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
