// DecisionTableStruct.ts

import { Signal, signal, WritableSignal } from "@angular/core";
import { Audit } from "./common.type";
import { VariablePackage } from "./variable_package";
import { apply, applyEach, form, required, schema } from "@angular/forms/signals";

export interface DecisionTable {
  nimb_id: string;
  description: string;
  no_of_rows: number;
  no_of_inputs: number;
  no_of_outputs: number;
  name: string;
  hit_policy: string;
  input_columns: TableInput[];
  output_columns: TableOutput[];
  variable_package: VariablePackage;
  rules: string[][];
  audit: Audit;
}

export interface TableInput {
  variable: string;
  label: string;
}

export interface TableOutput {
  variable: string;
  label: string;
  allowed_values?: string[];
  is_priority: boolean;
}



export const DTModelutils = {
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
  basicFormGroup(formModel: WritableSignal<DecisionTable>) {
    return form<DecisionTable>(formModel, (schema) => {
      required(schema.name, { message: 'Decision Table Name is required' });
      required(schema.description, { message: 'Description is required' });
      required(schema.hit_policy, { message: 'Hit Policy is required' });
      required(schema.variable_package.nimb_id, { message: 'Variable Package is required' });
    });
  },

  detailsFormGroup(formModel: WritableSignal<DecisionTable>) {
    return form<DecisionTable>(formModel, (schema) => {
      required(schema.name, { message: 'Decision Table Name is required' });
      required(schema.description, { message: 'Description is required' });
      required(schema.hit_policy, { message: 'Hit Policy is required' });
      required(schema.input_columns, { message: 'At least one Input Column is required' });
      required(schema.output_columns, { message: 'At least one Output Column is required' });
      required(schema.rules, { message: 'At least one Rule is required' });
      applyEach(schema.input_columns, (inputSchema) => {
        required(inputSchema.variable, { message: 'Input Variable is required' });
        required(inputSchema.label, { message: 'Input Label is required' });
      });
      applyEach(schema.output_columns, (outputSchema) => {
        required(outputSchema.variable, { message: 'Output Variable is required' });
        required(outputSchema.label, { message: 'Output Label is required' });
      });
    });
  }
}