import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogicFlowLayoutComponent } from './logic-flow-layout.component';

describe('LogicFlowLayoutComponent', () => {
  let component: LogicFlowLayoutComponent;
  let fixture: ComponentFixture<LogicFlowLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LogicFlowLayoutComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LogicFlowLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
