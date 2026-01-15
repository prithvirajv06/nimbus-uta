import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConditionConfigureComponent } from './condition-configure.component';

describe('ConditionConfigureComponent', () => {
  let component: ConditionConfigureComponent;
  let fixture: ComponentFixture<ConditionConfigureComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConditionConfigureComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ConditionConfigureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
