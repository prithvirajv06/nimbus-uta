import { Component, EventEmitter, Host, HostListener, Input, OnChanges, Output, signal, SimpleChanges } from '@angular/core';
import { InputFieldComponent } from '../../../../shared/components/form/input/input-field.component';
import { CommonModule, NgClass } from '@angular/common';
import { Field, form, required } from '@angular/forms/signals';
import { LabelComponent } from '../../../../shared/components/form/label/label.component';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';
import { ArrayFilter, Variable } from '../../../../shared/types/variable_package';
import { RulesCommons } from '../../../../shared/util-fulctions/options';
import { ClickOutsideDirective } from '../../../../shared/directives/clickoutside';
import { CheckboxComponent } from "../../../../shared/components/form/input/checkbox.component";

@Component({
  selector: 'app-variable-selector',
  imports: [NgClass,
    CommonModule,
    LabelComponent,
    InputFieldComponent,
    ModalComponent,
    ClickOutsideDirective,
    Field, CheckboxComponent],
  templateUrl: './variable-selector.component.html',
  styleUrl: './variable-selector.component.css',
})
export class VariableSelectorComponent extends RulesCommons implements OnChanges {


  @Input() variables: Variable[] | any = [];
  @Input() selectedVariablePath: string | any = null;
  selectedVariable: Variable | null = null;
  @Output() changeInFilter: EventEmitter<ArrayFilter> = new EventEmitter<ArrayFilter>();
  @Output() variableSelectedEvent: EventEmitter<Variable> = new EventEmitter<Variable>();
  @Output() closeEvent: EventEmitter<void> = new EventEmitter<void>();


  filteredVars: Variable[] = [];
  @Input() isSelectBox: boolean = false;
  openChild: { [key: string]: boolean } = {};
  selectedFilterVariable: ArrayFilter = {
    array_name: '',
    property: '',
    logical: '',
    op_value: ''
  };
  selectedOperator: string | null = null;
  isVariableFilterOpen: boolean = false;

  constructor() {
    super();
    this.filteredVars = this.variables;
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.highlightSelectedVariable(this.selectedVariablePath);
  }

  highlightSelectedVariable(path: string): void {
    // if (path && this.variables.length > 0) {
    //   const foundVar = this.variables.find((v: Variable) => v.var_key === path);
    //   if (foundVar) {
    //     this.selectedVariable = foundVar;
    //     return;
    //   } else {
    //     path = path.substring(0, path.lastIndexOf('[*') != -1 ? path.lastIndexOf('[*') : path.lastIndexOf('.'));
    //     this.highlightSelectedVariable(path);
    //   }
    // }
    for (let variable of this.variables) {
      const result = this.searchVariableByPath(variable, path);
      if (result) {
        this.selectedVariable = result;
        return;
      }
    }
  }

  searchVariableByPath(variable: Variable, path: string): Variable | null {
    if (variable.var_key === path) {
      return variable;
    }
    if (variable.children && variable.children.length > 0) {
      for (let child of variable.children) {
        const result = this.searchVariableByPath(child, path);
        if (result) {
          return result;
        }
      }
    }
    return null;
  }

  toggleChild(variable: Variable) {
    this.openChild[variable.var_key] = !this.openChild[variable.var_key];
  }


  setSelectVariable(variable: Variable) {
    this.selectedVariable = variable;
    this.isVariableFilterOpen = false
    if (variable.isClickable) {
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
    value: '',
    operator: 'equals',
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
        logical: this.formModel().operator,
        op_value: this.formModel().value
      };
      this.changeInFilter.emit(this.selectedFilterVariable);
      this.reset();
    }
  }
  reset() {
    this.formModel.set({
      var_key: '',
      value: '',
      operator: 'equals',
    });
  }

  filterSelection(event: any) {
    const filterValue = event.target.value.toLowerCase();
    if (filterValue === '') {
      this.resetVariables();
    } else {
      this.filteredVars = this.getFilteredVariables(this.variables, filterValue);
    }
  }
  getFilteredVariables(variables: Variable[], filterValue: string): Variable[] {
    const filteredVars: Variable[] = [];
    for (const variable of variables) {
      if (variable.label.toLowerCase().includes(filterValue) || variable.var_key.toLowerCase().includes(filterValue)) {
        filteredVars.push(variable);
      } else if (variable.children && variable.children.length > 0) {
        const filteredChildren = this.getFilteredVariables(variable.children, filterValue);
        if (filteredChildren.length > 0) {
          filteredVars.push({
            ...variable,
            children: filteredChildren
          });
        }
      }
    }
    return filteredVars;
  }

  resetVariables() {
    // Logic to reset variables to original list can be added here
    this.filteredVars = this.variables;
  }
}
