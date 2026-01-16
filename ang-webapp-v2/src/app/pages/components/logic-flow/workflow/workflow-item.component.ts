
import { CommonModule, JsonPipe, NgClass } from '@angular/common';
import { Component, Host, HostListener, input, Input } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatExpansionModule } from "@angular/material/expansion";
import { ɵEmptyOutletComponent } from "@angular/router";
import { DragDropModule } from '@angular/cdk/drag-drop';
import { moveItemInArray, CdkDragDrop } from '@angular/cdk/drag-drop';
import { ConditionConfigureComponent } from './condition-configure/condition-configure.component';
import { WorkflowStep } from '../../../../shared/types/workflow-step';
import { SideDrawComponent } from "../../../../shared/components/ui/side-draw/side-draw.component";
import { AssignmentConfigureComponent } from './assignment-configure/assignment-configure.component';
import { Variable, VariablePackage } from '../../../../shared/types/variable_package';
import { LoopConfigureComponent } from './loop-configure/loop-configure.component';

@Component({
  selector: 'app-workflow-item',
  imports: [CommonModule, FormsModule, MatExpansionModule, ɵEmptyOutletComponent,
    DragDropModule, ConditionConfigureComponent, NgClass, SideDrawComponent,
    LoopConfigureComponent,
    AssignmentConfigureComponent],
  styleUrls: ['./workflow-item.component.css'],
  templateUrl: './workflow-editor.component.html',
})
export class WorkflowItemComponent {
  @Input() workflowData: WorkflowStep[] = [] as WorkflowStep[];
  @Input() isRoot: boolean = true;
  @Input() variablePackage: VariablePackage = {} as VariablePackage;
  parentVar: Variable = {} as Variable;
  isAddStepOpen: { [key: string]: boolean } = {};
  configDrawTitle = '';
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
      target: {} as Variable,
      true_children: [],
      false_children: [],
      statement: '',
      statement_label: ''
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

  findActualParentArray(node: any, currentArray: any[]): any[] | null {
    for (let item of currentArray) {
      if (item.children && item.children.includes(node)) {
        return item.children;
      }
      if (item.true_children && item.true_children.includes(node)) {
        return item.true_children;
      }
      if (item.false_children && item.false_children.includes(node)) {
        return item.false_children;
      }
      if (item.children) {
        const found = this.findActualParentArray(node, item.children);
        if (found) return found;
      }
      if (item.true_children) {
        const found = this.findActualParentArray(node, item.true_children);
        if (found) return found;
      }
      if (item.false_children) {
        const found = this.findActualParentArray(node, item.false_children);
        if (found) return found;
      }
    }
    return null;
  }

  getDynamicBorder() {
    return this.isRoot ? 'border-indigo-300' : 'border-slate-300';
  }


  getAllowedVariables(node: WorkflowStep): any {
    const parentArray = this.findParentArray(node);
    if (parentArray) {
    }
  }




  //Condition vars and methods
  isConfigPanOpen: boolean = false;
  selectedNodeForConfig!: WorkflowStep;
  configType = 'condition';

  openConditionConfig(node: WorkflowStep, parentVar: Variable): void {
    this.selectedNodeForConfig = node;
    this.isConfigPanOpen = true;
    this.configType = 'condition';
    this.getAllowedVariables(node);
    this.parentVar = parentVar;
    this.configDrawTitle = 'Configure Condition';
  }

  closeConditionConfig(): void {
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }

  saveConditionConfig(updatedNode: WorkflowStep): void {
    Object.assign(this.selectedNodeForConfig, updatedNode);
    this.closeConditionConfig();
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }


  //Loop COnfig
  openLoopConfig(node: WorkflowStep, parentVar: Variable): void {
    this.selectedNodeForConfig = node;
    this.isConfigPanOpen = true;
    this.configType = 'for_each';
    this.parentVar = parentVar;
    this.configDrawTitle = 'Configure Loop';
  }

  closeLoopConfig(): void {
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }

  saveLoopConfig(updatedNode: WorkflowStep): void {
    Object.assign(this.selectedNodeForConfig, updatedNode);
    this.closeLoopConfig();
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }


  //Assignment Config
  openAssignmentConfig(node: WorkflowStep, parentVar: Variable): void {
    this.selectedNodeForConfig = node;
    this.isConfigPanOpen = true;
    this.configType = 'assignment';
    this.parentVar = parentVar;
    this.configDrawTitle = 'Configure Assignment';
  }

  closeAssignmentConfig(): void {
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }

  saveAssignmentConfig(updatedNode: WorkflowStep): void {
    Object.assign(this.selectedNodeForConfig, updatedNode);
    this.closeAssignmentConfig();
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }
  deleteConfig(deletedNode: WorkflowStep): void {
    const parentArray = this.findActualParentArray(deletedNode, this.workflowData);
    if (parentArray) {
      const index = parentArray.indexOf(deletedNode);
      if (index > -1) {
        parentArray.splice(index, 1);
      }
    }
    this.closeAssignmentConfig();
    this.isConfigPanOpen = false;
    this.selectedNodeForConfig = {} as WorkflowStep;
  }
}
