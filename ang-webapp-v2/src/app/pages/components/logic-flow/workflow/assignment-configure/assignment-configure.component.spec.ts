import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AssignmentConfigureComponent } from './assignment-configure.component';

describe('AssignmentConfigureComponent', () => {
  let component: AssignmentConfigureComponent;
  let fixture: ComponentFixture<AssignmentConfigureComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AssignmentConfigureComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AssignmentConfigureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
