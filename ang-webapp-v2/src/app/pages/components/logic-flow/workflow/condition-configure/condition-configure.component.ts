import { Component, EventEmitter, inject, Input, OnInit, Output, signal, WritableSignal } from '@angular/core';
import { LabelComponent } from '../../../../../shared/components/form/label/label.component';
import { InputFieldComponent } from '../../../../../shared/components/form/input/input-field.component';
import { ButtonComponent } from '../../../../../shared/components/ui/button/button.component';
import { ConditionConfig, WorkflowStep } from '../../../../../shared/types/workflow-step';
import { applyEach, form, required, Field } from '@angular/forms/signals';
import { SelectComponent } from "../../../../../shared/components/form/select/select.component";
import { ToggleSwitchComponent } from "../../../../../shared/components/form/form-elements/toggle-switch/toggle-switch.component";
import { SwitchComponent } from "../../../../../shared/components/form/input/switch.component";
import { JsonPipe } from '@angular/common';
import { VariableSelectorComponent } from "../../../variable-package/variable-selector/variable-selector.component";
import { Variable } from '../../../../../shared/types/variable_package';
import { NotificationService } from '../../../../../shared/services/notification.service';
import { SearchSelectComponent } from "../../../../../shared/components/ui/search-select/search-select.component";
import { ClickOutsideDirective } from '../../../../../shared/directives/clickoutside';

@Component({
  selector: 'app-condition-configure',
  imports: [LabelComponent, InputFieldComponent, ButtonComponent, SwitchComponent, VariableSelectorComponent, Field, SearchSelectComponent],
  templateUrl: './condition-configure.component.html',
  styleUrl: './condition-configure.component.css',
})
export class ConditionConfigureComponent implements OnInit {

  @Input({ required: true }) conditionNode: WorkflowStep | null = null;
  @Input({ required: true }) variables: Variable[] = [];
  @Input({ required: true }) parentVar: Variable = {} as Variable;
  @Output() onSave = new EventEmitter<WorkflowStep>();
  @Output() onCancel = new EventEmitter<void>();

  notificationService = inject(NotificationService);

  workflowStep: WritableSignal<WorkflowStep> = signal(this.conditionNode as WorkflowStep);

  formGroup = form(this.workflowStep, (schema) => {
    applyEach(schema.condition_config, (config) => {
      required(config.left_var, { message: 'Left Variable is required' });
      required(config.operator, { message: 'Operator is required' });
      required(config.right_value, { message: 'Right Value is required' });
    });
  });

  ngOnInit(): void {
    if (!this.conditionNode || !this.conditionNode.condition_config || this.conditionNode.condition_config.length == 0) {


      const defaultConditionConfg: ConditionConfig = {
        left_var: {
          var_key: '',
          label: '',
          type: '',
          is_required: false,
          value: undefined,
          children: []
        },
        operator: '',
        right_value: '',
        preceeding_logic: 'AND'
      };
      this.workflowStep = signal({
        ...this.conditionNode,
        condition_config: [defaultConditionConfg]
      } as WorkflowStep);
      this.formGroup().setControlValue(this.workflowStep());
      console.log('Initialized workflowStep:', this.workflowStep());
    }
  }

  addCondition() {
    const newCondition: ConditionConfig = {
      left_var: {
        var_key: '',
        label: '',
        type: '',
        is_required: false,
        value: undefined,
        children: []
      },
      operator: '',
      right_value: '',
      preceeding_logic: 'AND'
    };

    const currentConditions = this.formGroup.condition_config().value();
    this.formGroup.condition_config().setControlValue([...currentConditions, newCondition]);
  }

  removeCondition(index: number) {
    const currentConditions = this.formGroup.condition_config().value();
    currentConditions.splice(index, 1);
    this.formGroup.condition_config().setControlValue([...currentConditions]);
  }

  saveConfiguration() {
    if (this.formGroup().valid()) {
      this.notificationService.success('Condition configuration saved successfully.', 5);
    } else {
      this.notificationService.error('Form is invalid. Please correct the errors before saving.', 5);
      return;
    }
    this.workflowStep().statement_label = this.generateStatementLabel(this.formGroup.condition_config().value());
    this.workflowStep().statement = this.generateStatement(this.formGroup.condition_config().value());
    const updatedStep: WorkflowStep = {
      ...this.workflowStep(),
      condition_config: this.formGroup.condition_config().value()
    };
    this.onSave.emit(updatedStep);
  }

  cancelConfiguration() {
    this.onCancel.emit();
  }



  generateStatementLabel(conditionConfig: any[]): string {
    return conditionConfig.map(config => {
      const leftVar = config.left_var.label || 'undefined';
      const operator = config.operator || '==';
      const rightValue = config.right_value || 'undefined';
      return `${leftVar} ${operator} ${rightValue}`;
    }).join(` ${conditionConfig[0]?.preceeding_logic || 'AND'} `);
  }

  generateStatement(conditionConfig: any[]): string {
    return conditionConfig.map(config => {
      const leftVar = config.left_var.var_key || 'undefined';
      const operator = config.operator || '==';
      const rightValue = config.right_value || 'undefined';
      return `${leftVar} ${operator} ${rightValue}`;
    }).join(` ${conditionConfig[0]?.preceeding_logic || 'AND'} `);
  }
}