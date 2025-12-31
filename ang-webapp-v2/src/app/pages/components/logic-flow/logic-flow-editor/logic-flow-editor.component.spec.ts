import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogicFlowEditorComponent } from './logic-flow-editor.component';

describe('LogicFlowEditorComponent', () => {
  let component: LogicFlowEditorComponent;
  let fixture: ComponentFixture<LogicFlowEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LogicFlowEditorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LogicFlowEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
