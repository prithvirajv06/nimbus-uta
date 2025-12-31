import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogicFlowTableComponent } from './logic-flow-table.component';

describe('LogicFlowTableComponent', () => {
  let component: LogicFlowTableComponent;
  let fixture: ComponentFixture<LogicFlowTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LogicFlowTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LogicFlowTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
