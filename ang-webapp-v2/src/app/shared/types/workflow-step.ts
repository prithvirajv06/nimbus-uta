import { Variable } from "./variable_package";

export interface WorkflowStep {
    step_id: string;
    type: string; // 'start' | 'task' | 'logic' | 'loop' | 'message' | 'end'
    label: string;
    icon: string;
    target: Variable; // For Loop will use this to store the target array to loop through
    value?: any;
    children: WorkflowStep[];
    true_children: WorkflowStep[]; // For the "ELSE" branch
    false_children: WorkflowStep[]; // For the "IF" branch
    isOpen?: boolean;
    context_var?: string; // For loop context variable local variable, In backend will replace the content_var in the target array loop functions
    condition_config: ConditionConfig[],
    context_map: LoopContextMap[];    
    statement: string;
    statement_label: string;
    loop_level?: number;
}
export interface LoopContextMap {
    var_key: string;//policy.customers
    context_key: string; // customer_1
    level:number// 1 | 2 | 3
}

export interface ConditionConfig {
    left_var: Variable;
    operator: string;
    right_value: any;
    preceeding_logic: string; // 'AND' | 'OR'
}