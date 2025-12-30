import { Component, OnInit, inject, signal } from "@angular/core";
import { Field, form, required, applyEach } from "@angular/forms/signals";
import { Router, ActivatedRoute } from "@angular/router";
import { PageBreadcrumbComponent } from "../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component";
import { CheckboxComponent } from "../../../../shared/components/form/input/checkbox.component";
import { InputFieldComponent } from "../../../../shared/components/form/input/input-field.component";
import { TextAreaComponent } from "../../../../shared/components/form/input/text-area.component";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { SelectComponent } from "../../../../shared/components/form/select/select.component";
import { ButtonComponent } from "../../../../shared/components/ui/button/button.component";
import { NotificationService } from "../../../../shared/services/notification.service";
import { VariablePackageService } from "../../../../shared/services/variable-package.service";
import { ApiResponse } from "../../../../shared/types/common.type";
import { VariablePackage, Variable } from "../../../../shared/types/variable_package";
import { MatExpansionModule } from "@angular/material/expansion";
import { MatIconModule } from "@angular/material/icon";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatDatepickerModule } from "@angular/material/datepicker";

@Component({
  selector: 'app-variable-package-editor',
  imports: [PageBreadcrumbComponent, LabelComponent, InputFieldComponent, TextAreaComponent,
    ButtonComponent, CheckboxComponent, Field, SelectComponent,
    MatExpansionModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatDatepickerModule,],
  templateUrl: './variable-package-editor.component.html',
  styleUrl: './variable-package-editor.component.css',
})
export class VariablePackageEditorComponent implements OnInit {


  ngOnInit(): void {
    this.getVariablepackage();
  }

  variablePack!: VariablePackage | null;
  variablePackageService = inject(VariablePackageService);
  router = inject(Router);
  formModel = signal<VariablePackage>({
    nimb_id: '',
    package_name: '',
    description: '',
    variables: [] as Variable[],
    audit: {
      version: 0,
      created_at: '',
      modified_at: '',
      created_by: '',
      modified_by: '',
      status:'DRAFT',
      is_prod_candidate: false,
      active: '',
      is_archived: false,
      minor_version: 0,
      restore_archive: false
    },
    no_of_variables: 0
  });
  activatedRoute = inject(ActivatedRoute);
  notificationService = inject(NotificationService);
  formGroup = form(this.formModel, (schemaPath) => {
    required(schemaPath.package_name, { message: 'Package Name is required' });
    required(schemaPath.description, { message: 'Description is required' });
    applyEach(schemaPath.variables, (variableSchema) => {
      required(variableSchema.label, { message: 'Variable Name is required' });
      required(variableSchema.type, { message: 'Variable Type is required' });
      required(variableSchema.var_key, { message: 'Variable Value is required' });
    });
  });

  variableTypeOptions = [
    { value: 'string', label: 'String' },
    { value: 'number', label: 'Number' },
    { value: 'boolean', label: 'Boolean' },
    { value: 'object', label: 'Object' },
    { value: 'array', label: 'Array' },
    { value: 'json', label: 'JSON' },
  ];


  getVariablepackage() {
    this.activatedRoute.queryParams.subscribe(params => {
      const editId = params['nimb_id'];
      const editVersion = params['version'];
      if (editId && editVersion) {
        this.variablePackageService.get(editId, editVersion)
          .subscribe((response: ApiResponse<VariablePackage>) => {
            this.variablePack = response.data;
            this.formModel.set(this.variablePack);
          });
      }
    });
  }

  cancelEdit() {
    this.formGroup().reset();
    this.variablePack = null;
    this.router.navigate([], {
      queryParams: {}
    });
  }

  saveVariablePackage() {
    if (this.formGroup().valid()) {
      const formValue = this.formGroup().value();
      // Update existing variable package
      this.variablePackageService.update(formValue.nimb_id, formValue.audit.version, formValue)
        .subscribe((response: ApiResponse<VariablePackage>) => {
          this.notificationService.success('Variable Package updated successfully.', 5);
        });
    } else {
      this.notificationService.error('Please fill in all required fields.', 10);
    }
  }

  removeVariable(index: number) {
    const currentVariables = this.formModel().variables;
    currentVariables.splice(index, 1);
    this.formModel.set({
      ...this.formModel(),
      variables: currentVariables
    });
  }
}
