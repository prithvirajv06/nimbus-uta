import { WorkflowStep } from "../../../../../shared/types/workflow-step";

export class ConditionHooks {
 
  //Condition vars and methods
  isConditionConfigOpen: boolean = false;
  selectedConditionNode!: WorkflowStep;

  openConditionConfig(node: WorkflowStep): void {
    this.selectedConditionNode = node;
    this.isConditionConfigOpen = true;
  }

  closeConditionConfig(): void {
    this.isConditionConfigOpen = false;
    this.selectedConditionNode = {} as WorkflowStep;
  }

  saveConditionConfig(updatedNode: WorkflowStep): void {
    Object.assign(this.selectedConditionNode, updatedNode);
    this.selectedConditionNode!.statementLabel = this.generateStatementLabel(this.selectedConditionNode.condition_config);
    this.selectedConditionNode!.statement = this.generateStatement(this.selectedConditionNode.condition_config);
    this.closeConditionConfig();
    this.isConditionConfigOpen = false;
    this.selectedConditionNode = {} as WorkflowStep;
  }

  generateStatementLabel(conditionConfig: any[]): string {
    return conditionConfig.map(config => {
      const leftVar = config.left_var.label || 'undefined';
      const operator = config.operator || '==';
      const rightValue = config.right_value || 'undefined';
      return `${leftVar} ${operator} ${rightValue}`;
    }).join(` ${conditionConfig[0]?.preceeding_logic || 'AND'} `);
  }

  generateStatement(conditionConfig: any[]): string {
    return conditionConfig.map(config => {
      const leftVar = config.left_var.var_key || 'undefined';
      const operator = config.operator || '==';
      const rightValue = config.right_value || 'undefined';
      return `${leftVar} ${operator} ${rightValue}`;
    }).join(` ${conditionConfig[0]?.preceeding_logic || 'AND'} `);
  }
}