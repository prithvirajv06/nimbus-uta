import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VariablePackageEditorComponent } from './variable-package-editor.component';

describe('VariablePackageEditorComponent', () => {
  let component: VariablePackageEditorComponent;
  let fixture: ComponentFixture<VariablePackageEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VariablePackageEditorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VariablePackageEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
