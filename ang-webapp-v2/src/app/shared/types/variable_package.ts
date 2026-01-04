// Typescript interfaces for variable_pacakage_struct.go

import { signal, WritableSignal } from "@angular/core";
import { Audit, FormGroupUtilsContract } from "./common.type";
import { form, required, applyEach, FieldTree, pattern } from "@angular/forms/signals";

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
  value: any;
  children: Variable[];
  array_filters?: ArrayFilter[]
}

export interface ArrayFilter {
  array_name: string;
  property: string;
  logical: string;
  op_value: string;
}

export interface VariablePackageUtilContract extends FormGroupUtilsContract<VariablePackage> {
  options(variablePackages: VariablePackage[]): { label: string; value: string }[];
  newVariableFormGroup(): FieldTree<Variable>;
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
  options: function (variablePackages: VariablePackage[]) {
    return variablePackages.map((vp) => ({
      label: vp.package_name + " v" + vp.audit.version + "." + vp.audit.minor_version,
      value: vp.nimb_id + "~" + vp.audit.version
    }));
  },
  newVariableFormGroup: function (): FieldTree<Variable> {
    let field = signal<Variable>({
      var_key: '',
      label: '',
      type: 'string',
      is_required: false,
      value: "",
      children: []
    });
    return form(field, (schemaPath) => {
      required(schemaPath.label, { message: 'Variable Name is required' });
      required(schemaPath.type, { message: 'Variable Type is required' });
      required(schemaPath.var_key, { message: 'Variable Key is required' });
      pattern(
        schemaPath.var_key,
        /^([a-zA-Z_][a-zA-Z0-9_]*|\[\*\])(\[(\*|\d+)\])*(\.([a-zA-Z_][a-zA-Z0-9_]*|\[\*\])(\[(\*|\d+)\])*)*$/,
        {
          message:
            "Variable Key must be a valid JSON path (e.g., key, obj.key, arr[*].key). Each segment must start with a letter or underscore, or be an array wildcard [*].",
        }
      );
    });
  }
}