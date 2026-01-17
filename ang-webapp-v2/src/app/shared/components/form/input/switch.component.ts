import { CommonModule } from '@angular/common';
import { Component, Input, Output, EventEmitter, OnInit, input } from '@angular/core';
import { CommonControlComponent } from './common-control';

@Component({
  selector: 'app-switch',
  imports: [
    CommonModule
  ],
  template: `
   <label
      class="flex cursor-pointer select-none items-center gap-3 text-sm font-medium"
      [ngClass]="disabled() ? 'text-gray-400' : 'text-gray-700 dark:text-gray-400'"
    >
    {{leftLabel}}
      <div class="relative" (click)="handleToggle()">
        <div
          class="block transition duration-150 ease-linear h-6 w-11 rounded-full"
          [ngClass]="
            (disabled()
              ? 'bg-gray-100 pointer-events-none dark:bg-gray-800'
              : switchColors.background)
          "
        ></div>
        <div
          class="absolute left-0.5 top-0.5 h-5 w-5 rounded-full shadow-theme-sm duration-150 ease-linear transform"
          [ngClass]="switchColors.knob"
        ></div>
      </div>
      {{ label }}
    </label>
  `
})
export class SwitchComponent implements OnInit {
  disabled = input<boolean>(false);
  @Input() label!: string;
  @Input() leftLabel: string = 'Off';
  @Input() defaultChecked: boolean = false;
  @Input() color: 'blue' | 'gray' = 'blue';
  @Output() valueChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  isChecked: boolean = false;

  ngOnInit() {
    this.isChecked = this.defaultChecked;
  }

  handleToggle() {
    if (this.disabled()) return;
    this.isChecked = !this.isChecked;
    this.valueChange.emit(this.isChecked);
  }

  get switchColors() {
    if (this.color === 'blue') {
      return {
        background: this.isChecked
          ? 'bg-brand-500'
          : 'bg-gray-200 dark:bg-white/10',
        knob: this.isChecked
          ? 'translate-x-full bg-white'
          : 'translate-x-0 bg-white',
      };
    } else {
      return {
        background: this.isChecked
          ? 'bg-gray-800 dark:bg-white/10'
          : 'bg-gray-200 dark:bg-white/10',
        knob: this.isChecked
          ? 'translate-x-full bg-white'
          : 'translate-x-0 bg-white',
      };
    }
  }
}
