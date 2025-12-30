import { CommonModule } from '@angular/common';
import { Component, EventEmitter, inject, Output } from '@angular/core';
import { ButtonComponent } from '../../ui/button/button.component';
import { TableDropdownComponent } from '../../common/table-dropdown/table-dropdown.component';
import { BadgeComponent } from '../../ui/badge/badge.component';
import { NotificationService } from '../../../services/notification.service';
import { VariablePackage } from '../../../types/variable_package';
import { ManagementContract } from '../../../contract/management.contarct';
import { ApiResponse } from '../../../types/common.type';

interface Transaction {
  image: string;
  action: string;
  date: string;
  amount: string;
  category: string;
  status: "Success" | "Pending" | "Failed";
}

@Component({
  selector: 'app-basic-table-three',
  imports: [
    CommonModule
],
  templateUrl: './basic-table.component.html',
  styles: ``
})
export class BasicTableComponent<T extends { audit: { version: number; [key: string]: any }, nimb_id: string | number }> implements ManagementContract {

  // Type definition for the transaction data

  showArchive = false;
  service: any;
  transactionData: T[] = []
  notificationService = inject(NotificationService)
  @Output() edit = new EventEmitter<T>();
  isNewModalOpen = false;
  currentPage = 1;
  itemsPerPage = 5;

  constructor() {
    this.setService();
    this.loadtableView();
  }

  get totalPages(): number {
    return Math.ceil(this.transactionData.length / this.itemsPerPage);
  }

  get currentItems(): T[] {
    const start = (this.currentPage - 1) * this.itemsPerPage;
    return this.transactionData.slice(start, start + this.itemsPerPage);
  }

  goToPage(page: number) {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
    }
  }

  handleViewMore(item: Transaction) {
    // logic here
    console.log('View More:', item);
  }

  handleDelete(item: Transaction) {
    // logic here
    console.log('Delete:', item);
  }

  getBadgeColor(status: string): 'success' | 'warning' | 'error' {
    if (status === 'Success') return 'success';
    if (status === 'Pending') return 'warning';
    return 'error';
  }


  closeNewModal() {
    this.isNewModalOpen = false;
  }



  loadtableView() {
    this.service.getList({ is_archived: false }).subscribe({
      next: (response: ApiResponse<T[]>) => {
        this.transactionData = Array.isArray(response.data) ? response.data as T[] : [];
      }
    });
    this.showArchive = false;
  }

  editDetails(item: T) {
    this.edit.emit(item);
  }

  archive(data: T) {
    this.service.delete(data.nimb_id, data.audit.version).subscribe({
      next: (response: ApiResponse<T>) => {
        this.notificationService.success('Variable Package archived successfully.', 5);
        this.loadtableView();
      }
    });
  }

  loadtableArchivedView() {
    this.service.getList({ is_archived: true }).subscribe({
      next: (response: ApiResponse<T[]>) => {
        this.transactionData = Array.isArray(response.data) ? response.data as T[] : [];
        this.showArchive = true;
      }
    });
    this.showArchive = true;
  }

  restore(data: T) {
    const updatedData = { ...data, audit: { ...data.audit, is_archived: false, restore_archive: true } };
    this.service.update(data.nimb_id, data.audit.version, updatedData).subscribe({
      next: () => {
        this.notificationService.success('Variable Package restored successfully.', 5);
        this.loadtableArchivedView();
      }
    });
  }

  setService(): void {
    throw new Error('Method not implemented.');
  }
}
