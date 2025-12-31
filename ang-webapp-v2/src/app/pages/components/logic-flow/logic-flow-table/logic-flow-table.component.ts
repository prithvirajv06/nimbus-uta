import { Component, inject, OnInit, WritableSignal } from '@angular/core';
import { BasicTableComponent } from '../../../../shared/components/tables/basic-tables/basic-table.component';
import { DtService } from '../../../../shared/services/dt.service';
import { VariablePackageService } from '../../../../shared/services/variable-package.service';
import { VariablePackage, VariablePackageUtils } from '../../../../shared/types/variable_package';
import { LogicFlow, LogicFlowUtils } from '../../../../shared/types/logic-flow.type';
import { Option, SelectComponent } from '../../../../shared/components/form/select/select.component';
import { LogicFlowService } from '../../../../shared/services/logic-flow.service';
import { ApiResponse } from '../../../../shared/types/common.type';
import { NgClass, DatePipe } from '@angular/common';
import { Field } from '@angular/forms/signals';
import { PageBreadcrumbComponent } from '../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { InputFieldComponent } from '../../../../shared/components/form/input/input-field.component';
import { TextAreaComponent } from '../../../../shared/components/form/input/text-area.component';
import { LabelComponent } from '../../../../shared/components/form/label/label.component';
import { TableActionHeaderComponent } from '../../../../shared/components/table-action-header/table-action-header.component';
import { ButtonComponent } from '../../../../shared/components/ui/button/button.component';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';

@Component({
  selector: 'app-logic-flow-table',
  imports: [PageBreadcrumbComponent, NgClass, ButtonComponent,
    ModalComponent,
    DatePipe, TableActionHeaderComponent,
    LabelComponent, InputFieldComponent, TextAreaComponent, Field, SelectComponent],
  templateUrl: './logic-flow-table.component.html',
  styleUrl: './logic-flow-table.component.css',
})
export class LogicFlowTableComponent extends BasicTableComponent<LogicFlow> implements OnInit {

  variableService = inject(VariablePackageService);
  variablePackages: VariablePackage[] = [];
  variablePackagesOptions: Option[] = [];
  override setService(): void {
    this.service = inject(LogicFlowService);
  }

  formModule: WritableSignal<LogicFlow> = LogicFlowUtils.signalModel();
  formGroup = LogicFlowUtils.basicFormGroup(this.formModule);

  ngOnInit(): void {
    this.variableService.getList({ is_archived: false }).subscribe({
      next: (res) => {
        if (res.status === 'success' && res.data) {
          this.variablePackages = res.data;
          this.variablePackagesOptions = VariablePackageUtils.options(res.data);
        }
      }
    });
  }

  setVariablePackage(value: string) {
    const [nimb_id, versionStr] = value.split('~');
    const version = parseInt(versionStr, 10);
    const selectedVP = this.variablePackages.find(vp => vp.nimb_id === nimb_id && vp.audit.version === version);
    if (selectedVP) {
      this.formGroup.variable_package().setControlValue(selectedVP);
    }
  }

  onSubmit() {
    if (this.formGroup().valid()) {
      const newDT: LogicFlow = this.formGroup().value();
      this.service.create(newDT).subscribe({
        next: (res: ApiResponse<LogicFlow>) => {
          if (res.status === 'success') {
            this.notificationService.success('Decision Table created successfully!', 5);
            this.closeNewModal();
            this.loadtableView();
          }
        }
      });
    } else {
      this.formGroup().markAsDirty();
      this.notificationService.error('Please fill all required fields correctly.', 10);
    }
  }

}
