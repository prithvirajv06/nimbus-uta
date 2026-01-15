import { Variable } from "./variable_package";

export interface WorkflowStep {
    step_id: string;
    type: 'trigger' | 'loop' | 'condition' | 'action' | 'start' | 'end' | 'assignment' | 'push_array' | 'for_each';
    label: string;
    icon: string;
    target: string;
    value?: string;
    statement?: string;
    statementLabel?: string;
    children?: WorkflowStep[];
    true_children?: WorkflowStep[]; // For the "ELSE" branch
    false_children?: WorkflowStep[]; // For the "IF" branch
    isOpen?: boolean;
    context_var?: string; // For loop context variable
    condition_config: ConditionConfig[]
}

export interface ConditionConfig {
    left_var: Variable;
    operator: string;
    right_value: any;
    preceeding_logic: string; // 'AND' | 'OR'
}