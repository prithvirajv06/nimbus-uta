import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SideDrawComponent } from './side-draw.component';

describe('SideDrawComponent', () => {
  let component: SideDrawComponent;
  let fixture: ComponentFixture<SideDrawComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SideDrawComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SideDrawComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
