
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
    DragDropModule, ModalComponent, ConditionConfigureComponent, NgClass],
  styleUrls: ['./workflow-item.component.css'],
  templateUrl: './workflow-editor.component.html',
})
export class WorkflowItemComponent extends WorkflowStepHooks {
  @Input() workflowData: WorkflowStep[] = [] as WorkflowStep[];
  @Input() isRoot: boolean = true;
  isAddStepOpen: { [key: string]: boolean } = {};
  stepOptions = [
    { type: 'assignment', label: 'Assignment', icon: 'fas fa-edit' },
    { type: 'for_each', label: 'Loop for Each', icon: 'fas fa-sync-alt' },
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


  openStepOption(step: WorkflowStep, prefix: string = "") {
    this.isAddStepOpen[step.step_id + prefix] = true;
  }

  addStep(step: WorkflowStep, option: any, isTrue: boolean = false, isFalse: boolean = false, prefix: string = "") {
    const newStep: WorkflowStep = {
      step_id: this.generateUniqueId(),
      type: option.type,
      label: '',
      icon: '',
      condition_config: [],
      children: [],
      target: '',
      true_children: [],
      false_children: []
    };
    if (isTrue) {
      if (!step.true_children) {
        step.true_children = [];
      }
      step.true_children.push(newStep);
      this.isAddStepOpen[step.step_id + prefix] = false;
      return;
    }
    if (isFalse) {
      if (!step.false_children) {
        step.false_children = [];
      }
      step.false_children.push(newStep);
      this.isAddStepOpen[step.step_id] = false;
      return;
    }
    if (!step.children) {
      step.children = [];
    }
    step.children.push(newStep);
    this.isAddStepOpen[step.step_id] = false;
  }

  generateUniqueId(): string {
    return 'step-' + Math.random().toString(36).substr(2, 9);
  }


  moveStepUp(node: any) {
    const arr = this.findParentArray(node);
    const idx = arr.indexOf(node);
    if (idx > 0) {
      [arr[idx - 1], arr[idx]] = [arr[idx], arr[idx - 1]];
    }
  }

  moveStepDown(node: any) {
    const arr = this.findParentArray(node);
    const idx = arr.indexOf(node);
    if (idx < arr.length - 1) {
      [arr[idx], arr[idx + 1]] = [arr[idx + 1], arr[idx]];
    }
  }

  // Helper to find the parent array of a node
  findParentArray(node: any): any[] {
    // Example: search workflowData and all children recursively
    // You may need to adjust this logic based on your data structure
    const search = (arr: any[]): any[] | null => {
      if (arr.includes(node)) return arr;
      for (const item of arr) {
        if (item.children) {
          const found = search(item.children);
          if (found) return found;
        }
        if (item.true_children) {
          const found = search(item.true_children);
          if (found) return found;
        }
        if (item.false_children) {
          const found = search(item.false_children);
          if (found) return found;
        }
      }
      return null;
    };
    return search(this.workflowData) || this.workflowData;
  }

  getDynamicBorder() {
    return this.isRoot ? 'border-indigo-300' : 'border-slate-300';
  }
}
