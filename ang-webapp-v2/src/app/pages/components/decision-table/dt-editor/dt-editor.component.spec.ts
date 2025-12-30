import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DtEditorComponent } from './dt-editor.component';

describe('DtEditorComponent', () => {
  let component: DtEditorComponent;
  let fixture: ComponentFixture<DtEditorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DtEditorComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DtEditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
