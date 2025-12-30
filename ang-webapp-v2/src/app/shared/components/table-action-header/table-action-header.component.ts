import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ButtonComponent } from '../ui/button/button.component';

@Component({
  selector: 'app-table-action-header',
  imports: [ButtonComponent],
  templateUrl: './table-action-header.component.html',
  styleUrl: './table-action-header.component.css',
})
export class TableActionHeaderComponent {

  @Input({ required: true }) title: string = '';
  @Input({ required: true }) showArchive: boolean = false;
  @Output() loadtableView = new EventEmitter<void>();
  @Output() openNewFormModal = new EventEmitter<void>();
  @Output() loadtableArchivedView = new EventEmitter<void>();

}
