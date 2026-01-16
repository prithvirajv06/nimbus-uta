import { Form } from "@angular/forms";
import { Audit, FormGroupUtilsContract } from "./common.type";
import { Variable, VariablePackage } from "./variable_package";
import { signal, WritableSignal } from "@angular/core";
import { applyEach, form, required } from "@angular/forms/signals";
import { WorkflowStep } from "./workflow-step";

export interface LogicFlow {
    nimb_id: string;
    name: string;
    description: string;
    active: boolean;
    no_of_branches: number;
    variable_package: VariablePackage;
    logical_steps: WorkflowStep[];
    audit: Audit;
}

export interface LogicFormGroupUtils<T> {
    signalModel(): WritableSignal<T>,
    basicFormGroup(formModel: WritableSignal<T>): any,
    detailsFormGroup(formModel: WritableSignal<T>): any,
    formModuleWithNoveltyValidations(formModel: WritableSignal<T>): any
}

export const LogicFlowUtils: LogicFormGroupUtils<LogicFlow> = {
    signalModel(): any {
        return signal<LogicFlow>({
            nimb_id: "",
            name: "",
            description: "",
            active: false,
            no_of_branches: 0,
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
            logical_steps: [],
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
        } );
    },
    detailsFormGroup(formModel: WritableSignal<LogicFlow>): any {
        function validateStep(step: any) {
            required(step.target, { message: 'Step target is required' });
            applyEach(step.condition_config, (config: any) => {
                required(config.left_var, { message: 'Left Variable is required' });
                required(config.operator, { message: 'Operator is required' });
                required(config.right_value, { message: 'Right Value is required' });
            });
            if (step.children && Array.isArray(step.children)) {
                applyEach(step.children, validateStep);
            }
        }

        return form<LogicFlow>(formModel, (schema) => {
            required(schema.name, { message: "Logic Flow name is required." });
            required(schema.variable_package.nimb_id, { message: "Number of branches is required." });
            required(schema.description, { message: "Description is required." });
            applyEach(schema.logical_steps, validateStep);
        });
    },
    formModuleWithNoveltyValidations(formModel: WritableSignal<LogicFlow>): any {
        return form<LogicFlow>(formModel, (schema) => {
            required(schema.name, { message: "Logic Flow name is required." });
            required(schema.variable_package.nimb_id, { message: "Number of branches is required." });
            required(schema.description, { message: "Description is required." });
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