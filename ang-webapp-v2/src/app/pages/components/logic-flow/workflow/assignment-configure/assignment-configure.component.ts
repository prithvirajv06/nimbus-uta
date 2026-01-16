import { Component, Input } from '@angular/core';
import { WorkflowStep } from '../../../../../shared/types/workflow-step';
import { Variable } from '../../../../../shared/types/variable_package';

@Component({
  selector: 'app-assignment-configure',
  imports: [],
  templateUrl: './assignment-configure.component.html',
  styleUrl: './assignment-configure.component.css',
})
export class AssignmentConfigureComponent {

  @Input() assignmentNode: WorkflowStep = {} as WorkflowStep;
  @Input() variables: Variable[] = [];
}
