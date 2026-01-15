
import { CommonModule, JsonPipe, NgClass } from '@angular/common';
import { Component, Host, HostListener, input, Input } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatExpansionModule } from "@angular/material/expansion";
import { ɵEmptyOutletComponent } from "@angular/router";
import { DragDropModule } from '@angular/cdk/drag-drop';
import { moveItemInArray, CdkDragDrop } from '@angular/cdk/drag-drop';
import { ModalComponent } from "../../../../shared/components/ui/modal/modal.component";
import { ConditionConfigureComponent } from './condition-configure/condition-configure.component';
import { WorkflowStep } from '../../../../shared/types/workflow-step';
import { ConditionHooks } from './hooks/condition.hooks';
import { WorkflowStepHooks } from './hooks/workflow.step.hooks';

@Component({
  selector: 'app-workflow-item',
  imports: [JsonPipe, CommonModule, FormsModule, MatExpansionModule, ɵEmptyOutletComponent,
    DragDropModule, ModalComponent, ConditionConfigureComponent],
  styleUrls: ['./workflow-item.component.css'],
  templateUrl: './workflow-editor.component.html',
})
export class WorkflowItemComponent extends WorkflowStepHooks {
  @Input() workflowData: WorkflowStep[] = [] as WorkflowStep[];
  @Input() isRoot: boolean = true;
  isAddStepOpen: { [key: string]: boolean } = {};
stepOptions = [
    { type: 'assignment', label: 'Assignment', icon: 'fas fa-edit' },
    { type: 'push_array', label: 'Push to Array', icon: 'fas fa-plus-square' },
    { type: 'condition', label: 'Condition', icon: 'fas fa-code-branch' },
  ];
  // ... inside your component class
  drop(event: CdkDragDrop<any[]>) {
    // moveItemInArray is a helper from CDK that handles the index logic
    moveItemInArray(
      event.container.data,
      event.previousIndex,
      event.currentIndex
    );
  }


  openStepOption(step: WorkflowStep) {
    this.isAddStepOpen[step.step_id] = true;
  }

  addStep(step: WorkflowStep, option: any) {
    const newStep: WorkflowStep = {
      step_id: this.generateUniqueId(),
      type: option.type,
      label: '',
      icon: '',
      condition_config: [],
      children: [],
      target: ''
    };
    if (!step.children) {
      step.children = [];
    }
    step.children.push(newStep);
    this.isAddStepOpen[step.step_id] = false;
  }

  generateUniqueId(): string {
    return 'step-' + Math.random().toString(36).substr(2, 9);
  }
}
