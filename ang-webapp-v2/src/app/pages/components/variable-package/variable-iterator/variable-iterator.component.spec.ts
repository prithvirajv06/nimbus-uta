import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VariableIteratorComponent } from './variable-iterator.component';

describe('VariableIteratorComponent', () => {
  let component: VariableIteratorComponent;
  let fixture: ComponentFixture<VariableIteratorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VariableIteratorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VariableIteratorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
