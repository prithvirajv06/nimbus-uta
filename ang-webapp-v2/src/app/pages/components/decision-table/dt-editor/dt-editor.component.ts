import { Component, inject, OnInit, signal, input } from '@angular/core';
import { DtService } from '../../../../shared/services/dt.service';
import { PageBreadcrumbComponent } from '../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { MatExpansionModule } from "@angular/material/expansion";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { InputFieldComponent } from "../../../../shared/components/form/input/input-field.component";
import { TextAreaComponent } from "../../../../shared/components/form/input/text-area.component";
import { DecisionTable, DTUtils } from '../../../../shared/types/dt.type';
import { Field, form } from '@angular/forms/signals';
import { ButtonComponent } from '../../../../shared/components/ui/button/button.component';
import { CommonEditorComponent } from '../../../../shared/components/editor/common-editor.component';
import { ModalComponent } from '../../../../shared/components/ui/modal/modal.component';
import { Option, SelectComponent } from "../../../../shared/components/form/select/select.component";
import { ArrayFilter, Variable, VariablePackage } from '../../../../shared/types/variable_package';
import { VariablePackageService } from '../../../../shared/services/variable-package.service';
import { ApiResponse } from '../../../../shared/types/common.type';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { VariableSelectorComponent } from "../../variable-package/variable-selector/variable-selector.component";
import { JsonPipe } from '@angular/common';

@Component({
  selector: 'app-dt-editor',
  imports: [PageBreadcrumbComponent, MatExpansionModule, LabelComponent,
    FormsModule, ReactiveFormsModule,
    JsonPipe,
    InputFieldComponent, TextAreaComponent, Field, ButtonComponent, ModalComponent, SelectComponent, VariableSelectorComponent],
  templateUrl: './dt-editor.component.html',
  styleUrl: './dt-editor.component.css',
})
export class DtEditorComponent extends CommonEditorComponent<DecisionTable> implements OnInit {

  dtService = inject(DtService);
  varService = inject(VariablePackageService);

  isNewColumnModalOpen = false;
  newColumnType: 'Input' | 'Output' | null = null;
  variablePackage: VariablePackage | any = null;
  variableOptions: Option[] = [];
  variableToAdd: Variable | null = null;
  selectedVariableFilter: ArrayFilter[] | any = [];
  booleanOptions: Option[] = [
    { value: true, label: 'True' },
    { value: false, label: 'False' }
  ];
  override setService(): void {
    this.service = this.dtService
  }

  override setFormModel(): void {
    this.formModel = DTUtils.signalModel();
  }
  override setFormGroup(): void {
    this.formGroup = DTUtils.detailsFormGroup(this.formModel);
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

  openNewColumnModal(column: 'Input' | 'Output', variable?: Variable) {
    this.newColumnType = column;
    this.isNewColumnModalOpen = true;
    if (variable) {
      this.newColumnType = 'Input';
      this.variableToAdd = variable;
      this.isNewColumnModalOpen = true;
      if (variable.array_filters)
        this.selectedVariableFilter = variable.array_filters;
    } else {
      this.variableToAdd = null;
      this.selectedVariableFilter = [];
    }
  }

  closeNewColumnModal() {
    this.isNewColumnModalOpen = false;
    this.newColumnType = null;
  }

  setNewVariableToAdd(variable: Variable) {
    this.variableToAdd = variable;
  }

  updateVariableFilter(filter: ArrayFilter) {
    this.selectedVariableFilter.push(filter);
  }


  setNewVariableColumn(columnData: any, columnType: 'Input' | 'Output' | null) {
    if (this.variableToAdd) {
      this.variableToAdd['array_filters'] = [];
      this.variableToAdd.array_filters = this.selectedVariableFilter;
    }
    if (columnType === 'Input') {
      let exists = this.formModel().input_columns.find(ic => ic.var_key === columnData.var_key);
      if (exists) {
        this.formModel.update(dt => {
          const input_columns = dt.input_columns.map(ic => {
            if (ic.var_key === columnData.var_key) {
              return columnData;
            }
            return ic;
          }
          ); return {
            ...dt,
            input_columns
          };
        });
        this.notificationService.success('Input variable updated successfully in the decision table.', 5);
        return
      } else if (this.variableToAdd && (this.variableToAdd.type == 'object' || this.variableToAdd.type == 'array')) {
        this.notificationService.error('Input variable of type object/Array is not supported.', 10);
        return
      }
      const currentInputs = this.formModel().input_columns || [];
      this.formModel.update(dt => {
        const input_columns = [...currentInputs, columnData];
        // Update each rule by adding a new input at the correct index
        const rules = (dt.rules || []).map(rule => {
          const newRule = [...rule];
          // Insert the new variable at the end of input columns (before outputs)
          if (this.variableToAdd) {
            newRule.splice(input_columns.length - 1, 0, this.variableToAdd);
          }
          return newRule as any;
        });
        this.closeNewColumnModal();
        return {
          ...dt,
          input_columns,
          rules
        };
      });
    } else if (columnType === 'Output') {
      let exists = this.formModel().output_columns.find(ic => ic.var_key === columnData.var_key);
      if (exists) {
        this.formModel.update(dt => {
          const output_columns = dt.output_columns.map(ic => {
            if (ic.var_key === columnData.var_key) {
              return columnData;
            }
            return ic;
          }
          ); return {
            ...dt,
            output_columns
          };
        });
        this.notificationService.success('Output variable updated successfully in the decision table.', 5);
        return
      }
      const currentOutputs = this.formModel().output_columns || [];
      this.formModel.update(dt => {
        const output_columns = [...currentOutputs, columnData];
        // Update each rule by adding a new input at the correct index
        const rules = (dt.rules || []).map(rule => {
          const newRule = [...rule];
          // Insert the new variable at the end of input columns (before outputs)
          if (this.variableToAdd) {
            newRule.splice(output_columns.length + dt.input_columns.length - 1, 0, this.variableToAdd);
          }
          return newRule as any;
        });
        this.closeNewColumnModal();
        return {
          ...dt,
          output_columns,
          rules
        };
      });
    }
  }

  addRuleAndConditions() {
    const newRule = DTUtils.createEmptyRule(this.formModel().input_columns, this.formModel().output_columns);
    const currentRules = this.formModel().rules || [];
    this.formModel.update(dt => ({
      ...dt,
      rules: [...currentRules, newRule]
    }));
  }

  setBoolValue(event: any, ruleIndex: number, varIndex: number) {
    this.formModel.update(dt => {
      const updatedRules = dt.rules.map((rule, rIndex) => {
        if (rIndex === ruleIndex) {
          return rule.map((variable, vIndex) => {
            if (vIndex === varIndex) {
              return {
                ...variable,
                value: event === 'true' ? true : false
              };
            }
            return variable;
          });
        }
        return rule;
      });
      return {
        ...dt,
        rules: updatedRules
      };
    });
  }
}
