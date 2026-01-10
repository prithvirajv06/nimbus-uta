import { start } from '@popperjs/core';
import { NgClass } from '@angular/common';
import { Component, Host, HostListener, input, Input } from '@angular/core';

@Component({
  selector: 'app-workflow-item',
  imports: [NgClass],
  styles: [`
    .animate-in {
      animation: fadeIn 0.2s ease-out;
    }

    @keyframes fadeIn {
      from { opacity: 0; transform: translateX(-10px); }
      to { opacity: 1; transform: translateX(0); }
    }
    `],
  templateUrl: './workflow-item.component.html',
})
export class WorkflowItemComponent {
  @Input() data: any[] = [];

  @Input() selectedItem: string | null = null;

  density: 'compact' | 'comfortable' = 'comfortable';
  explainMode = false;
  editingId?: string;
  draggingId?: string;

  // Helper to define styles based on action type
  getConfigs(type: string) {
    const configs: FlowNode = {
      start: { icon: 'fa-play', color: 'bg-green-500', border: 'border-green-100', text: 'Start', isChildApplicable: true },
      assignment: { icon: 'fa-equals', color: 'bg-blue-500', border: 'border-blue-100', text: 'Set', isChildApplicable: false },
      push_array: { icon: 'fa-plus', color: 'bg-emerald-500', border: 'border-emerald-100', text: 'Push', isChildApplicable: false },
      condition: { icon: 'fa-code-branch', color: 'bg-amber-500', border: 'border-amber-200', text: 'If', isChildApplicable: true },
      for_each: { icon: 'fa-sync-alt', color: 'bg-indigo-600', border: 'border-indigo-200', text: 'For Each', isChildApplicable: true }
    };
    return (<any>configs)[type] || configs['assignment'];
  }


}

export interface FlowNode {
  start: TaskConfig;
  assignment: TaskConfig;
  push_array: TaskConfig;
  condition: TaskConfig;
  for_each: TaskConfig;
}
export interface TaskConfig {
  icon: string;
  color: string;
  border: string;
  text: string;
  isChildApplicable: boolean;
}