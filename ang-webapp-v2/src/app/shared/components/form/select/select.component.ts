
import { Component, Input, Output, EventEmitter, OnInit, OnChanges, SimpleChanges } from '@angular/core';
import { FormsModule } from '@angular/forms';

export interface Option {
  value: string | boolean | number;
  label: string;
  applicableTo?: string[];
}

@Component({
  selector: 'app-select',
  imports: [FormsModule],
  templateUrl: './select.component.html',
})
export class SelectComponent implements OnInit, OnChanges {

  @Input() options: Option[] = [];
  @Input() placeholder: string = 'Select an option';
  @Input() className: string = '';
  @Input() defaultValue: string = '';
  @Input() value: string = '';
  @Input() applicableTo: string = '';
  filteredOptions: Option[] = [];
  @Output() valueChange = new EventEmitter<string>();

  ngOnInit() {
    if (!this.value && this.defaultValue) {
      this.value = this.defaultValue;
    }
    this.filteredOptions = this.options;
  }
  ngOnChanges(changes: SimpleChanges): void {
    if (this.applicableTo) {
      this.filteredOptions = this.options.filter(option => {
        return !option.applicableTo || option.applicableTo.includes(this.applicableTo);
      });
    } else {
      this.filteredOptions = this.options;
    }
  }
  onChange(event: Event) {
    const value = (event.target as HTMLSelectElement).value;
    this.value = value;
    this.valueChange.emit(value);
  }
}