import { VariablePackageUtils } from './../../../../shared/types/variable_package';
import { Component, Input } from '@angular/core';
import { Variable } from '../../../../shared/types/variable_package';
import { NgClass } from '@angular/common';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';
import { LabelComponent } from '../../../../shared/components/form/label/label.component';
import { InputFieldComponent } from '../../../../shared/components/form/input/input-field.component';
import { CheckboxComponent } from "../../../../shared/components/form/input/checkbox.component";
import { SelectComponent } from "../../../../shared/components/form/select/select.component";
import { Field } from "@angular/forms/signals";

@Component({
  selector: 'app-variable-iterator',
  imports: [NgClass,
    LabelComponent,
    InputFieldComponent,
    ModalComponent,
    CheckboxComponent,
    SelectComponent,
    Field],
  templateUrl: './variable-iterator.component.html',
  styleUrl: './variable-iterator.component.css',
})
export class VariableIteratorComponent {


  @Input() variables: Variable[] | any = [];
  @Input() selectedVariablePath: string = '';
  formGroup = VariablePackageUtils.newVariableFormGroup();

  variableTypeOptions = [
    { value: 'string', label: 'String' },
    { value: 'number', label: 'Number' },
    { value: 'boolean', label: 'Boolean' },
    { value: 'object', label: 'Object' },
    { value: 'array', label: 'Array' }
  ];
  isVariableEditorOpen: boolean = false;
  selectedVariable: Variable | null = null;

  setSelectVariable(variable: Variable) {
    this.selectedVariable = variable;
  }

  isSelected(variable: Variable): boolean {
    return this.selectedVariable === variable;
  }

  onAddVariable() {
    this.formGroup().setControlValue({
      var_key: this.selectedVariablePath ? this.selectedVariablePath + '.' : '' + 'new_variable_' + (this.variables.length + 1),
      label: '',
      type: 'string',
      is_required: false,
      value: '',
      children: []
    });
    this.isVariableEditorOpen = true;
  }

  closeNewModal() {
    this.isVariableEditorOpen = false;
  }

  saveVariable() {
    if (this.formGroup().valid()) {
      this.variables.push(this.formGroup().value());
      this.isVariableEditorOpen = false;
      this.formGroup().setControlValue({
        var_key: '',
        label: '',
        type: 'string',
        is_required: false,
        value: '',
        children: []
      });
    } else {
      this.formGroup().markAsTouched();
      this.formGroup().markAsDirty();
    }
  }

  getVarTypeIcon(variable: Variable): string {
    switch (variable.type) {
      case 'string':
        return 'fa fa-font';
      case 'number':
        return 'fa fa-hashtag';
      case 'boolean':
        return 'fa fa-toggle-on';
      case 'object':
        return 'fa fa-cube';
      case 'array':
        return 'fa fa-list';
      default:
        return 'fa fa-question';
    }
  }
}
