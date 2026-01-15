// flow-builder.component.ts
import { Component, HostListener, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PageBreadcrumbComponent } from '../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { MatExpansionModule } from "@angular/material/expansion";
import { WorkflowItemComponent } from "../workflow/workflow-item.component";
import { start } from '@popperjs/core';
import { WorkflowStep } from '../../../../shared/types/workflow-step';

interface FlowNode {
  id: string;
  type: 'start' | 'task' | 'logic' | 'loop' | 'message' | 'end';
  label: string;
  config?: any;
  children?: FlowNode[];
  trueBranch?: FlowNode[];
  falseBranch?: FlowNode[];
}

@Component({
  selector: 'app-logic-flow-editor',
  standalone: true,
  imports: [CommonModule, FormsModule, PageBreadcrumbComponent, MatExpansionModule, WorkflowItemComponent],
  templateUrl: './logic-flow-editor.component.html',
  styles: []
})
export class LogicFlowEditorComponent {
  data: WorkflowStep[] | any = [
    {
      type: 'start',
      label: '',
      icon: '',
      target: ''
    },
  ];

  // Helper to define styles based on action type
  getConfigs(type: string) {
    const configs: any = {
      start: { icon: 'fa-play', color: 'bg-green-500', border: 'border-green-100', text: 'Start' },
      assignment: { icon: 'fa-equals', color: 'bg-blue-500', border: 'border-blue-100', text: 'Set' },
      push_array: { icon: 'fa-plus', color: 'bg-emerald-500', border: 'border-emerald-100', text: 'Push' },
      condition: { icon: 'fa-code-branch', color: 'bg-amber-500', border: 'border-amber-200', text: 'If' },
      for_each: { icon: 'fa-sync-alt', color: 'bg-indigo-600', border: 'border-indigo-200', text: 'For Each' }
    };
    return configs[type] || configs['assignment'];
  }
}