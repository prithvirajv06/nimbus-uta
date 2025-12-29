import { CommonModule } from '@angular/common';
import { Component, Input, Output, EventEmitter, model, ChangeDetectionStrategy, input } from '@angular/core';
import { DisabledReason, FormCheckboxControl, ValidationError } from '@angular/forms/signals';

@Component({
  selector: 'app-checkbox',
  imports: [CommonModule],
  template: `
  <label
  class="flex items-center space-x-3 group cursor-pointer"
  [ngClass]="{ 'cursor-not-allowed opacity-60': disabled() }"
>
  <div class="relative w-5 h-5">
    <input
      [id]="id"
      type="checkbox"
      class="w-5 h-5 appearance-none cursor-pointer dark:border-gray-700 border border-gray-300 checked:border-transparent rounded-md checked:bg-brand-500 disabled:opacity-60"
      [ngClass]="className"
      [checked]="checked()"
      (change)="toggle()"
      [disabled]="disabled()"
    />
    @if (checked()) {
    <ng-container>
      <svg
        class="absolute transform -translate-x-1/2 -translate-y-1/2 pointer-events-none top-1/2 left-1/2"
        xmlns="http://www.w3.org/2000/svg"
        width="14"
        height="14"
        viewBox="0 0 14 14"
        fill="none"
      >
        <path
          d="M11.6666 3.5L5.24992 9.91667L2.33325 7"
          stroke="white"
          stroke-width="1.94437"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </ng-container>
    }
    @if (disabled()) {
    <ng-container>
      <svg
        class="absolute transform -translate-x-1/2 -translate-y-1/2 pointer-events-none top-1/2 left-1/2"
        xmlns="http://www.w3.org/2000/svg"
        width="14"
        height="14"
        viewBox="0 0 14 14"
        fill="none"
      >
        <path
          d="M11.6666 3.5L5.24992 9.91667L2.33325 7"
          stroke="#E4E7EC"
          stroke-width="2.33333"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </ng-container>
    }
  </div>
  @if (label) {
  <span
    class="text-sm font-medium text-gray-800 dark:text-gray-200"
    >
      {{ label }}
  </span>
  }
</label>
  `,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CheckboxComponent implements FormCheckboxControl {

  @Input() label?: string;
  @Input() className = '';
  @Input() id?: string;

  checked = model<boolean>(false);
  toggle() {
    this.checked.update(val => !val);
  }

  disabled = input<boolean>(false);
  readonly = input<boolean>(false);
  hidden = input<boolean>(false);
  invalid = input<boolean>(false);
}
