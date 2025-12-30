// DecisionTableStruct.ts

import { Audit } from "./common.type";
import { VariablePackage } from "./variable_package";

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