import { Component, EventEmitter, Input, OnChanges, Output, signal, SimpleChanges } from '@angular/core';
import { InputFieldComponent } from '../../../../shared/components/form/input/input-field.component';
import { NgClass } from '@angular/common';
import { Field, form, required } from '@angular/forms/signals';
import { LabelComponent } from '../../../../shared/components/form/label/label.component';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';
import { ArrayFilter, Variable } from '../../../../shared/types/variable_package';

@Component({
  selector: 'app-variable-selector',
  imports: [NgClass,
    LabelComponent,
    InputFieldComponent,
    ModalComponent,
    Field],
  templateUrl: './variable-selector.component.html',
  styleUrl: './variable-selector.component.css',
})
export class VariableSelectorComponent implements OnChanges {


  @Input() variables: Variable[] | any = [];
  @Input() selectedVariablePath: string | any = null;
  selectedVariable: Variable | null = null;
  @Output() changeInFilter: EventEmitter<ArrayFilter> = new EventEmitter<ArrayFilter>();
  @Output() variableSelectedEvent: EventEmitter<Variable> = new EventEmitter<Variable>();
  @Output() closeEvent: EventEmitter<void> = new EventEmitter<void>();
  selectedFilterVariable: ArrayFilter = {
    array_name: '',
    property: '',
    logical: '',
    op_value: ''
  };

  variableTypeOptions = [
    { value: 'string', label: 'String' },
    { value: 'number', label: 'Number' },
    { value: 'boolean', label: 'Boolean' },
    { value: 'object', label: 'Object' },
    { value: 'array', label: 'Array' }
  ];
  isVariableFilterOpen: boolean = false;

  ngOnChanges(changes: SimpleChanges): void {
    this.highlightSelectedVariable(this.selectedVariablePath);
  }

  highlightSelectedVariable(path: string): void {
    if (path && this.variables.length > 0) {
      const foundVar = this.variables.find((v: Variable) => v.var_key === path);
      if (foundVar) {
        this.selectedVariable = foundVar;
        return;
      } else {
        path = path.substring(0, path.lastIndexOf('[') != -1 ? path.lastIndexOf('[') : path.lastIndexOf('.'));
        this.highlightSelectedVariable(path);
      }
    }
  }


  setSelectVariable(variable: Variable) {
    this.selectedVariable = variable;
    if (variable.type != 'array' && variable.type != 'object') {
      this.variableSelectedEvent.emit(variable);
    }
  }

  isSelected(variable: Variable): boolean {
    return this.selectedVariable === variable;
  }


  closeModal() {
    this.isVariableFilterOpen = false;
  }

  formModel = signal({
    var_key: '',
    value: ''
  });

  formGroup = form(this.formModel, (schemaPath) => {
    required(schemaPath.var_key, { message: 'Variable Key is required' });
    required(schemaPath.value, { message: 'Value is required' });
  });

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

  saveFilter() {
    if (this.formGroup().valid()) {
      // Logic to save the filter can be added here
      this.isVariableFilterOpen = false;
      this.selectedFilterVariable = {
        array_name: this.selectedVariable?.var_key || '',
        property: this.formModel().var_key.split('.').pop() || '',
        logical: 'equals',
        op_value: this.formModel().value
      };
      this.changeInFilter.emit(this.selectedFilterVariable);
    }
  }
}
