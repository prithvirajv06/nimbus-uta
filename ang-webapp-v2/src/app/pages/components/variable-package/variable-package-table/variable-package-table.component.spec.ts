import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VariablePackageTableComponent } from './variable-package-table.component';

describe('VariablePackageTableComponent', () => {
  let component: VariablePackageTableComponent;
  let fixture: ComponentFixture<VariablePackageTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VariablePackageTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VariablePackageTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
