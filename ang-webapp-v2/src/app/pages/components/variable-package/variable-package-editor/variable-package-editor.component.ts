import { Component, OnInit, inject } from "@angular/core";
import { Field } from "@angular/forms/signals";
import { PageBreadcrumbComponent } from "../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component";
import { InputFieldComponent } from "../../../../shared/components/form/input/input-field.component";
import { TextAreaComponent } from "../../../../shared/components/form/input/text-area.component";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { ButtonComponent } from "../../../../shared/components/ui/button/button.component";
import { VariablePackageService } from "../../../../shared/services/variable-package.service";
import { Variable, VariablePackage, VariablePackageUtils } from "../../../../shared/types/variable_package";
import { MatExpansionModule } from "@angular/material/expansion";
import { MatIconModule } from "@angular/material/icon";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatDatepickerModule } from "@angular/material/datepicker";
import { CommonEditorComponent } from "../../../../shared/components/editor/common-editor.component";
import { VariableIteratorComponent } from "../variable-iterator/variable-iterator.component";

@Component({
  selector: 'app-variable-package-editor',
  imports: [PageBreadcrumbComponent, LabelComponent, InputFieldComponent, TextAreaComponent,
    ButtonComponent, Field,
    MatExpansionModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatDatepickerModule, VariableIteratorComponent],
  templateUrl: './variable-package-editor.component.html',
  styleUrl: './variable-package-editor.component.css',
})
export class VariablePackageEditorComponent extends CommonEditorComponent<VariablePackage> implements OnInit {


  variablePack!: VariablePackage | null;
  variablePackageService = inject(VariablePackageService);
  
  override setService(): void {
    this.service = this.variablePackageService;
  }

  override setFormModel(): void {
    this.formModel = VariablePackageUtils.signalModel();
  }
  override setFormGroup(): void {
    this.formGroup = VariablePackageUtils.detailsFormGroup(this.formModel);
  }

  variableTypeOptions = [
    { value: 'string', label: 'String' },
    { value: 'number', label: 'Number' },
    { value: 'boolean', label: 'Boolean' },
    { value: 'object', label: 'Object' },
    { value: 'array', label: 'Array' },
    { value: 'json', label: 'JSON' },
  ];

  addVariable() {
    const currentVariables = this.formModel().variables;
    let newVariable: Variable = {
      var_key: '',
      label: '',
      type: 'string',
      is_required: false,
      value: "",
      children: []
    };
    this.formModel.update(prev => ({
      ...prev,
      variables: [...currentVariables, newVariable]
    }));
  }


  removeVariable(index: number) {
    const currentVariables = this.formModel().variables;
    currentVariables.splice(index, 1);
    this.formModel.set({
      ...this.formModel(),
      variables: currentVariables
    });
  }

  alert(message: string) {
    window.alert(message);
  }
}
