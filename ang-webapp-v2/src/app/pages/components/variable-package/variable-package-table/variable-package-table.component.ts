import { NgClass, DatePipe } from "@angular/common";
import { Component, AfterContentInit, inject, signal } from "@angular/core";
import { Field, form, required } from "@angular/forms/signals";
import { PageBreadcrumbComponent } from "../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component";
import { FileInputComponent } from "../../../../shared/components/form/input/file-input.component";
import { InputFieldComponent } from "../../../../shared/components/form/input/input-field.component";
import { TextAreaComponent } from "../../../../shared/components/form/input/text-area.component";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { BasicTableComponent } from "../../../../shared/components/tables/basic-tables/basic-table.component";
import { ButtonComponent } from "../../../../shared/components/ui/button/button.component";
import { ModalComponent } from "../../../../shared/components/ui/modal/modal.component";
import { VariablePackageService } from "../../../../shared/services/variable-package.service";
import { VariablePackage, VariablePackageRequest } from "../../../../shared/types/variable_package";
import { TableActionHeaderComponent } from "../../../../shared/components/table-action-header/table-action-header.component";


@Component({
  selector: 'app-variable-package-table',
  imports: [PageBreadcrumbComponent, NgClass, ButtonComponent,
    ModalComponent, Field, InputFieldComponent, LabelComponent,
    TextAreaComponent, FileInputComponent, DatePipe, TableActionHeaderComponent],
  templateUrl: './variable-package-table.component.html',
  styleUrl: './variable-package-table.component.css',
})
export class VariablePackageTableComponent extends BasicTableComponent<VariablePackage> implements AfterContentInit {

  /**
   * Lifecycle hook that is called after the component's content has been fully initialized.
   */
  ngAfterContentInit(): void {
    this.loadtableView();
  }


  //Componetn Variables
  override transactionData: VariablePackage[] = [];
  override service = inject(VariablePackageService)

  formModal = signal<VariablePackageRequest>({
    package_name: '',
    description: '',
    json_str: '',
  });
  formGroup = form(this.formModal, schema => {
    required(schema.package_name, { message: 'Package Name is required' });
    required(schema.description, { message: 'Description is required' });
    required(schema.json_str, { message: 'JSON File is required' });
  });

  /**
   * Functional Implentation
   */

  onFileUpload(event: Event) {
    const element = event.currentTarget as HTMLInputElement;
    let fileList: FileList | null = element.files;

    if (fileList && fileList.length > 0) {
      const file = fileList[0];
      if (file.type !== 'application/json') {
        alert('Please upload a valid JSON file.');
        return;
      }
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const content = e.target?.result as string;
          const jsonData = JSON.parse(content);
          console.log('JSON Extracted:', jsonData);
          this.formGroup.json_str().setControlValue(JSON.stringify(jsonData));
        } catch (error) {
          alert('Error parsing JSON. File might be corrupted.');
        }
      };
      reader.readAsText(file);
    }
  }


  onSubmit() {
    if (this.formGroup().valid()) {
      const formData = this.formModal();
      this.service.create(formData).subscribe({
        next: (response) => {
          this.notificationService.success(response.message, 5);
          this.formGroup().reset();
          this.isNewModalOpen = false;
          this.loadtableView();
        }
      });
    } else {
      this.formGroup().markAsDirty();
    }
  }

  override setService(): void {
    this.service = inject(VariablePackageService);
  }
}