import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LoopConfigureComponent } from './loop-configure.component';

describe('LoopConfigureComponent', () => {
  let component: LoopConfigureComponent;
  let fixture: ComponentFixture<LoopConfigureComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LoopConfigureComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LoopConfigureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
