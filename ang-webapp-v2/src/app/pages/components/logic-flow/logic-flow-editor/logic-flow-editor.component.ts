import { Component, inject, OnInit } from '@angular/core';
import { PageBreadcrumbComponent } from '../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { CommonEditorComponent } from '../../../../shared/components/editor/common-editor.component';
import { DecisionTable, DTUtils } from '../../../../shared/types/dt.type';
import { Condition, LogicalStep, LogicFlow, LogicFlowUtils, Operation } from '../../../../shared/types/logic-flow.type';
import { LogicFlowService } from '../../../../shared/services/logic-flow.service';
import { VariablePackageService } from '../../../../shared/services/variable-package.service';
import { Variable, VariablePackage } from '../../../../shared/types/variable_package';
import { Option, SelectComponent } from '../../../../shared/components/form/select/select.component';
import { Field } from '@angular/forms/signals';
import { MatExpansionModule } from '@angular/material/expansion';
import { InputFieldComponent } from '../../../../shared/components/form/input/input-field.component';
import { TextAreaComponent } from '../../../../shared/components/form/input/text-area.component';
import { LabelComponent } from '../../../../shared/components/form/label/label.component';
import { ButtonComponent } from '../../../../shared/components/ui/button/button.component';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';
import { CommonModule, JsonPipe } from '@angular/common';

@Component({
  selector: 'app-logic-flow-editor',
  imports: [PageBreadcrumbComponent, MatExpansionModule, LabelComponent, InputFieldComponent,
    JsonPipe,
    TextAreaComponent, Field, ButtonComponent, ModalComponent, SelectComponent, CommonModule],
  templateUrl: './logic-flow-editor.component.html',
  styleUrl: './logic-flow-editor.component.css',
})
export class LogicFlowEditorComponent extends CommonEditorComponent<LogicFlow> implements OnInit {

  dtService = inject(LogicFlowService);
  varService = inject(VariablePackageService);
  variablePackage: VariablePackage | null = null;
  variableOptions: Option[] = [];

  logicalOptions: Option[] = [
    { value: 'AND', label: 'AND' },
    { value: 'OR', label: 'OR' },
    { value: 'NOT', label: 'NOT' }
  ];

  logicalOperators: Option[] = [
    { label: "Equal (=)", value: "eq" },
    { label: "Not Equal (≠)", value: "neq" },
    { label: "Greater Than (>)", value: "gt" },
    { label: "Less Than (<)", value: "lt" },
    { label: "Greater Than or Equal (≥)", value: "gte" },
    { label: "Less Than or Equal (≤)", value: "lte" },
    { label: "Contains", value: "contains" },
    { label: "Starts With", value: "startswith" },
    { label: "Ends With", value: "endswith" },
    { label: "Matches (Regex)", value: "matches" },
    { label: "In Array", value: "inarray" },
    { label: "Not In Array", value: "notinarray" },
    { label: "Is Empty", value: "empty" },
  ];
  actionOperators: Option[] = [
    { label: "Set", value: "SET" },
    { label: "Add (+)", value: "ADD" },
    { label: "Multiply (×)", value: "MULT" },
    { label: "Subtract (−)", value: "SUB" },
    { label: "Set Object", value: "SET_OBJ" },
    { label: "Collect", value: "COLLECT" },
    { label: "Collect Sum", value: "COLLECT_SUM" },
    { label: "Collect Count", value: "COLLECT_COUNT" },
    { label: "Delete", value: "DELETE" },
    { label: "Push", value: "PUSH" },
    { label: "Remove", value: "REMOVE" },
    { label: "Clear", value: "CLEAR" },
    { label: "Uppercase", value: "UPPERCASE" },
    { label: "Lowercase", value: "LOWERCASE" },
    { label: "Trim", value: "TRIM" },
    { label: "Append", value: "APPEND" },
    { label: "Prepend", value: "PREPEND" },
    { label: "Increment (++)", value: "INCREMENT" },
    { label: "Decrement (−−)", value: "DECREMENT" },
    { label: "Toggle", value: "TOGGLE" },
    { label: "Reverse", value: "REVERSE" },
    { label: "Sort Ascending", value: "SORT_ASC" },
    { label: "Sort Descending", value: "SORT_DESC" }
  ];

  override setService(): void {
    this.service = this.dtService
  }

  override setFormModel(): void {
    this.formModel = LogicFlowUtils.signalModel();
  }
  override setFormGroup(): void {
    this.formGroup = LogicFlowUtils.detailsFormGroup(this.formModel);
  }

