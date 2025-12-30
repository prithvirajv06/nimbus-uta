// Typescript interfaces for variable_pacakage_struct.go

import { Audit } from "./common.type";

export interface VariablePackageRequest {
  package_name: string;
  description: string;
  json_str: string;
}

export interface VariablePackage {
  nimb_id: string;
  package_name: string;
  description: string;
  variables: Variable[];
  no_of_variables: number;
  audit: Audit;
}

export interface Variable {
  var_key: string;
  label: string;
  type: string;
  is_required: boolean;
}
