// DecisionTableStruct.ts

import { Signal, signal, WritableSignal } from "@angular/core";
import { Audit, FormGroupUtilsContract } from "./common.type";
import { Variable, VariablePackage } from "./variable_package";
import { apply, applyEach, form, required, schema } from "@angular/forms/signals";

export interface DecisionTable {
  nimb_id: string;
  description: string;
  no_of_rows: number;
  no_of_inputs: number;
  no_of_outputs: number;
  name: string;
  hit_policy: string;
  input_columns: Variable[];
  output_columns: Variable[];
  variable_package: VariablePackage;
  rules: Variable[][];
  audit: Audit;
}

export interface DTRuleUtilsContract extends FormGroupUtilsContract<DecisionTable> {
  createEmptyRule(inputColumns: Variable[], outputColumns: Variable[]): Variable[];
}

export const DTUtils: DTRuleUtilsContract = {
  signalModel(): WritableSignal<DecisionTable> {
    return signal<DecisionTable>({
      nimb_id: '',
      description: "",
      no_of_rows: 0,
      no_of_inputs: 0,
      no_of_outputs: 0,
      name: "",
      hit_policy: "",
      input_columns: [],
      output_columns: [],
      variable_package: {
        nimb_id: "",
        package_name: "",
        description: "",
        variables: [],
        no_of_variables: 0,
        audit: {
          created_at: new Date(),
          created_by: "",
          modified_at: new Date(),
          modified_by: "",
          is_prod_candidate: false,
          status: "",
          active: "",
          is_archived: false,
          restore_archive: false,
          version: 0,
          minor_version: 0
        }
      },
      rules: [],
      audit: {
        created_at: new Date(),
        created_by: "",
        modified_at: new Date(),
        modified_by: "",
        is_prod_candidate: false,
        status: "",
        active: "",
        is_archived: false,
        restore_archive: false,
        version: 0,
        minor_version: 0
      }
    });
  },
  basicFormGroup(formModel: WritableSignal<DecisionTable>): any {
    return form<DecisionTable>(formModel, (schema) => {
      required(schema.name, { message: 'Decision Table Name is required' });
      required(schema.description, { message: 'Description is required' });
      required(schema.hit_policy, { message: 'Hit Policy is required' });
      required(schema.variable_package.nimb_id, { message: 'Variable Package is required' });
    });
  },

  detailsFormGroup(formModel: WritableSignal<DecisionTable>): any {
    return form<DecisionTable>(formModel, (schema) => {
      required(schema.name, { message: 'Decision Table Name is required' });
      required(schema.description, { message: 'Description is required' });
      required(schema.hit_policy, { message: 'Hit Policy is required' });
      required(schema.input_columns, { message: 'At least one Input Column is required' });
      required(schema.output_columns, { message: 'At least one Output Column is required' });
      required(schema.rules, { message: 'At least one Rule is required' });
      applyEach(schema.input_columns, (inputSchema) => {
        required(inputSchema.var_key, { message: 'Input Variable is required' });
        required(inputSchema.label, { message: 'Input Label is required' });
      });
      applyEach(schema.output_columns, (outputSchema) => {
        required(outputSchema.var_key, { message: 'Output Variable is required' });
        required(outputSchema.label, { message: 'Output Label is required' });
      });
    });
  },
  createEmptyRule(inputColumns: Variable[], outputColumns: Variable[]): Variable[] {
    const rule: Variable[] = [];
    inputColumns.forEach((variable) => {
      variable.value = '';
      rule.push({ ...variable });
    });
    outputColumns.forEach((variable) => {
      rule.push({ ...variable });
    });
    return rule;
  }
}