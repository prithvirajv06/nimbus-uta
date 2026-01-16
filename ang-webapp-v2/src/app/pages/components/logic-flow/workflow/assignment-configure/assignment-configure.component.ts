import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import { WorkflowStep } from '../../../../../shared/types/workflow-step';
import { Variable } from '../../../../../shared/types/variable_package';
import { VariableSelectorComponent } from "../../../variable-package/variable-selector/variable-selector.component";
import { LabelComponent } from "../../../../../shared/components/form/label/label.component";
import { TextAreaComponent } from "../../../../../shared/components/form/input/text-area.component";
import { InputFieldComponent } from '../../../../../shared/components/form/input/input-field.component';
import { single } from 'rxjs';
import { Field, form, required } from "@angular/forms/signals";
import { JsonPipe } from '@angular/common';
import { SwitchComponent } from "../../../../../shared/components/form/input/switch.component";
import { ButtonComponent } from '../../../../../shared/components/ui/button/button.component';

@Component({
  selector: 'app-assignment-configure',
  imports: [VariableSelectorComponent, LabelComponent, TextAreaComponent, InputFieldComponent, Field,
    JsonPipe, SwitchComponent, ButtonComponent],
  templateUrl: './assignment-configure.component.html',
  styleUrl: './assignment-configure.component.css',
})
export class AssignmentConfigureComponent {

  @Input() assignmentNode: WorkflowStep = {} as WorkflowStep;
  @Input() variables: Variable[] = [];
  @Input() loopNode: WorkflowStep = {} as WorkflowStep;
  @Input() parentVar: Variable = {} as Variable;
  applicableVariables: Variable[] = [];
  selectedVar: Variable = {} as Variable;
  @Output() onSave = new EventEmitter<WorkflowStep>();
  @Output() onCancel = new EventEmitter<void>();
  @Output() onDelete = new EventEmitter<WorkflowStep>();
  variableAssignmentForm = signal({
    target_variable: this.assignmentNode.target || {} as Variable,
    value_to_assign: this.assignmentNode.value || ''
  });
  formGroup = form(this.variableAssignmentForm, (schema) => {
    required(schema.target_variable, { message: 'Target Variable is required' });
    required(schema.value_to_assign, { message: 'Value to assign is required' });
  });
  constructor() { }


  ngAfterContentInit(): void {
    this.applicableVariables = this.getApplicableTypeVars(this.variables, ['string', 'number', 'boolean', 'object'], false);
    this.formGroup().setControlValue({
      target_variable: this.assignmentNode.target || {} as Variable,
      value_to_assign: this.assignmentNode.value || ''
    });
  }

  getApplicableTypeVars(variable: Variable[], allowedType: string[], isCrossedArray: boolean): Variable[] {
    for (let varItem of variable) {
      if (varItem.children && varItem.children.length > 0) {
        varItem.children = this.getApplicableTypeVars(varItem.children, allowedType, isCrossedArray || varItem.type === 'array');
        varItem.isClickable = varItem.children.every(child => child.isClickable) && allowedType.includes(varItem.type) && !isCrossedArray;
      }
      if ((allowedType.includes(varItem.type) &&  !isCrossedArray)
        || (allowedType.includes(varItem.type) && this.parentVar && varItem.var_key.includes(this.parentVar.var_key))) {
        varItem.isClickable = true;
      }
    }
    return variable;
  }

  saveConfiguration() {
    this.assignmentNode.target = this.variableAssignmentForm().target_variable;
    this.assignmentNode.value = this.variableAssignmentForm().value_to_assign;
    this.assignmentNode.statement_label = `Assign ${this.assignmentNode.value} to ${this.assignmentNode.target.label}`;
    this.onSave.emit(this.assignmentNode);
  }


}
