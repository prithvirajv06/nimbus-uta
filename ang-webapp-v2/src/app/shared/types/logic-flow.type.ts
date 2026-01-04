import { Form } from "@angular/forms";
import { Audit, FormGroupUtilsContract } from "./common.type";
import { Variable, VariablePackage } from "./variable_package";
import { signal, WritableSignal } from "@angular/core";
import { applyEach, form, required } from "@angular/forms/signals";

export interface LogicFlow {
    nimb_id: string;
    name: string;
    description: string;
    active: boolean;
    no_of_branches: number;
    variable_package: VariablePackage;
    logical_steps: LogicalStep[];
    audit: Audit;
}

export interface LogicalStep {
    operation_name: string;
    condition: Condition;
    variable: Variable;
    logical: string;
    op_value: any;
    operation_if_true: Operation[];
    operation_if_false: Operation[];
}

export interface Condition {
    operator: string;
    conditions: Condition[];
    variable: Variable;
    logical: string;
    op_value: any;
    array_filters?: ArrayFilter[];
}

export interface Operation {
    variable: Variable;
    operation: string;
    op_value: any;
    value_is_path: boolean;
    array_filters?: ArrayFilter[];
}

export interface ArrayFilter {
    array_name: string;
    property: string;
    logical: string;
    op_value: any;
}


export const LogicFlowUtils: FormGroupUtilsContract<LogicFlow> = {
    signalModel(): any {
        return signal<LogicFlow>({
            nimb_id: '',
            name: '',
            description: '',
            active: true,
            no_of_branches: 0,
            variable_package: {
                nimb_id: '',
                description: '',
                no_of_variables: 0,
                variables: [],
                audit: {
                    version: 0,
                    created_at: new Date(),
                    modified_at: new Date(),
                    created_by: "",
                    modified_by: "",
                    is_prod_candidate: false,
                    status: "",
                    active: "",
                    is_archived: false,
                    restore_archive: false,
                    minor_version: 0
                },
                package_name: ""
            },
            logical_steps: [{
                operation_name: "",
                condition: {
                    operator: "",
                    conditions: [
                        {
                            operator: "",
                            conditions: [],
                            variable: {
                                var_key: "",
                                label: "",
                                type: "",
                                is_required: false,
                                value: "",
                                children: []
                            },
                            logical: "",
                            op_value: undefined
                        }
                    ],
                    variable: {
                        var_key: "",
                        label: "",
                        type: "",
                        is_required: false,
                        value: "",
                        children: []
                    },
                    logical: "",
                    op_value: undefined
                },
                variable: {
                    var_key: "",
                    label: "",
                    type: "",
                    is_required: false,
                    value: "",
                    children: []
                },
                logical: "",
                op_value: undefined,
                operation_if_true: [],
                operation_if_false: []
            }],
            audit: {
                version: 0,
                created_at: new Date(),
                modified_at: new Date(),
                created_by: "",
                modified_by: "",
                is_prod_candidate: false,
                status: "",
                active: "",
                is_archived: false,
                restore_archive: false,
                minor_version: 0
            }
        })
    },
    detailsFormGroup(formModel: WritableSignal<LogicFlow>): any {
        return form<LogicFlow>(formModel, (schema) => {
            required(schema.name, { message: "Logic Flow name is required." });
            required(schema.variable_package.nimb_id, { message: "Number of branches is required." });
            required(schema.description, { message: "Description is required." });
            applyEach(schema.logical_steps, (stepSchema) => {
                required(stepSchema.operation_name, { message: "Step name is required." });
                required(stepSchema.condition.logical, { message: "Step condition logical operator is required." });
                if (stepSchema.condition && stepSchema.condition.conditions) {
                    applyEach(stepSchema.condition.conditions, (condSchema) => {
                        required(condSchema.variable, { message: "Condition variable is required." });
                        required(condSchema.op_value, { message: "Condition value is required." });
                        required(condSchema.logical, { message: "Condition logical operator is required." });
                        if (condSchema.array_filters) {
                            applyEach(condSchema.array_filters, (arrayFilterSchema) => {
                                required(arrayFilterSchema.array_name, { message: "Array name is required." });
                                required(arrayFilterSchema.property, { message: "Array property is required." });
                                required(arrayFilterSchema.logical, { message: "Array filter logical operator is required." });
                                required(arrayFilterSchema.op_value, { message: "Array filter value is required." });
                            });
                        }
                    });
                }
                // Additional validations can be added here
                if (stepSchema.operation_if_true) {
                    applyEach(stepSchema.operation_if_true, (opSchema) => {
                        required(opSchema.variable, { message: "Operation variable is required." });
                        required(opSchema.operation, { message: "Operation type is required." });
                        required(opSchema.op_value, { message: "Operation value is required." });
                        if (opSchema.array_filters) {
                            applyEach(opSchema.array_filters, (arrayFilterSchema) => {
                                required(arrayFilterSchema.array_name, { message: "Array name is required." });
                                required(arrayFilterSchema.property, { message: "Array property is required." });
                                required(arrayFilterSchema.logical, { message: "Array filter logical operator is required." });
                                required(arrayFilterSchema.op_value, { message: "Array filter value is required." });
                            });
                        }
                    });
                }
                if (stepSchema.operation_if_false) {
                    applyEach(stepSchema.operation_if_false, (opSchema) => {
                        required(opSchema.variable, { message: "Operation variable is required." });
                        required(opSchema.operation, { message: "Operation type is required." });
                        required(opSchema.op_value, {
                            message: "Operation value is required."
                        });
                        if (opSchema.array_filters) {
                            applyEach(opSchema.array_filters, (arrayFilterSchema) => {
                                required(arrayFilterSchema.array_name, { message: "Array name is required." });
                                required(arrayFilterSchema.property, { message: "Array property is required." });
                                required(arrayFilterSchema.logical, { message: "Array filter logical operator is required." });
                                required(arrayFilterSchema.op_value, { message: "Array filter value is required." });
                            });
                        }
                    });
                }

            });
        });
    },
    basicFormGroup(formModel: WritableSignal<LogicFlow>): any {
        return form<LogicFlow>(formModel, (schema) => {
            required(schema.name, { message: "Logic Flow name is required." });
            required(schema.variable_package.nimb_id, { message: "Number of branches is required." });
            required(schema.description, { message: "Description is required." });
        });
    }
}