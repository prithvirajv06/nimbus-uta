// Typescript interfaces for variable_pacakage_struct.go

import { signal, WritableSignal } from "@angular/core";
import { Audit, FormGroupUtilsContract } from "./common.type";
import { form, required, applyEach } from "@angular/forms/signals";

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
  value: string;
}

export interface VariablePackageUtilContract extends FormGroupUtilsContract<VariablePackage> {
  options(variablePackages: VariablePackage[]): { label: string; value: string }[];
}
export const VariablePackageUtils: VariablePackageUtilContract = {
  signalModel(): WritableSignal<VariablePackage> {
    return signal<VariablePackage>({
      nimb_id: '',
      package_name: '',
      description: '',
      variables: [] as Variable[],
      audit: {
        version: 0,
        created_at: new Date(),
        modified_at: new Date(),
        created_by: '',
        modified_by: '',
        status: 'DRAFT',
        is_prod_candidate: false,
        active: '',
        is_archived: false,
        minor_version: 0,
        restore_archive: false
      },
      no_of_variables: 0
    });
  },
  basicFormGroup: function (formModel: WritableSignal<VariablePackage>) {
    return form(formModel, (schemaPath) => {
      required(schemaPath.package_name, { message: 'Package Name is required' });
      required(schemaPath.description, { message: 'Description is required' });
      applyEach(schemaPath.variables, (variableSchema) => {
        required(variableSchema.label, { message: 'Variable Name is required' });
        required(variableSchema.type, { message: 'Variable Type is required' });
        required(variableSchema.var_key, { message: 'Variable Value is required' });
      });
    });
  },
  detailsFormGroup: function (formModel: WritableSignal<VariablePackage>) {
    return form(formModel, (schemaPath) => {
      required(schemaPath.package_name, { message: 'Package Name is required' });
      required(schemaPath.description, { message: 'Description is required' });
      applyEach(schemaPath.variables, (variableSchema) => {
        required(variableSchema.label, { message: 'Variable Name is required' });
        required(variableSchema.type, { message: 'Variable Type is required' });
        required(variableSchema.var_key, { message: 'Variable Value is required' });
      });
    });
  },
  options: function(variablePackages: VariablePackage[]) {
    return variablePackages.map((vp) => ({
      label: vp.package_name + " v" + vp.audit.version + "." + vp.audit.minor_version,
      value: vp.nimb_id + "~" + vp.audit.version
    }));
  }
}