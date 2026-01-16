import { AfterContentInit, Component, EventEmitter, Input, Output } from '@angular/core';
import { Variable } from '../../../../../shared/types/variable_package';
import { WorkflowStep } from '../../../../../shared/types/workflow-step';
import { VariableSelectorComponent } from '../../../variable-package/variable-selector/variable-selector.component';
import { LabelComponent } from "../../../../../shared/components/form/label/label.component";
import { ButtonComponent } from '../../../../../shared/components/ui/button/button.component';

@Component({
  selector: 'app-loop-configure',
  imports: [VariableSelectorComponent, LabelComponent, ButtonComponent],
  templateUrl: './loop-configure.component.html',
  styleUrl: './loop-configure.component.css',
})
export class LoopConfigureComponent implements AfterContentInit {
  @Input() loopNode: WorkflowStep = {} as WorkflowStep;
  @Input() parentVar: Variable = {} as Variable;
  @Input() variables: Variable[] = [];
  applicableVariables: Variable[] = [];
  selectedVar: Variable = {} as Variable;
  @Output() onSave = new EventEmitter<WorkflowStep>();
  @Output() onCancel = new EventEmitter<void>();


  constructor() { }


  ngAfterContentInit(): void {
    this.getArrayVariables(this.variables);
    this.filterVariableBasedOnHirarchy();
  }

  getArrayVariables(variable: Variable[]) {
    for (let varItem of variable) {
      if (varItem.type === 'array') {
        varItem.children = [];
        this.applicableVariables.push(varItem);
      }
      if (varItem.children && varItem.children.length > 0) {
        this.getArrayVariables(varItem.children);
      }
    }
  }

  filterVariableBasedOnHirarchy() {
    if (this.parentVar)
      this.applicableVariables.forEach((varItem) => {
        if (this.parentVar && !this.parentVar.var_key.includes(varItem.var_key)) {
          const index = this.applicableVariables.indexOf(varItem);
          if (index > -1) {
            this.applicableVariables.splice(index, 1);
          }
        }
      });
  }


  saveLoopConfiguration() {
    this.loopNode.target = this.selectedVar;
    this.loopNode.context_var = this.selectedVar.var_key.split('.').pop() || '';
    this.loopNode.statement_label = `For Each details in ${this.selectedVar.label}`;
    this.onSave.emit(this.loopNode);
  }


}
