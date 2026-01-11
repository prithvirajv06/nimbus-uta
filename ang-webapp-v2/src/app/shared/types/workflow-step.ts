export interface WorkflowStep {
    type: 'trigger' | 'loop' | 'condition' | 'action';
    label: string;
    icon: string;
    colorClass?: string; // e.g., 'indigo', 'amber', 'emerald'
    conditionLabel?: string;
    children?: WorkflowStep[];
    elseChildren?: WorkflowStep[]; // For the "ELSE" branch
    isOpen?: boolean;
}