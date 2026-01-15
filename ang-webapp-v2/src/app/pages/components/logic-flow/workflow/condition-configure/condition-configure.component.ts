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

  @Input() conditionNode: WorkflowStep | null = null;
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
    if (!this.conditionNode || this.conditionNode.condition_config.length == 0) {


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
      this.notificationService.success('Condition configuration saved successfully.');
    } else {
      this.notificationService.error('Form is invalid. Please correct the errors before saving.');
      return;
    }
    const updatedStep: WorkflowStep = {
      ...this.workflowStep(),
      condition_config: this.formGroup.condition_config().value()
    };
    this.onSave.emit(updatedStep);
  }

  cancelConfiguration() {
    this.onCancel.emit();
  }

  //** Mock Variable */
  variables: Variable[] = [
    {
      var_key: 'user.age', label: 'User Age', type: 'number', children: [
        {
          var_key: 'user.age.years', label: 'User Age in Years', type: 'number',
          is_required: false,
          value: undefined,
          children: []
        },
        {
          var_key: 'user.age.months', label: 'User Age in Months', type: 'number',
          is_required: false,
          value: undefined,
          children: []
        },
      ],
      is_required: false,
      value: undefined
    },
    {
      var_key: 'user.name', label: 'User Name', type: 'string',
      is_required: false,
      value: undefined,
      children: []
    },
    {
      var_key: 'order.total', label: 'Order Total', type: 'number',
      is_required: false,
      value: undefined,
      children: [{
        var_key: 'order.total.amount', label: 'Order Total Amount', type: 'number',
        is_required: false,
        value: undefined,
        children: [
          {
            var_key: 'order.total.amount.currency', label: 'Currency', type: 'string',
            is_required: false,
            value: undefined,
            children: []
          },
        ]
      }]
    },
  ];
}
