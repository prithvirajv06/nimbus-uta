import { start } from '@popperjs/core';
import { CommonModule, JsonPipe, NgClass } from '@angular/common';
import { Component, Host, HostListener, input, Input } from '@angular/core';
import { WorkflowStep } from '../../../../shared/types/workflow-step';
import { FormsModule } from '@angular/forms';
import { MatExpansionModule } from "@angular/material/expansion";
import { ɵEmptyOutletComponent } from "@angular/router";

@Component({
  selector: 'app-workflow-item',
  imports: [JsonPipe, CommonModule, FormsModule, MatExpansionModule, ɵEmptyOutletComponent],
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
  @Input() workflowData: any;
  @Input() isRoot: boolean = true;
}
