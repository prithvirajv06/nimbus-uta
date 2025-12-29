import { CommonModule } from '@angular/common';
import { Component, Input, Output, EventEmitter, model } from '@angular/core';
import { FormValueControl } from '@angular/forms/signals';

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
        (input)="value.set($event.target.value)"
        [value]="value()"
      />

      @if (hint) {
      <p class="mt-1.5 text-xs"
        [ngClass]="{
          'text-error-500': error,
          'text-success-500': success,
          'text-gray-500': !error && !success
        }">
        {{ hint }}
      </p>
      }
    </div>
  `,
})
export class InputFieldComponent implements FormValueControl<string | number> {

  @Input() type: string = 'text';
  @Input() id?: string = '';
  @Input() placeholder?: string = '';
  @Input() step?: number;
  @Input() success: boolean = false;
  @Input() error: boolean = false;
  @Input() hint?: string;
  @Input() className: string = '';
  value = model<string | number>('');

  @Output() valueChange = new EventEmitter<string | number>();

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

  onInput(event: Event) {
    const input = event.target as HTMLInputElement;
    this.valueChange.emit(this.type === 'number' ? +input.value : input.value);
  }
}
