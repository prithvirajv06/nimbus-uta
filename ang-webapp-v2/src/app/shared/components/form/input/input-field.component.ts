import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { CommonControlComponent } from './common-control';

@Component({
  selector: 'app-input-field',
  imports: [CommonModule],
  template: `
    <div class="relative">
      <input
        [type]="type"
        [id]="id"
        [placeholder]="placeholder"
        [step]="step"
        [ngClass]="inputClasses"
        [value]="value()"
        (input)="value.set($event.target.value)"
      />

      @if (errors().length > 0 && (dirty()||touched())) {
      <p
        class="mt-2 text-sm"
        [ngClass]="errors().length > 0 ? 'text-error-500' : 'text-gray-500 dark:text-gray-400'">
        @for (err of errors(); track $index) {
        {{ err.message }}
        }
      </p>
      }
    </div>
  `,
})
export class InputFieldComponent extends CommonControlComponent{

  @Input() type: string = 'text';
  @Input() id?: string = '';
  @Input() placeholder?: string = '';
  @Input() step?: number;
  @Input() success: boolean = false;
  @Input() error: boolean = false;
  @Input() hint?: string;
  @Input() className: string = '';

  get inputClasses(): string {
    let inputClasses = `h-11 w-full rounded-lg border appearance-none px-4 py-2.5 text-sm shadow-theme-xs placeholder:text-gray-400 focus:outline-hidden focus:ring-3 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 ${this.className}`;

    if (this.error) {
      inputClasses += ` border-error-500 focus:border-error-300 focus:ring-error-500/20 dark:text-error-400 dark:border-error-500 dark:focus:border-error-800`;
    } else if (this.success) {
      inputClasses += ` border-success-500 focus:border-success-300 focus:ring-success-500/20 dark:text-success-400 dark:border-success-500 dark:focus:border-success-800`;
    } else {
      inputClasses += ` bg-transparent text-gray-800 border-gray-300 focus:border-brand-300 focus:ring-brand-500/20 dark:border-gray-700 dark:text-white/90  dark:focus:border-brand-800`;
    }
    return inputClasses;
  }

}
