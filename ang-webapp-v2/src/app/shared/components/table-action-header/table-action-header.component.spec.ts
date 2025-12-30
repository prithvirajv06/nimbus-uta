import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableActionHeaderComponent } from './table-action-header.component';

describe('TableActionHeaderComponent', () => {
  let component: TableActionHeaderComponent;
  let fixture: ComponentFixture<TableActionHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TableActionHeaderComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TableActionHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