  override afterGetDetails(): void {
    const nimb_id = this.formModel().variable_package.nimb_id;
    const version = this.formModel().variable_package.audit.version;
    if (version && nimb_id) {
      this.varService.get(nimb_id, version).subscribe({
        next: (res) => {
          if (res.status === 'success' && res.data) {
            this.variablePackage = res.data;
            this.variableOptions = <Option[]>res.data.variables.map((v) => ({
              label: v.label + " (" + v.type + ")",
              value: v.var_key
            }));
          }
        }
      });
    }
  }


  addLogicStep() {
    const currentSteps = this.formModel().logical_steps;
    const newStep: LogicalStep = {
      operation_name: '',
      condition: {
        operator: '',
        conditions: [],
        variable: {} as Variable,
        logical: '',
        op_value: ''
      },
      variable: {} as Variable,
      logical: '',
      op_value: '',
      operation_if_true: [],
      operation_if_false: []
    }
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: [...currentSteps, newStep]
    }));
  }

  addCondition(stepIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    const newCondition: Condition = {
      operator: '',
      conditions: [],
      variable: {} as Variable,
      logical: '',
      op_value: ''
    };
    // Create a new array for conditions to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        return {
          ...step,
          condition: {
            ...step.condition,
            conditions: [...step.condition.conditions, newCondition]
          }
        };
      }
      return step;
    });
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
    console.log(this.formGroup().value());
  }

  addGroupCondition(stepIndex: number, conditionIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        const updatedConditions = step.condition.conditions.map((cond, cIdx) => {
          if (cIdx === conditionIndex) {
            return {
              ...cond,
              conditions: [
                ...(cond.conditions || []),
                {
                  operator: '',
                  conditions: [],
                  variable: {
                    var_key: '',
                    label: '',
                    type: '',
                    is_required: false,
                    value: ''
                  },
                  logical: '',
                  op_value: ''
                } as Condition
              ]
            };
          }
          return cond;
        });
        return {
          ...step,
          condition: {
            ...step.condition,
            conditions: updatedConditions
          }
        };
      }
      return step;
    });
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  removeGroupCondition(stepIndex: number, conditionIndex: number, groupConditionIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        const updatedConditions = step.condition.conditions.map((cond, cIdx) => {
          if (cIdx === conditionIndex) {
            const newGroupConditions = cond.conditions?.filter((_, gIdx) => gIdx !== groupConditionIndex) || [];
            return {
              ...cond,
              conditions: newGroupConditions
            };
          }
          return cond;
        });
        return {
          ...step,
          condition: {
            ...step.condition,
            conditions: updatedConditions
          }
        };
      }
      return step;
    });
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  removeCondition(stepIndex: number, conditionIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        const newConditions = step.condition.conditions.filter((_, cIdx) => cIdx !== conditionIndex);
        return {
          ...step,
          condition: {
            ...step.condition,
            conditions: newConditions
          }
        };
      }
      return step;
    });
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  removeIfTrueAction(stepIndex: number, actionIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        const newActions = step.operation_if_true.filter((_, aIdx) => aIdx !== actionIndex);
        return {
          ...step,
          operation_if_true: newActions
        };
      }
      return step;
    }
    );
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  addIfTrueAction(stepIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    const newAction: Operation = {
      variable: {
        var_key: '',
        label: '',
        type: '',
        is_required: false,
        value: ''
      },
      operation: '',
      op_value: '',
      value_is_path: false
    };
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        return {
          ...step,
          operation_if_true: [...step.operation_if_true, newAction]
        };
      }
      return step;
    }
    );
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  removeIfFalseAction(stepIndex: number, actionIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        const newActions = step.operation_if_false.filter((_, aIdx) => aIdx !== actionIndex);
        return {
          ...step,
          operation_if_false: newActions
        };
      }
      return step;
    }
    );
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }
  addIfFalseAction(stepIndex: number) {
    const currentSteps = this.formModel().logical_steps;
    const newAction: Operation = {
      variable: {
        var_key: '',
        label: '',
        type: '',
        is_required: false,
        value: ''
      },
      operation: '',
      op_value: '',
      value_is_path: false
    };
    // Deep copy the steps to avoid direct mutation
    const updatedSteps = currentSteps.map((step, idx) => {
      if (idx === stepIndex) {
        return {
          ...step,
          operation_if_false: [...step.operation_if_false, newAction]
        };
      }
      return step;
    }
    );
    this.formModel.update(prev => ({
      ...prev,
      logical_steps: updatedSteps
    }));
  }

  setVarVaule(varField: any, selectedVarKey: string) {
    if (this.variablePackage) {
      const selectedVar = this.variablePackage.variables.find(v => v.var_key === selectedVarKey);
      if (selectedVar) {
        varField.setControlValue(selectedVar);
      }
    }
  }
}

