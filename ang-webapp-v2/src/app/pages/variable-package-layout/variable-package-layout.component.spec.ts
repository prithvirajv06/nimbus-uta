import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VariablePackageLayoutComponent } from './variable-package-layout.component';

describe('VariablePackageComponent', () => {
  let component: VariablePackageLayoutComponent;
  let fixture: ComponentFixture<VariablePackageLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VariablePackageLayoutComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VariablePackageLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
